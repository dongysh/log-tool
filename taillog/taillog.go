package taillog

import (
	"fmt"

	"github.com/dongysh/project/logagent/kafka"
	"github.com/hpcloud/tail"
)

var (
	tailObj *tail.Tail
	logChan chan string
)

// TailTask 一个日志收集的任务
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

// NewTailTask TailTask构造函数
func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}
	tailObj.Init()
	return
}

// Init 初始单个 tailtask
func (t TailTask) Init() {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,                                 //轮询
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed. err:", err)
	}
	go t.run()
}

// run ...
func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines: //从tailObj的通道中一行一行的读取日志数据
			kafka.SendToChan(t.topic, line.Text)
		default:
		}
	}
}
