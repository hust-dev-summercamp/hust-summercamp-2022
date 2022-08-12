package handler

import (
	"summercamp-filestore/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"

	rPool "summercamp-filestore/cache/redis"
	cmn "summercamp-filestore/common"
	cfg "summercamp-filestore/config"
	dblayer "summercamp-filestore/db"
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
	// 已经上传完成的分块索引列表
	ChunkExists []int
}

const (
	// ChunkDir : 上传的分块所在目录
	ChunkDir = cfg.ChunckLocalRootDir
	// MergeDir : 合并后的文件所在目录
	MergeDir = cfg.MergeLocalRootDir
	// ChunkKeyPrefix : 分块信息对应的redis键前缀
	ChunkKeyPrefix = "MP_"
	// HashUpIDKeyPrefix : 文件hash映射uploadid对应的redis键前缀
	HashUpIDKeyPrefix = "HASH_UPID_"
)

func init() {
	if err := os.MkdirAll(ChunkDir, 0744); err != nil {
		fmt.Println("无法指定目录用于存储分块文件: " + ChunkDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(MergeDir, 0744); err != nil {
		fmt.Println("无法指定目录用于存储合并后文件: " + MergeDir)
		os.Exit(1)
	}
}

// InitialMultipartUploadHandler : 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.Request.Form.Get("username")
	filehash := c.Request.Form.Get("filehash")
	filesize, err := strconv.Atoi(c.Request.Form.Get("filesize"))
	if err != nil {
		c.Writer.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 判断文件是否存在
	if dblayer.IsUserFileUploaded(username, filehash) {
		c.Writer.Write(util.NewRespMsg(int(cmn.FileAlreadExists), "OK", nil).JSONBytes())
		return
	}

	// 2. 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadID := ""

	// 3. 通过文件hash判断是否断点续传，并获取uploadID
	keyExists, _ := redis.Bool(rConn.Do("EXISTS", HashUpIDKeyPrefix+filehash))
	if keyExists {
		uploadID, err = redis.String(rConn.Do("GET", HashUpIDKeyPrefix+filehash))
		if err != nil {
			c.Writer.Write(util.NewRespMsg(-1, "Upload part failed", err.Error()).JSONBytes())
			return
		}
	}

	// 4.1 首次上传则新建uploadID
	// 4.2 断点续传则根据uploadID获取已上传的文件分块列表
	chunksExist := []int{}
	if uploadID == "" {
		uploadID = username + fmt.Sprintf("%x", time.Now().UnixNano())
	} else {
		chunks, err := redis.Values(rConn.Do("HGETALL", ChunkKeyPrefix+uploadID))
		if err != nil {
			c.Writer.Write(util.NewRespMsg(-2, "Upload part failed", err.Error()).JSONBytes())
			return
		}
		for i := 0; i < len(chunks); i += 2 {
			k := string(chunks[i].([]byte))
			v := string(chunks[i+1].([]byte))
			if strings.HasPrefix(k, "chkidx_") && v == "1" {
				// chkidx_6 -> 6
				chunkIdx, _ := strconv.Atoi(k[7:len(k)])
				chunksExist = append(chunksExist, chunkIdx)
			}
		}
	}

	// 5. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:    filehash,
		FileSize:    filesize,
		UploadID:    uploadID,
		ChunkSize:   5 * 1024 * 1024, // 5MB
		ChunkCount:  int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
		ChunkExists: chunksExist,
	}

	// 6. 将初始化信息写入到redis缓存
	if len(upInfo.ChunkExists) <= 0 {
		hkey := ChunkKeyPrefix + upInfo.UploadID
		rConn.Do("HSET", hkey, "chunkcount", upInfo.ChunkCount)
		rConn.Do("HSET", hkey, "filehash", upInfo.FileHash)
		rConn.Do("HSET", hkey, "filesize", upInfo.FileSize)
		rConn.Do("EXPIRE", hkey, 43200)
		rConn.Do("SET", HashUpIDKeyPrefix+filehash, upInfo.UploadID, "EX", 43200)
	}

	// 7. 将响应初始化数据返回到客户端
	c.Writer.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	//	username := r.Form.Get("username")
	uploadID := c.Request.Form.Get("uploadid")
	chunkIndex := c.Request.Form.Get("index")

	// 2. 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := ChunkDir + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.Writer.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4. 更新redis缓存状态
	rConn.Do("HSET", ChunkKeyPrefix+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	c.Writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(c *gin.Context) {
	// 1. 解析请求参数
	uploadID := c.Request.Form.Get("uploadid")
	username := c.Request.Form.Get("username")
	filehash := c.Request.Form.Get("filehash")
	filesize := c.Request.Form.Get("filesize")
	filename := c.Request.Form.Get("filename")

	// 2. 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		c.Writer.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		c.Writer.Write(util.NewRespMsg(-2, "Invalid request", nil).JSONBytes())
		return
	}

	// 4. 合并分块 (备注: 更新于2020/04/01; 此合并逻辑非必须实现，因后期转移到ceph/oss时也可以通过分块方式上传)
	if mergeSuc := util.MergeChuncksByShell(ChunkDir+uploadID, MergeDir+filehash, filehash); !mergeSuc {
		c.Writer.Write(util.NewRespMsg(-3, "Complete upload failed", nil).JSONBytes())
		return
	}

	// 5. 更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	// 更新于2020-04: 增加fileaddr参数的写入
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), MergeDir+filehash)
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	// 更新于2020-04: 删除已上传的分块文件及redis分块信息
	_, delHashErr := rConn.Do("DEL", HashUpIDKeyPrefix+filehash)
	delUploadID, delUploadInfoErr := redis.Int64(rConn.Do("DEL", ChunkKeyPrefix+uploadID))
	if delUploadID != 1 || delUploadInfoErr != nil || delHashErr != nil {
		c.Writer.Write(util.NewRespMsg(-4, "Complete upload part failed", nil).JSONBytes())
		return
	}

	delRes := util.RemovePathByShell(ChunkDir + uploadID)
	if !delRes {
		fmt.Printf("Failed to delete chuncks as upload comoleted, uploadID: %s\n", uploadID)
	}

	// 6. 响应处理结果
	c.Writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CancelUploadHandler : 文件取消上传接口
func CancelUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	filehash := c.Request.Form.Get("filehash")

	// 2. 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 检查uploadID是否存在，如果存在则删除
	uploadID, err := redis.String(rConn.Do("GET", HashUpIDKeyPrefix+filehash))
	if err != nil || uploadID == "" {
		c.Writer.Write(util.NewRespMsg(-1, "Cancel upload part failed", nil).JSONBytes())
		return
	}

	_, delHashErr := rConn.Do("DEL", HashUpIDKeyPrefix+filehash)
	_, delUploadInfoErr := rConn.Do("DEL", ChunkKeyPrefix+uploadID)
	if delHashErr != nil || delUploadInfoErr != nil {
		c.Writer.Write(util.NewRespMsg(-2, "Cancel upload part failed", nil).JSONBytes())
		return
	}

	// 4. 删除已上传的分块文件
	delChkRes := util.RemovePathByShell(ChunkDir + uploadID)
	if !delChkRes {
		fmt.Printf("Failed to delete chunks as upload canceled, uploadID:%s\n", uploadID)
	}

	// 5. 响应客户端
	c.Writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
