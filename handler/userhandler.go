package handler

import (
	"encoding/json"
	"fileserver/db"
	"fileserver/util"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	// "strings"
	"time"
)

// 注册和登录
const (
	pwd_salt = "*#890"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/sign.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	} else if r.Method == "POST" {
		fmt.Println("收到post")
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println(username)
		password := util.Sha1([]byte(r.Form.Get("password") + pwd_salt)) //密码加密
		phone := r.Form.Get("phone")
		email := r.Form.Get("email")
		fmt.Println(username, password, phone, email)
		if db.UserSignUp(username, password, phone, email) {
			// 注册成功
			fmt.Fprintf(w, `
				<script>
					alert("恭喜您，注册成功！");
					window.location.href="/login";
				</script>
			`)

		} else {
			// 注册失败
			fmt.Fprintf(w, `
				<script>
					alert("注册失败！请重新注册");
					window.location.href="/sign";
				</script>
			`)

		}

	}
}

// 登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/login.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	} else if r.Method == "POST" {
		fmt.Println("收到post")
		r.ParseForm()
		username := r.Form.Get("username")
		password := util.Sha1([]byte(r.Form.Get("password") + pwd_salt)) //密码加密
		if !db.UserLogin(username, password) {                           //验证密码
			w.Write([]byte("登录失败"))
		}
		// 登陆成功
		// 生成访问凭证（token）
		token := GenToken(username)
		b := db.UpdateToken(username, token)
		if b == false {
			w.Write([]byte("写入token失败 登录失败"))
		}
		fmt.Println("开始跳转")
		//登录成功后定向到首页
		resp := util.RespMsg{
			Code: 0,
			Msg:  "OK",
			Data: struct {
				Username string `json:"username"`
				Token    string `json:"token"`
			}{
				Username: username,
				Token:    token,
			},
		}
		// fmt.Printf("Location: http://%s/static/view/home1.html", r.Host)
		w.Write(resp.JSONBytes())
	}
}

// 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	//验证token
	if !IsTokenVaild(token) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//组装并相应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// 生成token
func GenToken(username string) string {
	//40位
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// 验证token
func IsTokenVaild(token string) bool {
	//判断token失效性  后八位是时间
	// 查询token
	// 对比token
	return true
}

// 用户首页
// 声明一个 struct 作为 HomePage 模板中要用到的数据类型
type HomePageData struct {
	Username     string
	RegDate      string
	CapacityUsed string
	CapacityMax  string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("获取username")
		username := r.URL.Query().Get("username")
		fmt.Printf("username: %v\n", username)
		user, err := db.GetUserInfo(username)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		files, err := db.QueryUserFileMetas(username, 10)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Username": username,
			"RegDate":  user.SignupAt[:10],
			"Capacity": "20GB / 100GB",
			"Files":    files,
		}

		t, err := template.ParseFiles("./static/view/home1.html")
		if err != nil {
			log.Printf("Parse home page failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, data); err != nil {
			log.Printf("Render home page failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	} else if r.Method == "POST" {
		fmt.Println("收到homePOST")
		var deleteReq DeleteRequest
		err := json.NewDecoder(r.Body).Decode(&deleteReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Access the username and sha1 values from the DeleteRequest struct
		username := deleteReq.Username
		status := deleteReq.Status
		sha1 := deleteReq.Sha1
		fmt.Printf("fileSha1: %v\n", sha1)
		fmt.Println(username, status)
		if status == 0 { //删除文件
			db.DeleteUserFile(username, sha1)
		} else if status == 1 { //下载文件
			fmt.Println("进行下载")
			http.Redirect(w, r, fmt.Sprintf("/download?filehash=%s", sha1), http.StatusFound)
		}

	}

	// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// 用来接受home发来的post的结构体
type DeleteRequest struct {
	Username string `json:"username"`
	Sha1     string `json:"sha1"`
	Status   int    `json:"status"`
}
