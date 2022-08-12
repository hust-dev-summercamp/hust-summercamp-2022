package main

// 测试分块上传 (场景：正常完成上传)

import (
	"encoding/json"
	cmn "summercamp-filestore/common"
	fsUtil "summercamp-filestore/util"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	jsonit "github.com/json-iterator/go"
)

func normalUsage() {
	fmt.Fprintf(os.Stderr, `<Test normal multipart upload>

	Usage: ./test_upload -tcase normal [-fpath filePath] [-user username] [-pwd password]

Options:
`)
	flag.PrintDefaults()
}

func normalTestMain() {
	// 判断参数是否有效
	if exist, err := fsUtil.PathExists(uploadFilePath); !exist || err != nil {
		fmt.Println("Error: 无效文件路径，请检查")
		normalUsage()
		return
	} else if len(username) == 0 || len(password) == 0 {
		fmt.Println("Error: 无效用户名/密码，请检查")
		normalUsage()
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
	fmt.Println("token: " + token)

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
	// fmt.Println("init_api_resp: " + string(body))
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
	var initResp UploadInitResponse
	err = json.Unmarshal(body, &initResp)
	if err != nil {
		fmt.Printf("Parse error: %s\n", err.Error())
		os.Exit(-1)
	}
	var chunksToUpload []int
	for idx := 1; idx <= initResp.Data.ChunkCount; idx++ {
		chunksToUpload = append(chunksToUpload, idx)
	}
	uploadChunkCount = len(chunksToUpload)
	tURL := apiUploadPart + "?username=" + username +
		"&token=" + token + "&uploadid=" + uploadID
	// 上传所有分块
	uploadPartsSpecified(uploadFilePath, tURL, chunkSize, chunksToUpload)

	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		apiUploadComplete,
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {fhash},
			"filesize": {strconv.Itoa(filesize)},
			"filename": {filepath.Base(uploadFilePath)},
			"uploadid": {uploadID},
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
	fmt.Printf("complete result: %s\n", string(body))
}
