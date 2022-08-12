package handler

import (
	cfg "summercamp-filestore/config"
	dblayer "summercamp-filestore/db"
	"summercamp-filestore/meta"
	"summercamp-filestore/store/ceph"
	"summercamp-filestore/store/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.Form.Get("filehash")
	// 从文件表查找记录
	row, _ := dblayer.GetFileMeta(filehash)

	fmt.Println("fileAddr: " + row.FileAddr.String)
	// 判断文件存在OSS，还是Ceph，还是在本地
	//if strings.HasPrefix(row.FileAddr.String, cfg.TempLocalRootDir) {
	//	username := c.Request.FormValue("username")
	//	token := c.Request.FormValue("token")
	//	tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
	//		c.Request.Host, filehash, username, token)
	//	c.Data(http.StatusOK, "octet-stream", []byte(tmpUrl))
	//} else
	if strings.HasPrefix(row.FileAddr.String, cfg.MergeLocalRootDir) ||
		strings.HasPrefix(row.FileAddr.String, "/ceph") {
		username := c.Request.Form.Get("username")
		token := c.Request.Form.Get("token")
		tmpURL := fmt.Sprintf(
			"http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host,
			filehash,
			username,
			token)
		c.Writer.Write([]byte(tmpURL))
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// oss下载url
		signedURL := oss.DownloadURL(row.FileAddr.String)
		c.Writer.Write([]byte(signedURL))
	} else {
		c.Writer.Write([]byte("Error: 下载链接暂时无法生成"))
	}
}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.Form.Get("filehash")
	username := c.Request.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	fmt.Println(userFile)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fileData []byte
	if strings.HasPrefix(fm.Location, cfg.MergeLocalRootDir) {
		f, err := os.Open(fm.Location)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fileData, err = ioutil.ReadAll(f)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fm.Location, "/ceph") {
		fmt.Println("to download file from ceph...")
		bucket := ceph.GetCephBucket("userfile")
		fileData, err = bucket.Get(fm.Location)
		if err != nil {
			fmt.Println(err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fm.Location, "oss") {
		fmt.Println("to download file from oss...")
		// TODO: to verify the code in this block
		fd, err := oss.Bucket().GetObject(fm.Location)
		if err == nil {
			fileData, err = ioutil.ReadAll(fd)
		}
		if err != nil {
			fmt.Println(err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		c.Writer.Write([]byte("File not found."))
		return
	}

	c.Writer.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Writer.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	c.Writer.Write(fileData)
}

// RangeDownloadHandler : 支持断点的文件下载接口
func RangeDownloadHandler(c *gin.Context) {
	fsha1 := c.Request.Form.Get("filehash")
	username := c.Request.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	fmt.Println(userFile.FileName)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 使用本地目录文件
	fpath := cfg.MergeLocalRootDir + fm.FileSha1
	fmt.Println("range-download-fpath: " + fpath)

	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	c.Writer.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Writer.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	c.File(fpath)
}
