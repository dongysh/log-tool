package kafka

import (
	"github.com/dongysh/project/logagent/etcd"
)

// 用于管理Kafak消费者 struct
type kafkaConsumerMgr struct {
	logEntry []*etcd.LogEntry
	taskMap  map[string]*ConsumerTask
}

var kafkaMgr *kafkaConsumerMgr

// InitConsumerMgr 初始化kakfa消费者的管理者
func InitConsumerMgr(project, addr string, logEntryConf []*etcd.LogEntry) {
	kafkaMgr = &kafkaConsumerMgr{
		logEntry: logEntryConf,
	}
	// 为每个topic添加任务
	for _, logEntry := range logEntryConf {
		NewConsumerTask(project, addr, logEntry.Topic)
	}
}
