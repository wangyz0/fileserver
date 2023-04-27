package db

import (
	"fileserver/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName     string
	FileHash     string
	FileName     string
	FileSize     int64
	UploadAt     string
	LastUploadAt string
}

// 将用户文件信息插入数据库用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filsize int64) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user_file(`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`)values(?,?,?,?,?)")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	_, err = stmt.Exec(username, filehash, filename, filsize, time.Now())
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	return true
}

// 查询用户网盘里的文件
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1, file_name, file_size, upload_at from tbl_user_file where user_name=? limit ?",
	)
	if err != nil {
		fmt.Printf("err1: %v\n", err)
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	var userfiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		userfiles = append(userfiles, ufile)
	}
	return userfiles, nil

}

// 删除用户的文件信息
func DeleteUserFile(username, filehash string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"DELETE FROM tbl_user_file WHERE user_name=? AND file_sha1=?")
	fmt.Println("删除", username, filehash)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	_, err = stmt.Exec(username, filehash)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	return true
}
