go get github.com/ichunt2019/go-msgserver

# 发送消息：
##### demo
    queueExchange := &rabbitmq.QueueExchange{
          "b_test_rabbit",
          "b_test_rabbit",
          "b_test_rabbit_mq",
          "direct",
          "amqp://guest:guest@192.168.2.232:5672/",
    }
    mq := rabbitmq.New(queueExchange)
    for i := 0;i<10;i++{
        mq.RegisterProducer("这是测试任务"+strconv.Itoa(i))
    }
    err := mq.Start()
    if(err != nil){
        fmt.Println("发送消息失败")
    }


# 消费消息:
##### 必须实现接口 Receiver
>       type Receiver interface {
    	Consumer([]byte)    error
    	FailAction([]byte)  error
    }


DEMO

    type RecvPro struct {
    
    }
    
    //// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
    func (t *RecvPro) Consumer(dataByte []byte) error {
    	fmt.Println(string(dataByte))
    	return errors.New("顶顶顶顶")
    	//return nil
    }
    
    //消息已经消费3次 失败了 请进行处理
    func (t *RecvPro) FailAction(dataByte []byte) error {
    	fmt.Println(string(dataByte))
    	fmt.Println("任务处理失败了，我要进入db日志库了")
    	fmt.Println("任务处理失败了，发送钉钉消息通知主人")
    	return nil
    }
    
    
    
    func main() {
    
    	//消费者实现 下面接口即可
    	//type Receiver interface {
    	//	Consumer([]byte)    error
    	//}
    	t := &RecvPro{}
    
    
    	queueExchange := &rabbitmq.QueueExchange{
    		"fengkong_dong_count",
    		"fengkong_dong_count",
    		"fengkong_exchange",
    		"direct",
    		"amqp://guest:guest@192.168.2.232:5672/",
    	}
    	for{
    
    		mq := rabbitmq.New(queueExchange)
    		mq.RegisterReceiver(t)
    		err :=mq.Start()
    		if err != nil{
    
    			fmt.Println(err)
    		}
    		time.Sleep(time.Second)
    	}
    }

###### 当我们消息3次都处理失败之后 ：
> {"com_credits_id":"2","erp_company_code":"LX001"}
消息处理失败，进入延时队列.....
{"com_credits_id":"2","erp_company_code":"LX001"}
消息处理失败，进入延时队列.....
{"com_credits_id":"2","erp_company_code":"LX001"}
消息处理失败，进入延时队列.....
{"com_credits_id":"2","erp_company_code":"LX001"}
消息处理失败 入库db{"com_credits_id":"2","erp_company_code":"LX001"}
任务处理失败了，我要进入db日志库了
任务处理失败了，发送钉钉消息通知主人



## 并发执行
    for{
    	var wg sync.WaitGroup
    	fmt.Println("开始执行任务....")
    	for i := 0;i<10;i++{
    	wg.Add(1)
    	go func(wg *sync.WaitGroup){
    	mq := rabbitmq.New(queueExchange)
    	mq.RegisterReceiver(t)
    	err :=mq.Start()
    	if err != nil{
    
    	fmt.Println(err)
    	}
    	wg.Done()
    	}(&wg)
    	}
    	wg.Wait()
    	fmt.Println("执行任务完成....")
    	time.Sleep(time.Microsecond*10)
    }


