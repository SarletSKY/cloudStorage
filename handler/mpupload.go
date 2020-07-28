package handler

// mp--->Multipart多部分

import (
	redisPool "filestore-server-study/cache/redis"
	"filestore-server-study/config"
	"filestore-server-study/db"
	"filestore-server-study/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct {
	FileHash    string
	FileSize    int
	UploadID    string
	BlockSize   int
	BlockCount  int
	BlockExists []int // 已经上传完成的分块索引列表
}

// 将数据string进行封装起来成常量
const (
	BlockKeyPrefix    = "MP_"                    // 分块上传的redis前缀
	BlockDir          = config.BlockLocalRootDir // 分块的所在目录
	MergeDir          = config.MergeLocalRootDir // 合并的目录
	HashUpIDKeyPrefix = "HASH_UPID_"             // 文件hash映射的uploadId对应的redis的前缀
)

// 初始化[也就是判断有没有这些文件]
func init() {
	if err := os.MkdirAll(BlockDir, 0744); err != nil {
		fmt.Println("not found mkdir file" + BlockDir)
		os.Exit(1)
	}
	if err := os.MkdirAll(MergeDir, 0744); err != nil {
		fmt.Println("found mkdir file" + MergeDir)
		os.Exit(1)
	}
}

// 初始化分块上传
func InitMultipartUpload(c *gin.Context) {

	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "params invalid",
		})
		return
	}
	// 2. 查询用户有没有上传过文件，判断文件是否存在
	if db.IsUserFileUpload(username, filehash) {
		c.JSON(http.StatusOK, gin.H{
			"code": 10006,
			"msg":  "OK",
		})
		return
	}

	// 3. 获取redis连接
	// 获取pool的连接
	conn := redisPool.GetRedisPool().Get()
	defer conn.Close()

	// 4. 通过文件hash判断是否断点续传，并获取uploadID
	uploadId := ""
	keyExists, _ := redis.Bool(conn.Do("EXISTS", HashUpIDKeyPrefix+filehash)) // redis是否存在该hash
	if keyExists {
		uploadId, err = redis.String(conn.Do("GET", HashUpIDKeyPrefix+filehash))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "Upload part failed",
			})
			return
		}
	}

	// 5.1 首次上传则新建uploadID
	// 5.2 断点续传则根据uploadID获取已上传的文件分块列表
	BlockExist := []int{} // 块完成的数量
	if uploadId == "" {
		uploadId = username + fmt.Sprintf("%x", time.Now().UnixNano())
	} else {
		blocks, err := redis.Values(conn.Do("HGETALL", BlockKeyPrefix+uploadId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "Upload part failed",
			})
			return
		}

		// 续传
		for i := 0; i < len(blocks); i += 2 {
			k := string(blocks[i].([]byte))
			v := string(blocks[i+1].([]byte))
			if strings.HasPrefix(k, "block_") && v == "1" {
				//block_6 -> 6
				blockIndex, _ := strconv.Atoi(k[7:]) //？？？？这个7？？？？
				BlockExist = append(BlockExist, blockIndex)
			}
		}
	}

	// 6. 初始化分块信息
	mpInfo := MultipartUploadInfo{
		FileHash:    filehash,
		FileSize:    filesize,
		UploadID:    uploadId,        //ID使用用户名加时间戳
		BlockSize:   5 * 1024 * 1024, // 5MB
		BlockCount:  int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
		BlockExists: BlockExist,
	}
	// 6. 上传到redis
	if len(mpInfo.BlockExists) <= 0 {
		hkey := BlockKeyPrefix + mpInfo.UploadID
		conn.Do("HSET", hkey, "blockcount", mpInfo.BlockCount)
		conn.Do("HSET", hkey, "filehash", mpInfo.FileHash)
		conn.Do("HSET", hkey, "filesize", mpInfo.FileSize)
		conn.Do("EXPIRE", hkey, 43200) // 半天的时间，时间一过就会清除。
		conn.Do("SET", HashUpIDKeyPrefix+filehash, mpInfo.UploadID, "EX", 43200)
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": mpInfo,
	})
}

// 实现分块上传
func MultipartUpload(c *gin.Context) {
	fmt.Println("??")
	//username := request.Form.Get("username")
	uploadID := c.Request.FormValue("uploadid")
	blockcount := c.Request.FormValue("index")
	blockhash := c.Request.FormValue("chkhash")
	//获取连接池
	conn := redisPool.GetRedisPool().Get()
	defer conn.Close()

	// 获取文件句柄，用于存储分块内容[在本地也创建起来]
	// 创建文件先要查找有没有该文件夹，不然会报错
	filepath := BlockDir + uploadID + "/" + blockcount
	os.MkdirAll(path.Dir(filepath), 0744)

	file, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "upload Multipart create file failed",
			"data": nil,
		})
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*1024) // 1MB
	for {
		n, err := c.Request.Body.Read(buf)
		file.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 对比下hash，配置正确才允许下一步
	cmpSha1, err := util.ComputeSha1ByShell(filepath)
	if err != nil || cmpSha1 != blockhash {
		fmt.Printf("Verify chunk sha1 failed, compare OK: %t, err:%+v\n",
			cmpSha1 == blockhash, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -2,
			"msg":  "Verify hash failed, chkIdx:" + blockcount,
			"data": nil,
		})
		return
	}

	// 更新redis状态
	conn.Do("HSET", BlockKeyPrefix+uploadID, "block_"+blockcount, 1)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": nil,
	})
	// 返回结果
}

// 实现上传分块合并
func CompleteMultipartUpload(c *gin.Context) {
	uploadId := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileSize, _ := strconv.ParseInt(c.Request.FormValue("filesize"), 10, 64)
	fileName := c.Request.FormValue("filename")

	// 获取连接
	conn := redisPool.GetRedisPool().Get()
	defer conn.Close()
	// 判断redis粉饼是否已经完成[设置两个变量进行对比是否已经完成]
	blockTotalCount := 0
	blockCount := 0
	// redis从key值获取数据
	data, err := redis.Values(conn.Do("HGETALL", BlockKeyPrefix+uploadId)) // 获取redis uploadid所有数据
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "complete upload failed",
			"data": nil,
		})
		return
	}
	for i := 0; i < len(data); i += 2 { // 这里为什么是加2？因为额每个data有key与value，所以下一层级是要+2
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "blockcount" {
			blockTotalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "block_") && v == "1" {
			blockCount++ // 将分块的实际上传的数量，在redis中消息获取出来查看
		}
	}
	// 分块上传出现问题
	if blockTotalCount != blockCount {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -2,
			"msg":  "invalid request",
			"data": nil,
		})
		return
	}
	// TODO: 6. 进行文件合并
	if suc := util.MergeBlocksByShell(BlockDir+uploadId, MergeDir+fileHash, fileHash); !suc {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -3,
			"msg":  "Complete upload failed",
			"data": nil,
		})
		return
	}

	// 加入到文件表与用户文件表数据库中
	_ = db.AddFileInfoTodb(fileHash, fileName, fileSize, MergeDir+fileHash)
	_ = db.OnUserFileUploadFinshedDB(username, fileName, fileHash, fileSize)

	// TODO: 6. 并且删除分块文件与redis数据库的分块文件
	_, delHashErr := conn.Do("DEL", HashUpIDKeyPrefix+fileHash)
	delUploadId, delUploadInfoErr := redis.Int64(conn.Do("DEL", BlockKeyPrefix+uploadId))
	if delHashErr != nil || delUploadInfoErr != nil || delUploadId != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -2,
			"msg":  "Complete upload part delete redis data failed",
			"data": nil,
		})
	}

	if suc := util.RemovePathByShell(BlockDir + uploadId); !suc {
		fmt.Printf("Failed to delete chunks as upload canceled, uploadID:%s\n", uploadId)
	}

	//返回请求
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": nil,
	})
	return
}

// 文件取消上传接口
func CancelUpload(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")

	// 2.获取连接池
	conn := redisPool.GetRedisPool().Get()
	defer conn.Close()

	// 3. 检查id是否存在，如果存在删除
	uploadId, err := redis.String(conn.Do("GET", HashUpIDKeyPrefix+filehash))
	if err != nil || uploadId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Cancel upload part failed",
			"data": nil,
		})
		return
	}

	// 4. 删除redis的upload与hash
	_, delHashErr := conn.Do("DEL", HashUpIDKeyPrefix+filehash)
	_, delUploadInfoErr := conn.Do("DEL", BlockKeyPrefix+uploadId)

	if delHashErr != nil || delUploadInfoErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Cancel upload part failed",
			"data": nil,
		})
		return
	}

	// 5. 删除上传的分块文件
	suc := util.RemovePathByShell(BlockKeyPrefix + uploadId)
	if !suc {
		fmt.Printf("Failed to delete chunks as upload canceled, uploadID:%s\n", uploadId)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": nil,
	})
	return
}
