package handler

import (
	"filestore-server-study/config"
	"filestore-server-study/db"
	"filestore-server-study/meta"
	"filestore-server-study/store/ceph"
	"filestore-server-study/store/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// 下载文件
func DownLoadFile(c *gin.Context) {

	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	// 获取具体文件信息
	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	// TODO: 6. 对下载接口进行修改
	userFileInfo, err := db.QueryUserFileDB(username, filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	// TODO: 7. 对下载方式进行判断ceph
	var fileBytes []byte
	if strings.HasPrefix(fileMeta.Location, config.MergeLocalRootDir) {
		// 打开文件
		file, err := os.Open(fileMeta.Location)
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
	} else if strings.HasPrefix(fileMeta.Location, "/ceph") { // 从ceph下载
		fmt.Println("to download file from ceph...")
		bucket := ceph.GetCephBucket("userfile")
		fileBytes, err = bucket.Get(fileMeta.Location)
		if err != nil {
			fmt.Println(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fileMeta.Location, "oss") { // 从oss下载
		fmt.Println("to download file from oss...")
		rc, err := oss.Bucket().GetObject(fileMeta.Location)
		if err != nil {
			fmt.Println(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		} else {
			fileBytes, err = ioutil.ReadAll(rc)
		}
	}

	// 写数据到前端页面去
	c.Writer.Header().Set("Content-Type", "application/octet-stream") // 下载二进制流
	c.Writer.Header().Set("content-disposition", "attachment; filename=\""+userFileInfo.FileName+"\"")
	c.FileAttachment(fileMeta.Location, userFileInfo.FileName)
	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}

// 支持断点的文件下载接口
func RangeDownload(c *gin.Context) {

	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	// 获取具体文件信息
	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	// TODO: 6. 对下载接口进行修改
	userFileInfo, err := db.QueryUserFileDB(username, filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// TODO: 8. 使用本地目录文件
	fpath := config.MergeLocalRootDir + fileMeta.FileSha1
	fmt.Println("range-download-fpath: " + fpath)

	// 打开文件
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 写数据到前端页面去
	c.Writer.Header().Set("Content-Type", "application/octet-stream") // 下载二进制流
	c.Writer.Header().Set("content-disposition", "attachment; filename=\""+userFileInfo.FileName+"\"")
	http.ServeFile(c.Writer, c.Request, fileMeta.Location)
}

// 生成文件下载地址
func DownloadURL(c *gin.Context) {
	// TODO: 8. 对下载地址进行修改
	filehash := c.Request.FormValue("filehash")
	// 从文件表查找信息
	row, _ := db.GetFileInfoTodb(filehash)

	// TODO: 8. 在oss中下载文件，不加上下载不了，因为要跨域请求。已经转移到auth.go文件
	// 进行判断是oss还是ceph下载地址
	if strings.HasPrefix(row.FileAdd.String, config.MergeLocalRootDir) || strings.HasPrefix(row.FileAdd.String, "/ceph") {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		downloadURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Data(http.StatusOK, "octet-stream", []byte(downloadURL))
	} else if strings.HasPrefix(row.FileAdd.String, "oss/") {
		signedURL := oss.DownloadURL(row.FileAdd.String)
		c.Data(http.StatusOK, "octet-stream", []byte(signedURL))
	} else {
		c.String(http.StatusOK, "ERROR: 下载链接错误")
	}
}
