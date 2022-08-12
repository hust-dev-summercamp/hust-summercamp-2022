在当前目录`test/`下进行测试前，需要先编译:

```bash
go build -o test_upload
```

用法说明:

```bash
<Test multipart upload>

Usage: ./test_upload [-tcase normal/cancel/resume] [-chkcnt chunkCount] [-fpath filePath] [-user username] [-pwd password]

Options:
  -chkcnt int
    	指定当次上传的分块数, 值<=0时表示上传所有分块
  -fpath string
    	指定上传的文件路径
  -pwd string
    	测试用户密码
  -tcase string
    	测试场景: normal/cancel/resume (default "normal")
  -user string
    	测试用户名
```

测试(正常分块上传):

```bash
# 根据实际情况修改参数： -fpath / -user / -pwd
./test_upload -tcase normal -fpath /data/pkg/xxx.zip -user test -pwd 123456
```

测试(分块上传，未完成时取消):

```bash
./test_upload -tcase cancel -fpath /data/pkg/xxx.zip -user test -pwd 123456
```

测试(断点续传，多次暂停/上传模拟续传):

```bash
./test_upload -tcase resume -chkcnt 10 -fpath /data/pkg/xxx.zip -user test -pwd 123456
# 文件合并结果: {"code":-2,"msg":"invalid request","data":null}
./test_upload -tcase resume -chkcnt 10 -fpath /data/pkg/xxx.zip -user test -pwd 123456
# 文件合并结果: {"code":-2,"msg":"invalid request","data":null}
./test_upload -tcase resume -chkcnt 10 -fpath /data/pkg/xxx.zip -user test -pwd 123456
# 文件合并结果: {"code":0,"msg":"OK","data":null} 

# 直到返回结果包含 `"msg":"OK"` 时，说明续传已经完成。
```
