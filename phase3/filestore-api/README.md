## 配置 Ceph (可参考 charter7, 不用 Ceph 则可跳过)

## 配置 OSS 连接参数 (config/oss.go)

```go

const (
	// OSSBucket : oss bucket名
	OSSBucket = "buckettest-filestore2"
	// OSSEndpoint : oss endpoint
	OSSEndpoint = "oss-cn-shenzhen.aliyuncs.com"
	// OSSAccesskeyID : oss访问key
	OSSAccesskeyID = "<你的AccesskeyId>"
	// OSSAccessKeySecret : oss访问key secret
	OSSAccessKeySecret = "<你的AccessKeySecret>"
)
```

## 管理依赖包 (go modules)

```bash
go mod init summercamp-filestore
```

## 应用启动

```bash
# 方式1: go build
go build main.go
./main
# 方式2: go run
go run main.go
```
