package mpupload

import (
	"encoding/json"
	"filestore-server-study/common"
	"filestore-server-study/util"
	"flag"
	"fmt"
	jsonit "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func normalUsage() {
	fmt.Fprintf(os.Stderr, `<Test normal multipart upload>

	Usage: ./test_upload -tcase normal [-fpath filePath] [-user username] [-pwd password]

Options:
`)
	flag.PrintDefaults()
}

// 正常上传函数
func normalTestMain() {
	// 判断文件路径是否存在,hash,size是否正确
	if exist, err := util.PathExists(uploadFilePath); !exist || err != nil {
		fmt.Println("Error: 无效文件路径，请检查")
		normalUsage()
		return
	} else if len(username) == 0 || len(password) == 0 {
		fmt.Println("Error: 无效的用户名或密码")
		normalUsage()
		return
	}

	filehash, err := util.ComputeSha1ByShell(uploadFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize, err := util.ComputeFileSizeByShell(uploadFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 登录
	token, err := signIn(username, password)
	if err != nil {
		fmt.Println(err)
		return
	} else if token == "" {
		fmt.Println("Error: 登录失败，请检查用户名/密码")
		return
	}
	fmt.Println("token" + token)

	// 请求初始化上传接口
	resp, err := http.PostForm(apiUploadInit, url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filesize": {strconv.Itoa(filesize)},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	respCode := jsonit.Get(body, "code").ToInt()
	if respCode == int(common.FileAlreadExists) {
		fmt.Println("文件已经存在，无需上传")
		return
	}

	// 得到uploadId服务指定的分块大小blockSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	blockSize := jsonit.Get(body, "data").Get("BlockSize").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadID, blockSize)

	// 请求分块上传接口

	var initResp UploadInitResponse
	err = json.Unmarshal(body, &initResp)

	if err != nil {
		fmt.Printf("Parse error: %s\n", err.Error())
		os.Exit(-1)
	}
	var blocksToUpload []int
	for index := 1; index <= initResp.Data.BlockCount; index++ {
		blocksToUpload = append(blocksToUpload, index)
	}

	uploadBlockCount = len(blocksToUpload)
	tURL := apiUploadPart + "?username=" + username +
		"&token=" + token + "&uploadid=" + uploadID
	// 上传所有分块
	uploadPartsSpecified(uploadFilePath, tURL, blockSize, blocksToUpload)

	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		apiUploadComplete,
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {strconv.Itoa(filesize)},
			"filename": {filepath.Base(uploadFilePath)},
			"uploadid": {uploadID},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	// 5. 打印分块上传结果
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}
