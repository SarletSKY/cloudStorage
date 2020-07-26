package mpupload

import (
	"flag"
	"fmt"
	"os"
)

// 提示命令
func usage() {
	fmt.Fprintf(os.Stderr, `<Test multipart upload>

Usage: ./test_upload [-tcase normal/cancel/resume] [-chkcnt chunkCount] [-fpath filePath] [-user username] [-pwd password]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&tCase, "tcase", "normal", "测试场景：normal/cancel/resume")
	flag.IntVar(&uploadBlockCount, "uplblkcnt", 0, "指定当次上传的分块数, 值<=0时表示上传所有分块")
	flag.StringVar(&uploadFilePath, "uplpath", "", "指定上传的文件路径")
	flag.StringVar(&username, "user", "", "测试用户名")
	flag.StringVar(&username, "pwd", "", "测试密码")
	flag.Parse()

	switch tCase {
	case "normal":
		normalTestMain() // 正常上传
		break
	case "cancel":
		cancelTestMain() // 取消上传
		break
	case "resume":
		resumableTestMain() // 断点续传
	default:
		usage()
	}
}
