package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"mqSimple")
	rabbitmq.ConsumeSimple()
}
