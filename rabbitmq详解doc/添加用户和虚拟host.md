# 概述

## rabbitmq的安装

### 1、安装方式一，具体请移步本人之前写的文章，

<http://cblog.1024.company/3-14.html> 

### 2、使用rpm安装

#### a、下载文件

wget http://www.rabbitmq.com/releases/erlang/erlang-19.0.4-1.el7.centos.x86_64.rpm

wget  http://www.rabbitmq.com/releases/rabbitmq-server/v3.6.10/rabbitmq-server-3.6.10-1.el7.noarch.rpm

#### b、安装erlang

rpm -ivh erlang-19.0.4-1.el7.centos.x86_64.rpm

#### c、安装RabbitMQ

rpm -ivh rabbitmq-server-3.6.10-1.el7.noarch.rpm

#### d、启动Rabbitmq

systemctl start rabbitmq-server

#### e、rabbitmq常用命令



1、systemctl start rabbitmq-server  启动

2、rabbitmqctl stop  停止

3、rabbitmq-plugins  list 查看插件列表

4、rabbitmq-plugins enable rabbitmq_management 启用rabbitmq_management列表

5、添加用户

rabbitmqctl add_user admin 123456添加用户

rabbitmqctl set_user_tags admin administrator  设置标签

rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"  设置权限

浏览器打开 “http://127.0.0.1:15672″， 用帐号“admin”, 密码“123456” 就可以登录了，默认是guest,guest但是只对localhost域名有效。



## 建立虚拟host，绑定用户

### 1、创建用户

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\3.png)

### 2、创建虚拟host

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\1.png)

### 3、绑定用户，点击虚拟host下的mq,进入设置详情页面，选择绑定User

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\2.png)

通过新建立的用户访问<http://127.0.0.1:15672/>

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\4.png)

## golang在rabbitmq实践

### 1、新建mq项目

### 2、RabbitMQ/RabbitMQ

```
package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//连接信息
const MQURL = "amqp://mqUser:mqUser@127.0.0.1:5672/mq"

//rabbitMQ结构体
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange  string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl     string
}

//创建结构体实例
func NewRabbitMQ(queueName string,exchange string ,key string) *RabbitMQ {
	return &RabbitMQ{QueueName:queueName,Exchange:exchange,Key:key,Mqurl:MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}
//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName,"","")
	var err error
    //获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err :=r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"",     // consumer
		//是否自动应答
		true,   // auto-ack
		//是否独有
		false,  // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false,  // no-local
		//列是否阻塞
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
    //启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// 接收done的信号, 没有信息过来则会一直阻塞，避免该函数退出
	<-forever
	// 关闭通道
	r.channel.Close()

}
//订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("",exchangeName,"")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err,"failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}
//订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
//订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"fanout",
		true,
		false,
		//YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")
    //2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}
//路由模式
//创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string,routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err,"failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string )  {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//要改成direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
//路由模式接受消息
func (r *RabbitMQ) RecieveRouting() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定key
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}
//话题模式
//创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string,routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err,"failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}
//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string )  {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//要改成topic
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 imooc.* 表示匹配 imooc.hello, 但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}
```

### 3、简单模式

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\5.png)

#### mainSimlpePublish.go

```
package main

import (
	"mq/RabbitMQ"
	"fmt"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"mqSimple")
	rabbitmq.PublishSimple("Hello mq!")
	fmt.Println("发送成功！")
}

```

运行上述命令代码，可以看到相关对列的生成，点击**mqSimple** 进入详情页面，Get message可以看到查看对应的信息

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\6.png)

#### mainSimpleRecieve.go

```
package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"mqSimple")
	rabbitmq.ConsumeSimple()
}

```

运行上述代码，可以获取到mqSimple信息

### 4、work模式，起到一个负载均衡的作用

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\7.png)

work模式跟simple模式一样的，只是创建两个不同的进程去消费消息

#### RabbitMQWork/mainWorkPublish.go

```
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

```

#### RabbitMQWork/mainWorkRecieve1.go

```
package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"mqSimple")
	rabbitmq.ConsumeSimple()
}

```

#### RabbitMQWork/mainWorkRecieve2.go

```
package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"mqSimple")
	rabbitmq.ConsumeSimple()
}

```

分别运行mainWorkRecieve1.go与mainWorkRecieve2.go等待生产者生成消息，运行mainWorkPublish.go查看消费者不重复的消费

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\8.png)

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\9.png)

### 5、Publish/Subscribe,订阅模式

![1583309642023](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\13.png)

RabbitMQPubSub/mainPub.go

```
package main

import (
	"fmt"
	"mq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		fmt.Println("订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1 * time.Second)
	}

}

```

RabbitMQPubSub/mainSub.go

```
package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	rabbitmq.RecieveSub()
}

```

RabbitMQPubSub/mainSub1.go

```
package main

import "mq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	rabbitmq.RecieveSub()
}

```

### 6、路由模式

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\14.png)

#### RabbitMQRouting/publishRouting.go

```
package main

import (
	"fmt"
	"mq/RabbitMQ"
	"strconv"
	"time"
)

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQRouting("exMq","mq_one")
	imoocTwo:=RabbitMQ.NewRabbitMQRouting("exMq","mq_two")
	for i := 0; i <= 10; i++ {
		imoocOne.PublishRouting("Hello mq one!" + strconv.Itoa(i))
		imoocTwo.PublishRouting("Hello mq Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
	
}

```

#### RabbitMQRouting/recievemqOne.go

```
package main

import "mq/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQRouting("exMq","mq_one")
	imoocOne.RecieveRouting()
}

```

#### RabbitMQRouting/recievemqTwo.go

```
package main

import "mq/RabbitMQ"

func main()  {
   imoocOne:=RabbitMQ.NewRabbitMQRouting("exMq","mq_two")
   imoocOne.RecieveRouting()
}
```

分别运行recievemqOne.go，recievemqTwo.go等待消费信息，运行publishRouting.go生产消息

### 7、Topic模式

![](C:\Users\Administrator\Desktop\go语言\go\rabbitmq\images\15.png)

#### RabbitMQTopic/publishTopic.go

```
package main

import (
	"mq/RabbitMQ"
	"strconv"
	"time"
	"fmt"
)

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQTopic("exMqTopic","mq.topic.one")
	imoocTwo:=RabbitMQ.NewRabbitMQTopic("exMqTopic","mq.topic.two")
	for i := 0; i <= 10; i++ {
		imoocOne.PublishTopic("Hello mq topic one!" + strconv.Itoa(i))
		imoocTwo.PublishTopic("Hello mq topic Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
	
}

```

#### RabbitMQTopic/recieveImoocTopicAll.go 

routingkey为#是获取所有

```
package main

import "mq/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQTopic("exMqTopic","#")
	imoocOne.RecieveTopic()
}

```

#### RabbitMQTopic/recieveImoocTopicTwo.go

```
package main

import "mq/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQTopic("exImoocTopic","mq.*.two") // 一个
	imoocOne.RecieveTopic()
}

```

分别运行recieveImoocTopicAll.go ,recieveImoocTopicTwo.go等待消费信息，运行publishTopic.go生产消息