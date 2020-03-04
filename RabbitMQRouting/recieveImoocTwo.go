package main

import "mq/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQRouting("exImooc","imooc_two")
	imoocOne.RecieveRouting()
}
