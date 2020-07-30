package handler

import (
	"encoding/json"
	"filestore-server-study/common"
	"filestore-server-study/config"
	"filestore-server-study/db"
	"filestore-server-study/meta"
	"filestore-server-study/mq"
	"filestore-server-study/store/ceph"
	"filestore-server-study/store/oss"
	"filestore-server-study/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 上传文件[GET]
func UploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		c.String(http.StatusBadRequest, "网页不存在")
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(data))
}

// 上传文件[POST]
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "Upload failed",
			})
		}
	}()

	// 2.1 接受get文件的数据 FormFile 是与前端对接的数据
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err:%s\n", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	// 2.2 对数据进行储存[就是对元信息进行初始化赋值]
	// TODO: 7. 对存进本地/tmp/路径修改打牌ceph
	fileMeta := meta.FileMeta{
		FileName:       header.Filename,
		Location:       config.TempLocalRootDir + header.Filename,
		UpdateFileTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 2.1 备份数据到本地 （利用copy进行处理）
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		errCode = -1
		fmt.Println(err.Error())
		return
	}
	defer newFile.Close()

	// 2.1 io.Copy返回fileSize数据
	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		errCode = -2
		fmt.Println(err.Error())
		return
	}
	// 2.2 文件进行sha1加密,并且添加到Meta元信息map中 注意：计算hash之前，一定要将seek移动到开头
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	//TODO: 7. 同步或异步将文件转移到ceph/oss
	// 7.1 ceph集群的设置
	newFile.Seek(0, 0)
	mergePath := config.MergeLocalRootDir + fileMeta.FileSha1
	if config.CurrentStoreType == common.StoreCeph {
		//文件存储到ceph
		// 读出文件数据
		data, _ := ioutil.ReadAll(newFile)
		cephFilePath := "/ceph/" + fileMeta.FileSha1
		err = ceph.PutObject("userfile", cephFilePath, data)
		if err != nil {
			fmt.Println("upload ceph err: " + err.Error())
			errCode = -3
			fmt.Println(err.Error())
			return
		}
		fileMeta.Location = cephFilePath
	} else if config.CurrentStoreType == common.StoreOSS {
		ossPath := "oss/" + fileMeta.FileSha1
		// oss存储分两种，异步与同步
		if !config.AsyncTransferEnable {
			err = oss.Bucket().PutObject(ossPath, newFile)
			if err != nil {
				fmt.Println("upload ceph err: " + err.Error())
				errCode = -4
				fmt.Println(err.Error())
				return
			}
			fileMeta.Location = ossPath
		} else {
			// TODO: 9. 加入rabbitMQ队列，先经过mq，再经过oss
			/*				// 注意：文件会先存入本地，将任务加入队列，加入oss之前，在将本地路径修改掉
							fileMeta.Location = mergePath*/
			// 解析msg数据,序列化数据.
			data := mq.TransferData{
				FileHash:      fileMeta.FileSha1,
				CurLocation:   fileMeta.Location,
				DestLocation:  ossPath,
				DestStoreType: common.StoreOSS,
			}
			msg, _ := json.Marshal(data)
			// 先生成生产者
			suc := mq.Publish(config.TransExchangeName,
				config.TransOSSRoutingKey,
				msg,
			)
			fmt.Println(suc)
			if !suc {
				// TODO: 当前发送转移信息失败，稍后重试
			}
		}
	} else {
		fileMeta.Location = mergePath
	}
	/*		// 读出文件数据
			data, _ := ioutil.ReadAll(newFile)
			bucket := ceph.GetCephBucket("userFile")
			// 设置ceph文件路径
			cephFilePath := "/ceph/" + fileMeta.FileSha1

			// 写入到ceph集群
			_ = bucket.Put(cephFilePath, data, "octet-stream", s3.PublicRead)
			// 路径改成ceph,以后提取往这提取
			fileMeta.Location = cephFilePath*/

	//meta.UploadFileMeta(fileMeta)
	_ = meta.UploadFileMetaDB(fileMeta)

	// 5.3 升级上传接口,将文件上传到用户文件表上
	// 解析上下文获取username

	username := c.Request.FormValue("username")
	suc := db.OnUserFileUploadFinshedDB(username, fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize)
	if !suc {
		errCode = -5
	}

	// 2.1 处理成功页面
	// 2.1 成功上传就进行重定向
	//http.Redirect(writer, request, "/file/upload/success", http.StatusFound) // 重定向的状态码
	// 5.1 跳转到登录页面
	c.Redirect(http.StatusFound, "/static/view/home.html")
}

// 秒上传的接口
func FastUploadUserFile(c *gin.Context) {

	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.ParseInt(c.Request.FormValue("filesize"), 10, 64)
	filehash := c.Request.FormValue("filehash")

	// 向文件表中查找有没有上传过
	fileMeta, err := db.GetFileInfoTodb(filehash)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 查不到数据，秒传数据失败
	if fileMeta == nil {
		// 返回前端
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请使用普通上传功能",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	// 成功则秒传[上传用户文件表]
	suc := db.OnUserFileUploadFinshedDB(username, filename, filehash, filesize)
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败,请稍后重试",
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	return
}

//返回成功页面
func UploadSucHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Upload Finish!",
	})
}

//获取元信息路由 获取时在终端使用命令[sha1sum 文件路径]获取sha1加密
func GetFileMetaInfo(c *gin.Context) {

	// 2.3 根据url上的filehash参数来去Meta中寻找元信息， 并进行json序列化返回
	filehash := c.Request.FormValue("filehash") // 根据url上的参数对应来赋值
	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -2,
			"msg":  "Upload failed",
		})
		return
	}

	if fileMeta != nil {
		fileMtaBytes, err := json.Marshal(fileMeta)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -3,
				"msg":  "Upload failed",
			})
			return
		}
		c.Data(http.StatusOK, "application/json", fileMtaBytes)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "No sub file",
		})
	}
}

//获取多个元信息数据
func GetManyFileMetaInfo(c *gin.Context) {

	limitCount := c.Request.FormValue("limit")
	// 转换int类型
	limit, _ := strconv.Atoi(limitCount)

	//fileMeta := meta.GetLastFileMeta(limit)
	//5.4 升级为批量查询用户文件接口
	//fileMeta, err := meta.GetLastFileMetaDB(limit)
	username := c.Request.FormValue("username")
	userFile, err := db.QueryMantUserFileDB(username, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "Query failed",
		})
		return
	}
	// 序列化数据
	fileMetaBytes, err := json.Marshal(userFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -2,
			"msg":  "Query failed",
		})
		return
	}
	c.Data(http.StatusOK, "application/json", fileMetaBytes)
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

	// TODO: 6. 添加修改用户文件表, 并更改接口：文件元信息更新。也就是重命名，不需要更改文件表
	_ = db.UpdateUserFileInfoDB(username, fileName, filehash)
	/*	//curFileMeta := meta.GetFileMeta(filehash)
		curFileMeta, err := meta.GetFileMetaDB(filehash)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}


		// 修改文件名称名进行sha1哈系修改, 并加到fileMetaMap中
		curFileMeta.FileName = fileName
		//meta.UploadFileMeta(curFileMeta)
		meta.UpdateFileMetaDB(curFileMeta)*/

	// 返回成功页面
	// 序列化数据
	// TODO: 6. 将用户文件表中更改的那条数据重新获取出来，序列化返回
	userFile, err := db.QueryUserFileDB(username, filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	//fileMetaBytes, err := json.Marshal(curFileMeta)
	fileMetaBytes, err := json.Marshal(userFile)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, fileMetaBytes)
}

// 删除文件元信息
func DeleteFile(c *gin.Context) {

	// 删除备份
	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")

	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = os.Remove(fileMeta.Location)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// 删除元信息的map
	//meta.DeleteFileMeta(filehash)
	/*	if suc := meta.DeleteFileMetaDB(filehash); !suc {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}*/
	// TODO: 6. 修改删除接口，不要删除文件表的数据，而是删除文件表的数据
	suc := db.DeleteUserFileDB(username, filehash)
	if !suc {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
