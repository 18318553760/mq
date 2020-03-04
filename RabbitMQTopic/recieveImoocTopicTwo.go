package main

import "mq/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQTopic("exImoocTopic","imooc.*.two")
	imoocOne.RecieveTopic()
}
