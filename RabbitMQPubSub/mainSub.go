package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	rabbitmq.RecieveSub()
}
