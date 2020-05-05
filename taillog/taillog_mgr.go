package taillog

import (
	"github.com/dongysh/project/logagent/etcd"
)

var taskMgr *tailLogMgr

// tailLogMgr 全局日志管理者
type tailLogMgr struct {
	logEntry []*etcd.LogEntry
	taskMap  map[string]*TailTask
}

// Init 初始化日志管理者
func Init(logEntryConf []*etcd.LogEntry) {
	taskMgr = &tailLogMgr{
		logEntry: logEntryConf,
	}
	// 为每个日志添加任务
	for _, logEntry := range logEntryConf {
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
}
