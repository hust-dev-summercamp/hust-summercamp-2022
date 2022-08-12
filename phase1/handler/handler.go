package handler

import (
	"encoding/json"
	cfg "summercamp-filestore/config"
	dblayer "summercamp-filestore/db"
	"summercamp-filestore/meta"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func init() {
	if err := os.MkdirAll(cfg.TempLocalRootDir, 0744); err != nil {
		fmt.Println("无法指定目录用于存储临时文件: " + cfg.TempLocalRootDir)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.MergeLocalRootDir, 0744); err != nil {
		fmt.Println("无法指定目录用于存储合并后文件: " + cfg.MergeLocalRootDir)
		os.Exit(1)
	}
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {
	filehash := c.Request.Form["filehash"][0]
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMeta)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(data)
}

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.Form.Get("limit"))
	username := c.Request.Form.Get("username")
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(data)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.Request.Form.Get("op")
	fileSha1 := c.Request.Form.Get("filehash")
	username := c.Request.Form.Get("username")
	newFileName := c.Request.Form.Get("filename")

	if opType != "0" || len(newFileName) < 1 {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
	if c.Request.Method != "POST" {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 更新用户文件表tbl_user_file中的文件名，tbl_file的文件名不用修改
	_ = dblayer.RenameFileName(username, fileSha1, newFileName)

	// 返回最新的文件信息
	userFile, err := dblayer.QueryUserFileMeta(username, fileSha1)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFile)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(c *gin.Context) {
	username := c.Request.Form.Get("username")
	fileSha1 := c.Request.Form.Get("filehash")

	// 删除本地文件
	fm, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fm.Location)

	// 删除文件表中的一条记录
	suc := dblayer.DeleteUserFile(username, fileSha1)
	if !suc {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
