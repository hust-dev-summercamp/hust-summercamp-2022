package main

// 测试分块上传及断点续传

import (
	"encoding/json"
	cmn "summercamp-filestore/common"
	"summercamp-filestore/util"
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

// 完成当次的分块上传 (指定分块)
func uploadFile(username, token, fhash, fpath string, filesize int) string {
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
		return ""
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	chunkCount := jsonit.Get(body, "data").Get("ChunkCount").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d chunkCount: %d\n", uploadID, chunkSize, chunkCount)

	// 3. 请求分块上传接口
	tURL := apiUploadPart + "?username=" + username +
		"&token=" + token + "&uploadid=" + uploadID

	if uploadChunkCount <= 0 {
		uploadPartsSpecified(fpath, tURL, chunkSize, []int{})
	} else {
		var initResp UploadInitResponse
		err = json.Unmarshal(body, &initResp)
		if err != nil {
			fmt.Printf("Parse error: %s\n", err.Error())
			return ""
		}
		var chunksToUpload []int
		for idx := 1; idx <= initResp.Data.ChunkCount; idx++ {
			if len(chunksToUpload) >= uploadChunkCount {
				break
			}
			if contained, _ := util.Contain(initResp.Data.ChunkExists, idx); contained {
				continue
			}
			chunksToUpload = append(chunksToUpload, idx)
		}
		fmt.Printf("将要上传的分块: %+v\n", chunksToUpload)
		uploadPartsSpecified(fpath, tURL, chunkSize, chunksToUpload)
	}

	return uploadID
}

func resumeUsage() {
	fmt.Fprintf(os.Stderr, `<Test resumable multipart upload>

	Usage: ./test_upload -tcase resume [-chkcnt chunkCount] [-fpath filePath] [-user username] [-pwd password]

Options:
`)
	flag.PrintDefaults()
}

func resumableTestMain() {
	// 判断参数是否有效
	if exist, err := fsUtil.PathExists(uploadFilePath); !exist || err != nil {
		fmt.Println("Error: 无效文件路径，请检查")
		resumeUsage()
		return
	} else if len(username) == 0 || len(password) == 0 {
		fmt.Println("Error: 无效用户名/密码，请检查")
		resumeUsage()
		return
	}

	// 需要上传的文件名及文件hash
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

	// 1: 登录，获取token
	token, err := signin(username, password)
	if err != nil {
		fmt.Println(err)
		return
	} else if token == "" {
		fmt.Println("登录失败，请检查用户名、密码")
		return
	}

	// 2: 分块上传(断点续传)
	uploadID := uploadFile(username, token, fhash, uploadFilePath, filesize)
	if uploadID == "" {
		return
	}

	// 3. 请求分块完成接口
	resp, err := http.PostForm(
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
	// 4. 打印分块上传结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("文件合并结果: %s\n", string(body))
}
