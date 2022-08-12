package handler

import (
	cmn "summercamp-filestore/common"
	cfg "summercamp-filestore/config"
	dblayer "summercamp-filestore/db"
	"summercamp-filestore/meta"
	"summercamp-filestore/store/ceph"
	"summercamp-filestore/store/oss"
	"summercamp-filestore/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(c *gin.Context) {
	//c.Redirect(http.StatusFound, "http://"+c.Request.Host+"/static/view/index.html")
	// 返回上传html页面
	data, err := ioutil.ReadFile("./static/view/index.html")
	if err != nil {
		c.Writer.WriteString("internel server error")
		return
	}
	c.Writer.WriteString(string(data))
	//io.WriteString(w, string(data))
}

// UploadHandler ： 处理文件上传
func DoUploadHandler(c *gin.Context) {
	// 接收文件流及存储到本地目录
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err:%s\n", err.Error())
		return
	}
	defer file.Close()

	tmpPath := cfg.TempLocalRootDir + head.Filename
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: tmpPath,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	fmt.Println("test for")
	fmt.Println(tmpPath)

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	// 5. 同步或异步将文件转移到Ceph/OSS
	newFile.Seek(0, 0) // 游标重新回到文件头部
	mergePath := cfg.MergeLocalRootDir + fileMeta.FileSha1
	fmt.Println(mergePath)
	if cfg.CurrentStoreType == cmn.StoreCeph {
		// 文件写入Ceph存储
		data, _ := ioutil.ReadAll(newFile)
		cephPath := "/ceph/" + fileMeta.FileSha1
		err = ceph.PutObject("userfile", cephPath, data)
		if err != nil {
			fmt.Println("upload ceph err: " + err.Error())
			c.Writer.Write([]byte("Upload failed!"))
			return
		}
		fileMeta.Location = cephPath
	} else if cfg.CurrentStoreType == cmn.StoreOSS {
		// 文件写入OSS存储
		ossPath := "oss/" + fileMeta.FileSha1
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			fmt.Println("upload oss err: " + err.Error())
			c.Writer.Write([]byte("Upload failed!"))
			return
		}
		fileMeta.Location = ossPath
	} else {
		fileMeta.Location = mergePath
	}
	// (普通上传/分块上传)本地文件统一存储到mergePath
	err = os.Rename(tmpPath, mergePath)
	if err != nil {
		fmt.Println("move local file err: " + err.Error())
		c.Writer.Write([]byte("Upload failed!"))
		return
	}

	// TODO: 处理异常情况，比如跳转到一个上传失败页面
	_ = meta.UpdateFileMetaDB(fileMeta)

	username := c.Request.FormValue("username")
	suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
		fileMeta.FileName, fileMeta.FileSize)
	fmt.Println("suc")
	fmt.Println(suc)
	if suc {
		c.Redirect(http.StatusFound, "/static/view/home.html")
	} else {
		c.Writer.Write([]byte("Upload Failed."))
	}
}

// UploadSucHandler : 上传已完成
func UploadSucHandler(c *gin.Context) {
	c.Writer.WriteString("Upload finished!")
}

// TryFastUploadHandler : 尝试秒传接口
func TryFastUploadHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta.FileSha1 == "" {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Writer.Write(resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	suc := dblayer.OnUserFileUploadFinished(
		username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Writer.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	c.Writer.Write(resp.JSONBytes())
	return
}
