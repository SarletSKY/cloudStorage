package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 静态资源处理
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//http.HandleFunc("/file/upload", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.UploadHandler)))
	//http.HandleFunc("/file/upload/success", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.UploadSucHandler)))
	//http.HandleFunc("/file/meta", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.GetFileMetaInfo)))
	//http.HandleFunc("/file/query", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.GetManyFileMetaInfo)))
	//http.HandleFunc("/file/download", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.DownLoadFile)))
	//http.HandleFunc("/file/download/range", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.RangeDownload)))
	//http.HandleFunc("/file/update", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.UpdateFileInfo)))
	//http.HandleFunc("/file/delete", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.DeleteFile)))
	//http.HandleFunc("/file/downloadurl", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.DownloadURL)))

	//http.HandleFunc("/file/fastupload", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.FastUploadUserFile)))

	//http.HandleFunc("/file/mpupload/init", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.InitMultipartUpload)))
	//http.HandleFunc("/file/mpupload/upload", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.MultipartUpload)))
	//http.HandleFunc("/file/mpupload/complete", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.CompleteMultipartUpload)))
	//http.HandleFunc("/file/mpupload/delete", handler.ReceiveClientRequest(handler.HTTPInterceptor(handler.CancelUpload)))

	fmt.Println("已启动：8080端口......")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("starting server failed")
	}
}
