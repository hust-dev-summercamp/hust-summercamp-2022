package main

import (
	"summercamp-filestore/store/ceph"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/amz.v1/s3"
)

func main() {
	bucket := ceph.GetCephBucket("userfile")

	// 创建一个新的bucket
	err := bucket.PutBucket(s3.PublicRead)
	fmt.Printf("create bucket err: %v\n", err)

	// 查询这个bucket下面指定条件的object keys
	res, err := bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("bucket list err: %s\n", err.Error())
	} else {
		fmt.Printf("object keys: %+v\n", res)
	}

	// return

	// 上传字节对象A
	objAPath := "/testupload/a.txt"
	err = bucket.Put(objAPath, []byte("just for test"), "octet-stream", s3.PublicRead)
	fmt.Printf("upload object err: %+v\n", err)

	// 下载字节对象A
	objA, err := bucket.Get(objAPath)
	if err != nil {
		fmt.Printf("Get object A err: %s\n", err.Error())
		return
	}
	fmt.Printf("object A body: %s\n", string(objA))

	// 上传文件B
	localPath := "/data/pkg/armory-20.01.zip"
	objBPath := "/ceph/5a4e4375ac922a7d2d971fe81d7e6356d73f39a1"
	fd, err := os.Open(localPath)
	if err != nil {
		fmt.Printf("Open file failed, %+v\n", err)
		return
	}
	fileBody, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Printf("Read file failed, %+v\n", err)
		return
	}
	err = bucket.Put(objBPath, fileBody, "octet-stream", s3.PublicRead)
	fmt.Printf("upload file err: %+v\n", err)

	// 下载文件B
	objB, err := bucket.Get(objBPath)
	if err != nil {
		fmt.Printf("Get object B err: %s\n", err.Error())
		return
	}
	tmpFile, err := os.Create(localPath + ".copy")
	if err != nil {
		fmt.Printf("Write object B to file err: %s\n", err.Error())
		return
	}
	tmpFile.Write(objB)

	// 查询这个bucket下面指定条件的object keys
	res, err = bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("bucket list err: %s\n", err.Error())
	} else {
		fmt.Printf("object keys: %+v\n", res)
	}
}
