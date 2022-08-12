package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, `<Test multipart upload>

Usage: ./test_upload [-tcase normal/cancel/resume] [-chkcnt chunkCount] [-fpath filePath] [-user username] [-pwd password]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	// 解析外部参数
	flag.StringVar(&tcase, "tcase", "normal", "测试场景: normal/cancel/resume")
	flag.IntVar(&uploadChunkCount, "chkcnt", 0, "指定当次上传的分块数, 值<=0时表示上传所有分块")
	flag.StringVar(&uploadFilePath, "fpath", "", "指定上传的文件路径")
	flag.StringVar(&username, "user", "", "测试用户名")
	flag.StringVar(&password, "pwd", "", "测试用户密码")
	flag.Parse()

	switch tcase {
	case "normal":
		normalTestMain()
		break
	case "cancel":
		cancelTestMain()
		break
	case "resume":
		resumableTestMain()
		break
	default:
		usage()
	}
}
