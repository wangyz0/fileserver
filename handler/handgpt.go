package handler

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		fmt.Println("Received POST request")

// 		// 解析POST请求中表单键值对数据，包括上传的文件
// 		err := r.ParseMultipartForm(32 << 20) // 最大文件大小设为32MB
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}

// 		// 获取第一个指定名称的文件，本例中默认"name"为"file"
// 		file, handler, err := r.FormFile("file")
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Bad Request", http.StatusBadRequest)
// 			return
// 		}
// 		defer file.Close()

// 		// 创建一个新文件，并将上传的内容复制到它里面
// 		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}
// 		defer f.Close()
// 		io.Copy(f, file)

// 		// 重定向到一个 "Upload Successful" 页面
// 		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
// 	} else {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 	}
// }

// // 在主函数中设置HTTP路由
// func main() {
// 	http.HandleFunc("/file/upload", uploadFileHandler)
// 	http.ListenAndServe(":8000", nil)
// }
