package handler

import (
	"fileserver/db"
	"fileserver/util"
	"fmt"
	"io/ioutil"
	"net/http"
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
		if db.UserLogin(username, password) {
			w.Write([]byte("登录成功"))
		} else {
			w.Write([]byte("登录失败"))
		}
	}
}
