package db

import (
	"fileserver/db/mysql"
	"fmt"
	"log"
)

// 用户注册
func UserSignUp(username string, password string, phone string, email string) bool {
	fmt.Printf("username=%v\npassword=%v\nphone=%v\nemail=%v\n", username, password, phone, email)
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
	if i, err := r.RowsAffected(); err == nil && i > 0 {
		return true
	} else {
		fmt.Printf("i: %v\n", i)
		fmt.Printf("err: %v\n", err)
		return false
	}

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

// 更新用户token
func UpdateToken(username string, token string) bool {
	// 连接到数据库，根据用户名创建用户或更新用户 Token。
	stmt, err := mysql.DBConn().Prepare(`
		insert into tbl_user_token (user_name, user_token)
		values(?, ?)
		on duplicate key update user_token = values(user_token)
	`)
	if err != nil { // 如果有错误，则打印出来，并返回 false。
		fmt.Printf("err: %v\n", err)
		return false
	}
	defer stmt.Close()                      // 最后关闭 SQL 语句的连接。
	res, err := stmt.Exec(username, token)  // 执行 SQL 语句，如果有错误，则打印出来，并返回 false。
	rowsAffected, err := res.RowsAffected() // 获取受影响的行数
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	if rowsAffected == 0 { // 如果没有受影响的行数，说明没有更新数据，就需要创建新记录
		stmt, err = mysql.DBConn().Prepare(`
			insert into tbl_user_token (user_name,user_token)
			values (?,?)
		`)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return false
		}
		defer stmt.Close()
		_, err = stmt.Exec(username, token)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return false
		}
	}

	return true // 最后返回 true
}

// 查询用户信息
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare(
		"SELECT user_name, signup_at FROM tbl_user WHERE user_name=? LIMIT 1")
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return user, err
	}
	return user, nil

}
