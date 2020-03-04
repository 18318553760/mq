package main

import (
	"fmt"
	"mq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"mqSimple")

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello mq!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
