package main

import (
	"mq/RabbitMQ"
	"fmt"
)

func main() {
	for i:=0;i<10;i++ {
		rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
			"mqSimple")
		rabbitmq.PublishSimple(fmt.Sprintf("Hello mq %d",i))
		fmt.Println("发送成功！")
	}

}
