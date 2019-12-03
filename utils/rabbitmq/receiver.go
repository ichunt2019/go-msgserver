package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)


// 定义全局变量,指针类型
var mqConn *amqp.Connection
var mqChan *amqp.Channel

// 定义生产者接口
type Producer interface {
	MsgContent() string
}

// 定义生产者接口
type RetryProducer interface {
	MsgContent() string
}

// 定义接收者接口
type Receiver interface {
	Consumer([]byte)    error
}

// 定义RabbitMQ对象
type RabbitMQ struct {
	connection *amqp.Connection
	channel *amqp.Channel
	dns string
	queueName   string            // 队列名称
	routingKey  string            // key名称
	exchangeName string           // 交换机名称
	exchangeType string           // 交换机类型
	producerList []Producer
	retryProducerList []RetryProducer
	receiverList []Receiver
	mu  sync.RWMutex
	wg  sync.WaitGroup
}

// 定义队列交换机对象
type QueueExchange struct {
	QuName  string           // 队列名称
	RtKey   string           // key值
	ExName  string           // 交换机名称
	ExType  string           // 交换机类型
	Dns     string			  //链接地址
}

// 链接rabbitMQ
func (r *RabbitMQ)mqConnect() (err error){
	//var err error
	//RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", "guest", "guest", "192.168.2.232", 5672)
	//mqConn, err = amqp.Dial(RabbitUrl)
	mqConn, err = amqp.Dial(r.dns)
	r.connection = mqConn   // 赋值给RabbitMQ对象
	if err != nil {
		return err
		//fmt.Printf("MQ打开链接失败:%s \n", err)
	}
	mqChan, err = mqConn.Channel()
	r.channel = mqChan  // 赋值给RabbitMQ对象
	if err != nil {
		return err
		//fmt.Printf("MQ打开管道失败:%s \n", err)
	}
	return err
}

// 关闭RabbitMQ连接
func (r *RabbitMQ)mqClose() {
	// 先关闭管道,再关闭链接
	err := r.channel.Close()
	if err != nil {
		fmt.Printf("MQ管道关闭失败:%s \n", err)
	}
	err = r.connection.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
	}
}

// 创建一个新的操作对象
func New(q *QueueExchange) *RabbitMQ {
	return &RabbitMQ{
		queueName:q.QuName,
		routingKey:q.RtKey,
		exchangeName: q.ExName,
		exchangeType: q.ExType,
		dns:q.Dns,
	}
}

// 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) Start() (err error){
	// 开启监听生产者发送任务
	for _, producer := range r.producerList {
		err = r.listenProducer(producer)
	}


	// 开启监听接收者接收任务
	for _, receiver := range r.receiverList {
		//r.listenReceiver(receiver)
		r.wg.Add(1)
		go func() {
			err = r.listenReceiver(receiver)
		}()

	}
	r.wg.Wait()
	time.Sleep(time.Microsecond*100)
	return err
}




type SendRbmqPro struct {
	msgContent   string
}

// 实现生产者
func (t *SendRbmqPro) MsgContent() string {
	return t.msgContent
}



// 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(msg string) {
	a := &SendRbmqPro{msgContent:msg}
	a.MsgContent()
	r.producerList = append(r.producerList, a)
}



// 发送任务
func (r *RabbitMQ) listenProducer(producer Producer) (err error){
	// 验证链接是否正常,否则重新链接
	if r.channel == nil {
		err = r.mqConnect()
		if err !=nil {
			return err
		}
	}
	err =  r.channel.ExchangeDeclare(r.exchangeName, r.exchangeType, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("MQ注册交换机失败:%s \n", err)
	}


	_, err = r.channel.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("MQ注册队列失败:%s \n", err)
	}


	// 队列绑定
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true,nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
	}

	header := make(map[string]interface{},1)

	header["retry_nums"] = int32(0)

	// 发送任务消息
	err =  r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.MsgContent()),
		Headers:header,
	})

	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
	}
	return err
}


func (r *RabbitMQ) listenRetryProducer(producer RetryProducer,retry_nums int32 ,args ...string) {
	fmt.Println("消息处理失败，进入延时队列.....")
	//defer r.mqClose()
	// 验证链接是否正常,否则重新链接
	if r.channel == nil {
		r.mqConnect()
	}

	err :=  r.channel.ExchangeDeclare(r.exchangeName, r.exchangeType, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("MQ注册交换机失败:%s \n", err)
		return
	}


	//原始路由key
	oldRoutingKey := args[0]
	//原始交换机名
	oldExchangeName := args[1]

	table := make(map[string]interface{},3)
	table["x-dead-letter-routing-key"] = oldRoutingKey
	table["x-dead-letter-exchange"] = oldExchangeName

	table["x-message-ttl"] = int64(50000)

	_, err = r.channel.QueueDeclare(r.queueName, true, false, false, false, table)
	if err != nil {
		fmt.Printf("MQ注册队列失败:%s \n", err)
		return
	}


	// 队列绑定
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true,nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return
	}

	header := make(map[string]interface{},1)

	header["retry_nums"] = retry_nums + int32(1)

	// 发送任务消息
	err =  r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.MsgContent()),
		Headers:header,
	})

	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
}

// 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.mu.Lock()
	r.receiverList = append(r.receiverList, receiver)
	r.mu.Unlock()
}

// 监听接收者接收任务 消费者
func (r *RabbitMQ) listenReceiver(receiver Receiver) (err error) {
	// 处理结束关闭链接
	defer r.mqClose()
	defer r.wg.Done()
	//defer
	// 验证链接是否正常
	if r.channel == nil {
		err = r.mqConnect()
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = r.channel.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("MQ注册队列失败:%s \n", err)
		return errors.New("bbvv")
	}
	// 绑定任务
	err =  r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, false, nil)
	if err != nil {
		fmt.Printf("绑定队列失败:%s \n", err)
		return errors.New("bbvv")
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err =  r.channel.Qos(1, 0, false)
	msgList, err :=  r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
		return errors.New("bbvv")
	}
	for msg := range msgList {

		retry_nums := msg.Headers["retry_nums"].(int32)
		// 处理数据
		err := receiver.Consumer(msg.Body)
		if err!=nil {
			//消息处理失败 进入延时尝试机制
			if retry_nums < 3{
				r.retry_msg(msg.Body,retry_nums)
			}else{
				//消息失败 入库db
				fmt.Printf("消息处理失败 入库db")
			}
			err = msg.Ack(true)
			if err != nil {
				fmt.Printf("确认消息未完成异常:%s \n", err)
				return errors.New("bbvv")
			}
		}else {
			// 确认消息,必须为false
			err = msg.Ack(true)

			if err != nil {
				fmt.Printf("确认消息完成异常:%s \n", err)
				return errors.New("bbvv")
			}
			return err
		}
	}
	return err
}



type retryPro struct {
	msgContent   string
}

// 实现生产者
func (t *retryPro) MsgContent() string {
	return t.msgContent
}

//消息处理失败之后 延时尝试
func(r *RabbitMQ) retry_msg(Body []byte,retry_nums int32){
	queueName := r.queueName+"_retry_3"
	routingKey := r.queueName+"_retry_3"
	exchangeName := r.exchangeName
	queueExchange := &QueueExchange{
		queueName,
		routingKey,
		exchangeName,
		"direct",
		r.dns,
	}
	mq := New(queueExchange)
	msg := fmt.Sprintf("%s",Body)
	t := &retryPro{
		msg,
	}
	mq.listenRetryProducer(t,retry_nums,r.routingKey,exchangeName)
}