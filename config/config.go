package config

// AppConf App全局配置
type AppConf struct {
	Project     string `ini:"project"`
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
	EtcdConf    `ini:"etcd"`
	EsConf      `ini:"es"`
}

// KafkaConf ...
type KafkaConf struct {
	Address     string `ini:"address"`
	Topic       string `ini:"topic"`
	ChanMaxSize int    `ini:"chan_max_size"`
}

// TaillogConf ...
type TaillogConf struct {
	FileName string `ini:"fileName"`
}

// EtcdConf ...
type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"key"`
}

// EsConf ...
type EsConf struct {
	Address string `ini:"address"`
}
