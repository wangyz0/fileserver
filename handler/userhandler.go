package handler

import (
	"fileserver/db"
	"fileserver/util"
	"fmt"
	"io/ioutil"
	"net/http"
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
			w.Write([]byte("注册成功"))
		} else {
			w.Write([]byte("注册失败"))
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
		//登录成功后定向到首页
		resp := util.RespMsg{
			Code: 0,
			Msg:  "OK",
			Data: struct {
				Location string
				Username string
				Token    string
			}{
				Location: "http://" + r.Host + "/static/view/home1.html",
				Username: username,
				Token:    token,
			},
		}
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
func HomeHandeler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/home1.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
}
