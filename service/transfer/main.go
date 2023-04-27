package main

import (
	"encoding/json"
	"fileserver/config"
	"fileserver/db"
	"fileserver/mq"
	"fileserver/store/oss"
	"fmt"
	"os"
	"time"
)

func ProcessTransfer(msg []byte) bool {
	// 1.解析msg
	pubData := mq.TransferData{}
	fmt.Println(string(msg))
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// 2 根据临时存储文件路径  创建文件句柄
	filed, err := os.Open("../." + pubData.CurLocation)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// 3 通过文件句柄读取文件并传入oss
	fmt.Printf("pubData.DestLocation: %v\n", pubData.DestLocation)
	//filed.Seek(0, 0)
	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		filed,
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// 4 更新oss文件存储路径到文件表
	b := db.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	if !b {
		fmt.Println("上传失败")
		return false
	}
	// 5. 关闭句柄和删除临时文件
	filed.Close()
	fmt.Println("删除临时文件")
	os.Remove(pubData.CurLocation)
	return true
}
func main() {
	// fmt.Println("开始接收")
	// _, err := os.Open("../../tmp/新建文本文档(2).txt")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	time.Sleep(time.Second * 10)
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer,
	)
}
