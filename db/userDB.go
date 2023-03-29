package db

import (
	"fileserver/db/mysql"
	"fmt"
	"log"
)

// 用户注册
func UserSignUp(username string, password string, phone string, email string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`,`phone`,`email`)values(?,?,?,?)")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	defer stmt.Close()
	r, err := stmt.Exec(username, password, phone, email)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	if r, err := r.RowsAffected(); err == nil && r > 0 {
		return true
	}
	return false
}

// 用户登录
func UserLogin(username string, password string) bool {
	row := mysql.DBConn().QueryRow("SELECT user_pwd FROM tbl_user WHERE user_name = ?", username)
	var pwd string
	err := row.Scan(&pwd)
	if err != nil {
		log.Println(err)
		return false
	}
	if pwd == password {
		return true
	}
	return false

}
