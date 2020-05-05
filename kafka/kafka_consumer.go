package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/dongysh/project/logagent/es"
)

var consumer sarama.Consumer

// ConsumerTask 一个消费kafka的任务
type ConsumerTask struct {
	topic    string
	instance sarama.Consumer
}

// NewConsumerTask 创建一个消费任务
func NewConsumerTask(project, addr, topic string) (consumerTask *ConsumerTask) {
	consumerTask = &ConsumerTask{
		topic: topic,
	}
	consumerTask.InitConsumer(project, addr)
	return
}

// InitConsumer 初始化消费者
func (t ConsumerTask) InitConsumer(project, addr string) (err error) {
	consumer, err = sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		fmt.Printf("connect kafka consumer failed, err:%v", err)
		return
	}

	fmt.Printf("connect kafka consumer success\n")
	go sendToES(project, t.topic)
	return
}

// sendToES 将数据发送给ES
func sendToES(project, topic string) (err error) {
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有分区
	if err != nil {
		fmt.Println("list of partition falied, err:", err)
		return err
	}
	fmt.Println("分区列表", partitionList)
	for partition := range partitionList { // 遍历所有分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d, err:%v\n", partition, err)
			return err
		}
		//异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("partition:%d offset:%d key:%v value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				// 直接发送给ES
				ld := map[string]interface{}{
					"data": string(msg.Value),
				}
				err = es.SendToES(project, topic, ld)
				if err != nil {
					fmt.Println("send to es failed, err:", err)
				}
			}
		}(pc)
	}
	return err
}
