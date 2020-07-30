package handler

import (
	"context"
	userProto "filestore-server-study/service/account/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//获取多个元信息数据
func GetManyFileMetaInfo(c *gin.Context) {

	limitCount := c.Request.FormValue("limit")
	// 转换int类型
	limit, _ := strconv.Atoi(limitCount)

	//fileMeta := meta.GetLastFileMeta(limit)
	//5.4 升级为批量查询用户文件接口
	//fileMeta, err := meta.GetLastFileMetaDB(limit)
	username := c.Request.FormValue("username")
	rpcResp, err := userCli.UserFiles(context.TODO(), &userProto.ReqUserFile{
		Username: username,
		Limit:    int32(limit),
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}

	fmt.Println("------")
	fmt.Println(rpcResp)
	fmt.Println(rpcResp.FileData)

	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}

// 更新文件元信息(重命名)
func UpdateFileInfo(c *gin.Context) {

	// 通过sha1获取文件的元信息 op是指客户端需要操作的类型的标志
	opType := c.Request.FormValue("op")
	filehash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")

	// TODO: 6. 进行优化[将更新用户文件表元信息同时也要修改]
	username := c.Request.FormValue("username")

	if opType != "0" || len(fileName) == 0 { // 0表示复制或者修改操作
		c.Status(http.StatusForbidden)
		return
	}

	rpcResp, err := userCli.UserFileRename(context.TODO(), &userProto.ReqUserFileRename{
		Username:    username,
		Filehash:    filehash,
		NewFileName: fileName,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}

	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}
