package mq

import (
	"fmt"
	"log"
)

var done chan bool

func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// 1 检查channel
	if !initChannel() {
		return
	}
	// 1． 通过channel．Consume获得消息信道
	fmt.Printf("qName: %v\n", qName)
	fmt.Printf("cName: %v\n", cName)
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//2．循环获取队列的消息
	done = make(chan bool)
	go func() {
		fmt.Println("循环接收")
		for msg := range msgs {
			// 用callback处理新消息
			procssSuc := callback(msg.Body)
			if !procssSuc {
				//将任务写到另一个队列 用于异常重试
			}
		}
	}()
	<-done
	channel.Close()
}
