package util

import (
	"filestore-server-study/config"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	// 通过shell合并分块文件
	MergeFileCMD = `
	#!/bin/bash
	# 需要进行合并的分片所在的目录
	chunkDir=$1
	# 合并后的文件的完成路径(目录＋文件名)
	mergePath=$2
	
	echo "分块合并，输入目录: " $chunkDir
	
	if [ ! -f $mergePath ]; then
			echo "$mergePath not exist"
	else
			rm -f $mergePath
	fi
	
	for chunk in $(ls $chunkDir | sort -n)
	do
			cat $chunkDir/${chunk} >> ${mergePath}
	done
	
	echo "合并完成，输出：" $mergePath
	`
	// shell命令计算hash
	FileSha1CMD = `
	#!/bin/bash
	sha1sum $1 | awk '{print $1}'
	`

	// shell命令计算size
	FileSizeCMD = `
	#!/bin/bash
	ls -l $1 | awk '{print $5}'
	`

	// 删除分块的文件[取消操作]
	FileBlockDelCMD = `
	#!/bin/bash
	blockDir="#CHUNKDIR#"
	testDir="$1"
	# 进行判断
	if [[ $testDir =~ $blockDir ]] && [[ $testDir != $blockDir ]]; then
		rm -fr $testDir
	fi
	`
)

// 使用bash命令来将文件数据删除,某路径下的文件
func RemovePathByShell(delPath string) bool {
	// 1. 将接收的文件路径进行传惨到变量中
	cmdStr := strings.Replace(FileBlockDelCMD, "$1", delPath, 1)
	// 2. 提取该bash命令进行cmd执行
	cmd := exec.Command("bash", "-c", cmdStr)
	if _, err := cmd.Output(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// 通过shell命令将文件合并
func MergeBlocksByShell(blockDir string, mergeDir string, filehash string) bool {
	// 合并分块
	cmdStr := strings.Replace(MergeFileCMD, "$1", blockDir, 1)
	cmdStr = strings.Replace(MergeFileCMD, "#CHUNKDIR#", config.BlockLocalRootDir, 1)
	// 执行shell
	cmd := exec.Command("bash", "-c", cmdStr)
	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		return false
	}

	// 计算合并后的hash，判断是否一致
	if MergeFilehash, err := ComputeSha1ByShell(mergeDir); err != nil {
		fmt.Println(err)
		return false
	} else if MergeFilehash != filehash {
		fmt.Println("hash different")
		fmt.Printf("mergeFileHash:%v\n sourceHash:%v\n", MergeFilehash, filehash)
		return false
	} else {
		fmt.Println("check sha1: " + mergeDir + " " + MergeFilehash + " " + filehash)
		return true
	}
}

// 判断hash是狗一致
func ComputeSha1ByShell(mergeDir string) (string, error) {
	cmdStr := strings.Replace(FileSha1CMD, "$1", mergeDir, 1)
	cmd := exec.Command("bash", "-c", cmdStr)
	if filehash, err := cmd.Output(); err != nil {
		fmt.Println(err)
		return "", err
	} else {
		// 正则表达式
		reg := regexp.MustCompile("\\s+")
		return reg.ReplaceAllString(string(filehash), ""), nil
	}
}

// 判断size大小
func ComputeFileSizeByShell(mergeDir string) (int, error) {
	cmdStr := strings.Replace(FileSizeCMD, "$1", mergeDir, 1)
	cmd := exec.Command("bash", "-c", cmdStr)
	if filehash, err := cmd.Output(); err != nil {
		fmt.Println(err)
		return -1, err
	} else {
		// 正则表达式
		reg := regexp.MustCompile("\\s+")
		size, err := strconv.Atoi(reg.ReplaceAllString(string(filehash), ""))
		if err != nil {
			return -1, err
		}
		return size, nil
	}
}
