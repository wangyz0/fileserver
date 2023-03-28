package main

import (
	"fileserver/handler"
	"fmt"
	"net/http"
)

func main() {
	// 处理static目录下的静态文件路径
	http.Handle("/accets/", http.StripPrefix("/accets/", http.FileServer(http.Dir("static/accets"))))

	// 注册 upload 页面路由处理函数
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/sign", handler.Sign)

	// 启动服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("fail to start server")
	}
}
