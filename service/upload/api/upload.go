package api

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

// 上传/秒传

// 初始化api
func init() {
	if err := os.MkdirAll(TempDir, 0744); err != nil {
		fmt.Println("not found mkdir file" + TempDir)
		os.Exit(1)
	}
	if err := os.MkdirAll(MergeDir, 0744); err != nil {
		fmt.Println("found mkdir file" + MergeDir)
		os.Exit(1)
	}
}

// 上传文件[POST]
func DoUploadHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	errCode := 0
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传成功",
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

	// TODO: 上传文件之前要确保用户文件名不会重复
	exist := db.QueryUserFileNameExist(username, fileMeta.FileName)
	if exist {
		errCode = -7
		return
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
	suc := meta.UploadFileMetaDB(fileMeta)
	if !suc {
		errCode = -6
		return
	}

	// 5.3 升级上传接口,将文件上传到用户文件表上
	// 解析上下文获取username

	suc = db.OnUserFileUploadFinshedDB(username, fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize, fileMeta.UpdateFileTime)
	if !suc {
		errCode = -6
	} else {
		errCode = 0
	}

	// 2.1 处理成功页面
	// 2.1 成功上传就进行重定向
	//http.Redirect(writer, request, "/file/upload/success", http.StatusFound) // 重定向的状态码
	// 5.1 跳转到登录页面

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
	suc := db.OnUserFileUploadFinshedDB(username, filename, filehash, filesize, time.Now().Format("2006-01-02 15:04:05"))
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
