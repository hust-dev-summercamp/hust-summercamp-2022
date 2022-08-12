package main

import (
	"summercamp-filestore/config"
	"summercamp-filestore/route"
	"fmt"
	"net/http"
)

func main() {
	// 静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	// 动态接口路由设置
	//http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	//http.HandleFunc("/file/upload/suc", handler.HTTPInterceptor(handler.UploadSucHandler))
	//http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	//http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	//http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	//http.HandleFunc("/file/download/range", handler.HTTPInterceptor(handler.RangeDownloadHandler))
	//http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	//http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))
	//http.HandleFunc("/file/downloadurl", handler.HTTPInterceptor(
	//	handler.DownloadURLHandler))

	// 秒传接口
	//http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(
	//	handler.TryFastUploadHandler))

	// 分块上传接口
	//http.HandleFunc("/file/mpupload/init",
	//	handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	//http.HandleFunc("/file/mpupload/uppart",
	//	handler.HTTPInterceptor(handler.UploadPartHandler))
	//http.HandleFunc("/file/mpupload/complete",
	//	handler.HTTPInterceptor(handler.CompleteUploadHandler))
	//http.HandleFunc("/file/mpupload/cancel",
	//	handler.HTTPInterceptor(handler.CancelUploadHandler))

	// 用户相关接口
	// http.HandleFunc("/", handler.SignInHandler)
	//http.HandleFunc("/user/signup", handler.SignupHandler)
	//http.HandleFunc("/user/signin", handler.SignInHandler)
	//http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	router := route.Router()

	fmt.Println("上传服务正在启动, 监听端口:8080...")
	// 启动服务并监听端口
	err := router.Run(config.UploadServiceHost)
	// 监听端口
	//err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
