package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123123@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping to mysql, err:" + err.Error())
		os.Exit(1)
	}

	db.SetMaxOpenConns(1000)
}

func DBConn() *sql.DB {
	return db
}
