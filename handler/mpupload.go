package handler

import (
	myredis "fileserver/cache/redis"
	"fileserver/db"
	"fileserver/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// 分块上传结构体
type MultipartUpload struct {
	FileHash   string // 文件哈希值
	FileSize   int64  // 文件大小
	UploadID   string // 上传任务的唯一标识符
	ChunkSize  int    // 分块大小，单位字节
	ChunkCount int    // 分块总数
}

// 初始化分块上传  j将分块信息发送给客户端
func InitiaMultipartUpload(w http.ResponseWriter, r *http.Request) {
	// 1.解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")                    // 用户名
	filehash := r.Form.Get("filehash")                    // 文件哈希值
	filesize, err := strconv.Atoi(r.Form.Get("filesize")) // 获取文件原始大小并转换为int类型
	if err != nil {
		// 参数错误返回json信息
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 2.获得redis链接
	rConn := myredis.RedisPool().Get() // 从 Redis 连接池中获取一个连接
	defer rConn.Close()                // 连接使用完毕后，关闭连接

	// 3.生成分块上传初始化信息
	upinfo := MultipartUpload{
		FileHash:   filehash,                                              // 将文件哈希值放入 upinfo 结构体
		FileSize:   int64(filesize),                                       // 将文件大小转为 int64 型，并加入 upinfo 结构体
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),   // 生成上传 ID，并加入 upinfo 结构体
		ChunkSize:  5 * 1024 * 1024,                                       // 默认分块大小为 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))), // 计算出分块数量并加入 upinfo 结构体
	}

	// 4.将初始化信息写入redis
	rConn.Do("HSET", "MP_"+upinfo.UploadID, "chunkcount", upinfo.ChunkCount) // 将上传ID的分块数量写入 Redis
	rConn.Do("HSET", "MP_"+upinfo.UploadID, "filehash", upinfo.FileHash)     // 将上传ID对应的文件哈希值写入 Redis
	rConn.Do("HSET", "MP_"+upinfo.UploadID, "filesize", upinfo.FileSize)     // 将上传ID对应的文件大小写入 Redis

	// 5.将响应数据返回客户端
	w.Write(util.NewRespMsg(0, "OK", upinfo).JSONBytes()) // 响应 http 请求，同时返回写入 Redis 后的上传信息
}

// 分块上传
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1．解析用户请求参数
	r.ParseForm()
	//username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid") // 上传任务唯一标识符
	chunklndex := r.Form.Get("index")  // 分块编号
	//2．获得redis连接池中的一个连接
	rConn := myredis.RedisPool().Get()
	defer rConn.Close()
	//3．获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunklndex // 分块所在路径
	os.MkdirAll(path.Dir(fpath), 0744)              // 创建分块文件目录
	fd, err := os.Create(fpath)                     // 创建并打开分块文件
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes()) // 返回错误信息到客户端
		return
	}
	defer fd.Close()              // 函数结束前关闭文件描述符
	buf := make([]byte, 10241024) // 缓存区，指定 1MB 大小
	for {
		n, err := r.Body.Read(buf) // 从 request body 中读取数据到 buf 中，返回字节数和错误信息
		fd.Write(buf[:n])          // 将读取到的部分写入到文件
		if err != nil {            // 如果读到末尾或者出错就退出循环
			break
		}
	}
	//4 更新redis缓存状态，标记该分块已经上传
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunklndex, 1)
	//5 返回处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler 通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1．解析请求参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")     // 获取上传ID
	username := r.Form.Get("username") // 获取用户名
	filehash := r.Form.Get("filehash") // 获取文件哈希值
	filesize := r.Form.Get("filesize") // 获取文件大小
	filename := r.Form.Get("filename") // 获取文件名

	//2．获得redis连接池中的一个连接
	rConn := myredis.RedisPool().Get() // 使用 Redis 连接池
	defer rConn.Close()                // 在 defer 中释放连接
	//3． 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid)) // 从 Redis 中取出上传任务信息
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes()) // 如果查询信息失败，返回错误状态码
		return
	}
	totalCount := 0 // 总共要上传的分块数量
	chunkCount := 0 // 已经上传的分块数量
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))   // 将二进制转化为字符串格式
		v := string(data[i+1].([]byte)) // 将二进制转化为字符串格式
		chunkCount := 0                 // 已经上传的分块数量
		if k == "chunkcount" {          // 如果是总分块数信息
			totalCount, _ = strconv.Atoi(v) // 将 v 转化为 int 类型
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount { // 检查分块数量是否和上传文件分块数量相同
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}
	//4．合并分块
	//5．更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	db.OnFileUploadFinished(filehash, filename, int64(fsize), "")           // 将文件信息写入文件表中
	db.OnUserFileUploadFinished(username, filehash, filename, int64(fsize)) // 将文件信息写入用户文件表中

	//6．响应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes()) // 返回成功状态码与空白数据

}
