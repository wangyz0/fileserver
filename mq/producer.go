package mq

import (
	"fileserver/config"
	"fmt"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func initChannel() bool {
	//1判断channel是否已经创建
	if channel != nil {
		return true
	}
	//2获得rabbitmq链接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	//3打开channel 用于消息接受和发布
	channel, err = conn.Channel()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	return true

}

// 发布消息
func Publish(exchange, routinKey string, msg []byte) bool {
	// 1 检查channel
	if !initChannel() {
		return false
	}
	fmt.Printf("exchange: %v\n", exchange)
	fmt.Printf("routinKey: %v\n", routinKey)
	// 2 发布消息
	err := channel.Publish(
		exchange,
		routinKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	fmt.Println("发送", string(msg))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}
	return true
}
