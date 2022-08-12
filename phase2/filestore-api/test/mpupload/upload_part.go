package main

import (
	"bufio"
	"bytes"
	fsUtil "summercamp-filestore/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// 实际上传分块逻辑
// chunkIdxs : 实际需要上传的分块
func uploadPartsSpecified(filename string, targetURL string, chunkSize int, chunkIdxs []int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		// 判断当前所在的块是否需要上传
		if contained, err := fsUtil.Contain(chunkIdxs, index); err != nil || !contained {
			continue
		}

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			resp, err := http.Post(
				targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("%+v %+v\n", string(body), err)
			}
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < len(chunkIdxs); idx++ {
		select {
		case res := <-ch:
			fmt.Printf("完成传输块index: %d\n", res)
		}
	}

	fmt.Printf("全部完成以下分块传输: %+v\n", chunkIdxs)
	return nil
}
