package main

import (
	"fileserver/handler"
	"fmt"
	"net/http"
)

func main() {
	// 处理static目录下的静态文件路径
	http.Handle("/accets/", http.StripPrefix("/accets/", http.FileServer(http.Dir("static/accets"))))

	//  uploads 上传页面路由处理函数
	http.HandleFunc("/file/upload", handler.UploadHandler)
	// login登录页面上传
	http.HandleFunc("/login", handler.LoginHandler)
	// sign注册页面
	http.HandleFunc("/sign", handler.SignupHandler)
	// 查询文件接口
	http.HandleFunc("/query", handler.GetFileMetaHandler)
	//文件下载接口
	http.HandleFunc("/download", handler.DownloadHandler)
	// 用户信息查询
	http.HandleFunc("/user/info", handler.UserInfoHandler)
	// // 主页
	// http.HandleFunc("/home", handler.HomeHandeler)

	// 启动服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("fail to start server")
	}

}
