package main

// 测试分块上传 (场景：上传过程中取消)

import (
	cmn "summercamp-filestore/common"
	fsUtil "summercamp-filestore/util"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	jsonit "github.com/json-iterator/go"
)

func cancelUsage() {
	fmt.Fprintf(os.Stderr, `<Test cancel multipart upload>

	Usage: ./test_upload cancel [-fpath filePath] [-user username] [-pwd password]

Options:
`)
	flag.PrintDefaults()
}

func cancelTestMain() {
	// 判断参数是否有效
	if exist, err := fsUtil.PathExists(uploadFilePath); !exist || err != nil {
		fmt.Println("Error: 无效文件路径，请检查")
		cancelUsage()
		return
	} else if len(username) == 0 || len(password) == 0 {
		fmt.Println("Error: 无效用户名/密码，请检查")
		cancelUsage()
		return
	}

	fhash, err := fsUtil.ComputeSha1ByShell(uploadFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize, err := fsUtil.ComputeFileSizeByShell(uploadFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 0: 登录，获取token
	token, err := signin(username, password)
	if err != nil {
		fmt.Println(err)
		return
	} else if token == "" {
		fmt.Println("登录失败，请检查用户名、密码")
		return
	}

	// 1. 请求初始化分块上传接口
	resp, err := http.PostForm(
		apiUploadInit,
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {fhash},
			"filesize": {strconv.Itoa(filesize)},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	respCode := jsonit.Get(body, "code").ToInt()
	if respCode == int(cmn.FileAlreadExists) {
		fmt.Println("文件已存在，无需上传")
		return
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// 3. 请求分块上传接口
	tURL := apiUploadPart + "?username=" + username +
		"&token=" + token + "&uploadid=" + uploadID
	// 只上传第一个分块后，取消上传
	uploadChunkCount = 1
	uploadPartsSpecified(uploadFilePath, tURL, chunkSize, []int{1})

	// 4. 取消分块上传接口
	resp, err = http.PostForm(
		apiUploadCancel,
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {fhash},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	// 5. 打印分块上传结果
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("cancel upload result: %s\n", string(body))
}
