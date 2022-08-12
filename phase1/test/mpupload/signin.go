package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	jsonit "github.com/json-iterator/go"
)

// 外部(命令行)参数

// 测试场景: 正常上传/取消上传/断点续传
var tcase string

// 当前测试只上传的分片数量, 默认为全部上传
var uploadChunkCount int

// 当前用于测试上传的本地文件路径
var uploadFilePath string

// 测试用户名
var username string

// 测试用户密码
var password string

const (
	apiHost           = "http://localhost:8080/"
	apiUserSignin     = apiHost + "user/signin"
	apiUploadInit     = apiHost + "file/mpupload/init"
	apiUploadPart     = apiHost + "file/mpupload/uppart"
	apiUploadComplete = apiHost + "file/mpupload/complete"
	apiUploadCancel   = apiHost + "file/mpupload/cancel"
)

// MultipartUploadInfo : 初始化的分片信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
	// 已经存在的分块，告诉客户端可以跳过这些分块，无需重复上传
	ChunkExists []int
}

// UploadInitResponse : 初始化接口返回的数据
type UploadInitResponse struct {
	Code int                 `json:code`
	Msg  string              `json:msg`
	Data MultipartUploadInfo `json:data`
}

// 登录
func signin(username, password string) (token string, err error) {
	resp, err := http.PostForm(
		apiUserSignin,
		url.Values{
			"username": {username},
			"password": {password},
		})
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || string(body) == "FAILED" {
		return "", err
	}
	token = jsonit.Get(body, "data").Get("Token").ToString()
	fmt.Println("signin token: " + token)
	return token, nil
}
