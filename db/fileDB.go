package db

import (
	"database/sql"
	"fileserver/db/mysql"
	"fmt"
)

// 将文件的元信息保存到数据库
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false

	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("文件已存在")
		}
		return true
	}
	return false
}

// 查询信息
// 因为我们在meta包里用了本包  所以不能引用 因此要新建一个类filemeta结构体
type DBfilemeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

func GetFileMeta(fileSha1 string) (*DBfilemeta, error) {

	meta := DBfilemeta{}
	meta.FileSha1 = fileSha1
	// 从 MySQL 中获取文件的元信息
	err := mysql.DBConn().QueryRow("SELECT file_name, file_size, file_addr, created_at FROM tbl_file WHERE file_sha1=?", fileSha1).Scan(&meta.FileName, &meta.FileSize, &meta.Location, &meta.UploadAt)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Failed to select from mysql: %s\n", err.Error())
		return nil, err
	}
	return &meta, nil
}

// 在文件转移到OSS上后  修改地址

func UpdateFileLocation(filehash, location string) bool {

	stmt, err := mysql.DBConn().Prepare("update tbl_file set file_addr=? where file_sha1=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(location, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
