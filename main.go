package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dongysh/project/logagent/config"
	"github.com/dongysh/project/logagent/es"
	"github.com/dongysh/project/logagent/etcd"
	"github.com/dongysh/project/logagent/kafka"
	"github.com/dongysh/project/logagent/taillog"
	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.AppConf)
)

func main() {
	// 1、加载配置并初始化kafka、etcd、es
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	// 初始化kafka生产者
	err = kafka.InitProducer([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Println("init kafka failed, err:", err)
		return
	}
	fmt.Println("init kafka success")
	// 初始化etcd连接
	err = etcd.Init([]string{cfg.EtcdConf.Address}, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd failed, err:", err)
		return
	}
	fmt.Println("init etcd success")
	// 初始化es连接
	err = es.Init(cfg.EsConf.Address)
	if err != nil {
		fmt.Println("init es failed, err:", err)
		return
	}
	fmt.Println("init es success")

	// 2、读取etcd配置，获取日志相关配置项
	key := fmt.Sprintf(cfg.EtcdConf.Key, cfg.Project) //获取项目配置
	logEntryConf, err := etcd.GetConf(key)
	if err != nil {
		fmt.Println("get conf from etcd failed, err:", err)
	}
	fmt.Println("get conf from etcd success!")

	// 3、收集日志并发往kafka
	taillog.Init(logEntryConf)

	// 4、消费kafka，将日志发往es
	kafka.InitConsumerMgr(cfg.Project, cfg.KafkaConf.Address, logEntryConf)

	select {}
}
