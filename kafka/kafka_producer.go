package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type logData struct {
	topic string
	line  string
}

var (
	producer    sarama.SyncProducer //声明一个全局的连接kafka的生产者producer
	logDataChan chan *logData
)

// InitProducer 初始化生产者
func InitProducer(addr []string, maxSize int) (err error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll          //发送完整数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true                   //成功交付的消息将在 success_channel返回

	//连接 kafka
	producer, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		fmt.Println("producer closed failed, err:", err)
		return
	}
	//初始化logDataChan
	logDataChan = make(chan *logData, maxSize)
	//开启后台的goroutine从通道中取数据发往kafka
	go sendToKafka()
	return
}

// SendToChan 给外部暴露的一个函数，把日志数据发送到内部的channel中
func SendToChan(topic, line string) {
	msg := &logData{
		topic: topic,
		line:  line,
	}
	logDataChan <- msg

}

// sendToKafka 真正往kafka发送日志
func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			//构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.line)
			//发送消息
			pid, offset, err := producer.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("send to kafaka topic:%s msg:%s pid:%v offset:%v\n", ld.topic, ld.line, pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}
