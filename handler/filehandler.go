package handler

import (
	"encoding/json"
	"fileserver/meta"
	"fileserver/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
		fmt.Println("收到post 开始存储") // 输出调试信息

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

		newFile.Seek(0, 0)                      // 文件指针归零
		fmeta.FileSha1 = util.FileSha1(newFile) // 通过计算哈希值来获取文件的唯一标识

		// 将文件属性存储到元数据库中
		meta.UpdataFileMetaDB(fmeta)

		// 跳转上传成功页面。
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
	//解析请求参数
	r.ParseForm()
	fas1 := r.Form.Get("filehash")  //获取文件的Hash信息
	fmeta := meta.GetFileMeta(fas1) //根据Hash查询Meta信息
	// fmt.Println(fmeta)
	// f, err := os.Open(fmeta.Location)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// defer f.Close()
	// data, err := ioutil.ReadAll(f) //读取文件内容到byte数组
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// w.Header().Set("Content-Type", "application/octect-stream") //设置header中的Content-Type
	// // "application/octect-stream"，它表示返回的内容是二进制流。
	// w.Header().Set("Content-Disposition", "attachment;filename=\""+fmeta.FileName+"\"") //设置header中的Content-Descrption
	// // Content-Disposion 用于告诉浏览器以下载的形式展示而非在浏览器中展示。
	// // attachment;filename属性可以指定下载后保存的文件名。
	// // 由于服务器默认输出都是inline（在当前页面开启显示）,attachment则代表以附件形式下载
	// // w.Write(data) //将byte数据作为响应内容返回给客户端
	// _, err = io.Copy(w, file)
	// if err != nil {
	// 	// 处理错误
	// 	return
	// }
	// 下面是gpt的方法
	// 3. 打开文件并返回给客户端
	file, err := os.Open(fmeta.Location)
	if err != nil {
		// 处理错误
		return
	}
	defer file.Close()

	// 设置header
	w.Header().Set("Content-Type", "application/octet-stream") //设置header中的Content-Type
	// "application/octect-stream"，它表示返回的内容是二进制流。
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fmeta.FileName+"\"") //在Content-Disposition中指定下载的文件名
	// Content-Disposion 用于告诉浏览器以下载的形式展示而非在浏览器中展示。
	// attachment;filename属性可以指定下载后保存的文件名。
	// 由于服务器默认输出都是inline（在当前页面开启显示）,attachment则代表以附件形式下载
	_, err = io.Copy(w, file)
	if err != nil {
		// 处理错误
		return
	}
}
