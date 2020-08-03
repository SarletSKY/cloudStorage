package api

import (
	"filestore-server-study/common"
	"filestore-server-study/config"
	dbCli "filestore-server-study/service/dbproxy/client"
	"filestore-server-study/store/ceph"
	"filestore-server-study/store/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// 下载的handler

// 下载文件
func DownLoadFile(c *gin.Context) {

	fmt.Println("zhaoweixiong：")
	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	// 获取具体文件信息
	//fileMeta := meta.GetFileMeta(filehash)
	dbfResp, fErr := dbCli.GetFileInfoTodb(filehash)
	dbfUserResp, fUserErr := dbCli.QueryUserFileDB(username, filehash)
	if fErr != nil || fUserErr != nil || !dbfResp.Suc || !dbfUserResp.Suc {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server error",
		})
		return
	}

	fileMeta := dbCli.ToTableFile(dbfResp.Data)
	userFileInfo := dbCli.ToTableUserFile(dbfUserResp.Data)
	// TODO: 6. 对下载接口进行修改
	// TODO: 7. 对下载方式进行判断ceph
	var fileBytes []byte
	if strings.HasPrefix(fileMeta.FileAdd.String, config.MergeLocalRootDir) {
		// 打开文件
		file, err := os.Open(fileMeta.FileAdd.String)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 读出文件数据
		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fileMeta.FileAdd.String, config.CephRootDir) { // 从ceph下载
		fmt.Println("to download file from ceph...")
		bucket := ceph.GetCephBucket("userfile")
		var err error
		fileBytes, err = bucket.Get(fileMeta.FileAdd.String)
		if err != nil {
			fmt.Println(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fileMeta.FileAdd.String, config.OSSRootDir) { // 从oss下载
		fmt.Println("to download file from oss...")
		var err1 error
		var err2 error
		var rc io.ReadCloser
		rc, err1 = oss.Bucket().GetObject(fileMeta.FileAdd.String)
		if err1 == nil {
			fileBytes, err2 = ioutil.ReadAll(rc)
			if err2 == nil {
			}
		}
		if err1 != nil || err2 != nil {
			c.Header("content-disposition", "attachment; filename=\""+userFileInfo.FileName+"\"")
			c.Data(http.StatusInternalServerError, "application/octet-stream", fileBytes)
		}
	}

	// 写数据到前端页面去
	c.Header("content-disposition", "attachment; filename=\""+userFileInfo.FileName+"\"")
	c.FileAttachment(fileMeta.FileAdd.String, userFileInfo.FileName)
	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}

// 支持断点的文件下载接口
func RangeDownload(c *gin.Context) {

	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	// 获取具体文件信息
	//fileMeta := meta.GetFileMeta(filehash)
	dbfResp, fErr := dbCli.GetFileInfoTodb(filehash)
	dbfUserResp, fUserErr := dbCli.QueryUserFileDB(username, filehash)

	if fErr != nil || fUserErr != nil || !dbfResp.Suc || !dbfUserResp.Suc {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server error",
		})
		return
	}
	// TODO: 6. 对下载接口进行修改   userFileInfo
	userFileInfo := dbCli.ToTableUserFile(dbfUserResp.Data)

	// TODO: 8. 使用本地目录文件
	fpath := config.MergeLocalRootDir + filehash
	fmt.Println("range-download-fpath: " + fpath)

	// 打开文件
	file, err := os.Open(fpath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  err.Error(),
		})
		return
	}
	defer file.Close()

	// 写数据到前端页面去
	c.Writer.Header().Set("Content-Type", "application/octet-stream") // 下载二进制流
	c.Writer.Header().Set("content-disposition", "attachment; filename=\""+userFileInfo.FileName+"\"")
	c.File(fpath)
}

// 生成文件下载地址
func DownloadURL(c *gin.Context) {
	// TODO: 8. 对下载地址进行修改
	filehash := c.Request.FormValue("filehash")
	// 从文件表查找信息
	dbResp, err := dbCli.GetFileInfoTodb(filehash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server error",
		})
	}

	row := dbCli.ToTableFile(dbResp.Data)

	// TODO: 8. 在oss中下载文件，不加上下载不了，因为要跨域请求。已经转移到auth.go文件
	// 进行判断是oss还是ceph下载地址
	if strings.HasPrefix(row.FileAdd.String, config.MergeLocalRootDir) || strings.HasPrefix(row.FileAdd.String, config.CephRootDir) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		downloadURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Data(http.StatusOK, "application/octet-stream", []byte(downloadURL))
	} else if strings.HasPrefix(row.FileAdd.String, config.OSSRootDir) {
		signedURL := oss.DownloadURL(row.FileAdd.String)
		c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))
	} else {
		c.Data(http.StatusOK, "application/octet-stream", []byte("ERROR: 下载链接错误"))
	}
}
