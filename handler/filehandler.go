package handler

import (
	"encoding/json"
	"fileserver/config"
	"fileserver/db"
	"fileserver/meta"
	"fileserver/mq"
	"fileserver/store/ceph"
	"fileserver/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"fileserver/store/oss"
	// "gopkg.in/amz.v1/s3"
)

// // 注册
// func Sign(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		//返回上传页面
// 		data, err := ioutil.ReadFile("./static/view/sign.html")
// 		if err != nil {
// 			io.WriteString(w, "internel server error")
// 			return
// 		} else {
// 			io.WriteString(w, string(data))
// 		}

// 	} else if r.Method == "POST" {

// 	}
// }

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 如果请求方式为 GET，则返回页面
	if r.Method == "GET" {
		fmt.Println("GET获取用户名")
		username := r.URL.Query().Get("username")
		fmt.Println(username)
		// 读取上传页面
		// 如果读取失败，返回服务器错误状态码；如果成功，则将页面写入 ResponseWriter 中
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		} else {
			// 将读取到的内容写入 ResponseWriter 中。
			io.WriteString(w, string(data))
		}

	} else if r.Method == "POST" {
		username := r.FormValue("username")
		userhash := r.FormValue("filehash") //实现秒传 暂时不用
		fmt.Printf("userhash: %v\n", userhash)
		//接受文件并存储到本地
		file, head, err := r.FormFile("file") // 解析上传文件数据
		if err != nil {
			fmt.Println("failed to get data", err.Error())
			return
		}
		defer file.Close() // 函数结束之前关闭文件链接

		// 构建 FileMeta 结构体，描述文件的属性
		fmeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		fmt.Println("新建文件")                                 // 输出调试信息
		newFile, err := os.Create("./tmp/" + head.Filename) // 在本地创建一个相同名字的文件，用于保存上传的文件数据。
		if err != nil {
			fmt.Println("faild to create file", err.Error())
		}
		defer newFile.Close() // 函数结束之前关闭新建的本地文件链接。

		fmt.Println("复制文件")                          // 输出调试信息
		fmeta.FileSize, err = io.Copy(newFile, file) // 从上传的文件中读取数据，并将数据写入到在本地创建的文件中。同时会返回已经拷贝了多少字节。
		if err != nil {
			fmt.Println("faild to save file", err.Error())
			return
		}
		// 计算hash值
		newFile.Seek(0, 0)                      // 文件指针归零
		fmeta.FileSha1 = util.FileSha1(newFile) // 通过计算哈希值来获取文件的唯一标识

		// 同时将文件写入ceph存储
		// newFile.Seek(0,0)
		// data,_:=ioutil.ReadAll(newFile)
		// bucket:= ceph.GetCephBucket("userfile")
		// cephPath:="/ceph/"+fmeta.FileSha1
		// bucket.Put(cephPath,data,"octet-stream",s3.PublicRead)
		// fmeta.Location=cephPath //路径改为ceph里的路径

		// 将文件属性存储到元数据库中文件表里
		meta.UpdataFileMetaDB(fmeta)

		// 将用户文件属性插入到用户文件表
		suc := db.OnUserFileUploadFinished(username, fmeta.FileSha1, fmeta.FileName, fmeta.FileSize)
		// 跳转上传成功页面。
		if !suc {
			w.Write([]byte("注册失败"))
			return
		}
		fmt.Println("弹出成功提示框")
		http.Redirect(w, r, "/", 200) // 将用户重定向到首页，前端展示上传成功的提示信息。
	}

}

// 上传文件到阿里云OSS
// 处理文件上传
func OssUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 如果请求方式为 GET，则返回页面
	if r.Method == "GET" {
		fmt.Println("GET获取用户名")
		username := r.URL.Query().Get("username")
		fmt.Println(username)
		// 读取上传页面
		// 如果读取失败，返回服务器错误状态码；如果成功，则将页面写入 ResponseWriter 中
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		} else {
			// 将读取到的内容写入 ResponseWriter 中。
			io.WriteString(w, string(data))
		}

	} else if r.Method == "POST" {
		username := r.FormValue("username")
		userhash := r.FormValue("filehash") //实现秒传 暂时不用
		fmt.Printf("userhash: %v\n", userhash)
		//接受文件并存储到本地
		file, head, err := r.FormFile("file") // 解析上传文件数据
		if err != nil {
			fmt.Println("failed to get data", err.Error())
			return
		}
		defer file.Close() // 函数结束之前关闭文件链接

		// 构建 FileMeta 结构体，描述文件的属性
		fmeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		fmt.Println("新建文件")                                 // 输出调试信息
		newFile, err := os.Create("./tmp/" + head.Filename) // 在本地创建一个相同名字的文件，用于保存上传的文件数据。
		if err != nil {
			fmt.Println("faild to create file", err.Error())
		}
		defer newFile.Close() // 函数结束之前关闭新建的本地文件链接。

		fmt.Println("复制文件")                          // 输出调试信息
		fmeta.FileSize, err = io.Copy(newFile, file) // 从上传的文件中读取数据，并将数据写入到在本地创建的文件中。同时会返回已经拷贝了多少字节。
		if err != nil {
			fmt.Println("faild to save file", err.Error())
			return
		}
		// 计算hash值
		newFile.Seek(0, 0)                      // 文件指针归零
		fmeta.FileSha1 = util.FileSha1(newFile) // 通过计算哈希值来获取文件的唯一标识
		// 将文件属性存储到元数据库中文件表里
		meta.UpdataFileMetaDB(fmeta)

		// 同时将文件上传到阿里云
		newFile.Seek(0, 0)
		ossPath := "oss/" + fmeta.FileSha1
		// err = oss.Bucket().PutObject(ossPath, newFile)
		// if err != nil {
		// 	fmt.Printf("err: %v\n", err)
		// 	w.Write([]byte("upload oss faild"))
		// 	return
		// }
		// fmeta.Location = ossPath //路径改成阿里云上
		// 改为使用队列的方式上传
		data := mq.TransferData{
			FileHash:      fmeta.FileSha1,
			CurLocation:   fmeta.Location,
			DestLocation:  ossPath,
			DestStoreType: "oss",
		}
		pubData, err := json.Marshal(data)
		suc := mq.Publish(config.TransExchangeName, config.TransOSSRoutingKey, pubData)
		if !suc {
			fmt.Println("rabbitmq消息发送失败")
		}
		// 将用户文件属性插入到用户文件表
		suc = db.OnUserFileUploadFinished(username, fmeta.FileSha1, fmeta.FileName, fmeta.FileSize)
		// 跳转上传成功页面。
		if !suc {
			w.Write([]byte("注册失败"))
			return
		}
		fmt.Println("弹出成功提示框")
		http.Redirect(w, r, "/", 200) // 将用户重定向到首页，前端展示上传成功的提示信息。
	}

}

// 查询文件信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fmeta)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	w.Write(data)
}

// 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("收到下载请求")
	//解析请求参数
	r.ParseForm()
	fas1 := r.Form.Get("filehash")  //获取文件的Hash信息
	fmeta := meta.GetFileMeta(fas1) //根据Hash查询Meta信息
	// 3. 打开文件并返回给客户端
	file, err := os.Open(fmeta.Location)
	if err != nil {
		// 处理错误
		return
	}
	defer file.Close()

	// 设置header
	fmt.Println("设置header")
	w.Header().Set("Content-Type", "application/octet-stream") //设置header中的Content-Type
	// "application/octect-stream"，它表示返回的内容是二进制流。
	fmt.Println(fmeta.FileName)
	filename := fmeta.FileName
	filename = url.QueryEscape(filename) // 对文件名进行URL编码
	filename = path.Base(filename)       // 只获取文件名部分  因为中文乱码  要前端解码
	w.Header().Set("Content-Disposition", "attachment;filename=\""+filename+"\"")
	// w.Header().Set("Content-Disposition", "attachment;filename=\""+fmeta.FileName+"\"") //在Content-Disposition中指定下载的文件名
	// Content-Disposion 用于告诉浏览器以下载的形式展示而非在浏览器中展示。
	// attachment;filename属性可以指定下载后保存的文件名。
	// 由于服务器默认输出都是inline（在当前页面开启显示）,attachment则代表以附件形式下载
	fmt.Println("返回文件")
	_, err = io.Copy(w, file)
	if err != nil {
		// 处理错误
		return
	}
}

// 从阿里云直接下载
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	row, _ := db.GetFileMeta(filehash)
	singedURL := oss.DownloadURL(row.Location)
	w.Write([]byte(singedURL))
}

// 从ceph下载文件
func DownloadCephHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("收到下载请求")
	//解析请求参数
	r.ParseForm()
	fas1 := r.Form.Get("filehash")
	//获取文件的Hash信息
	fmeta := meta.GetFileMeta(fas1)
	//根据Hash查询Meta信息
	// 打开文件并返回给客户端
	bucket := ceph.GetCephBucket("userfile")
	data, _ := bucket.Get("/ceph/" + fas1)
	// 设置header
	fmt.Println("设置header")
	w.Header().Set("Content-Type", "application/octet-stream")
	//设置header中的Content-Type
	//"application/octect-stream"，它表示返回的内容是二进制流。
	fmt.Println(fmeta.FileName)
	filename := fmeta.FileName
	filename = url.QueryEscape(filename)
	// 对文件名进行URL编码
	filename = path.Base(filename)
	// 只获取文件名部分 因为中文乱码 要前端解码
	w.Header().Set("Content-Disposition", "attachment;filename=\""+filename+"\"")
	//在Content-Disposition中指定下载的文件名
	// Content-Disposion 用于告诉浏览器以下载的形式展示而非在浏览器中展示。
	// attachment;filename属性可以指定下载后保存的文件名。
	// 由于服务器默认输出都是inline（在当前页面开启显示）,attachment则代表以附件形式下载
	fmt.Println("返回文件")
	// 写入ResponseWriter
	_, err := w.Write(data)
	if err != nil {
		// 处理错误
		w.Write([]byte("下载失败"))
	}
}
