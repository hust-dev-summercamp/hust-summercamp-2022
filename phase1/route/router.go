package route

import (
	"summercamp-filestore/handler"
	"github.com/gin-gonic/gin"
)

// Router ：路由规则定义
func Router() *gin.Engine {
	// gin framework
	router := gin.Default()

	// 静态资源处理
	router.Static("/static/", "./static")

	// 不需验证的接口
	router.GET("/user/signup", handler.SignupHandler)
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	router.POST("/user/signin", handler.DoSignInHandler)
	//router.GET("/user/exists", hdl.UserExistsHandler)

	// 加入auth认证中间件
	router.Use(handler.Authorize())

	// 文件存取接口
	router.GET("/file/upload", handler.UploadHandler)
	router.POST("/file/upload", handler.DoUploadHandler)
	router.GET("/file/upload/suc", handler.UploadSucHandler)
	router.GET("/file/meta", handler.GetFileMetaHandler)
	router.POST("/file/query", handler.FileQueryHandler)
	//router.GET("/file/download", handler.DownloadHandler)
	router.GET("/file/download", handler.DownloadHandler)
	router.GET("/file/download/range", handler.RangeDownloadHandler)
	router.POST("/file/update", handler.FileMetaUpdateHandler)
	router.DELETE("/file/delete", handler.FileDeleteHandler)
	router.POST("/file/downloadurl",
		handler.DownloadURLHandler)

	// 秒传接口
	router.POST("/file/fastupload",
		handler.TryFastUploadHandler)

	// 分块上传接口
	router.POST("/file/mpupload/init",
		handler.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart",
		handler.UploadPartHandler)
	router.POST("/file/mpupload/complete",
		handler.CompleteUploadHandler)
	router.POST("/file/mpupload/cancle",
		handler.CancelUploadHandler)

	// 用户相关接口
	router.POST("/user/info", handler.UserInfoHandler)

	return router
}
