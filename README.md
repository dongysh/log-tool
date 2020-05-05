# log_tool
 日志小工具，基于go实现

背景
项目涉及多个系统开发，各系统日志分布在不同目录，不能集中展示，严重影响联调和测试效率

任务
收集项目中涉及的所有系统的日志文件，可视化展示

方案
GO语言开发、Etcd(配置中心)、Tail (go module、日志收集)、Kafka(日志缓存)、Elastic Search(日志持久化)、Kibana(日志可视化展示、检索)

软件版本
Go:1.42.2  Kafka:2.5.0   Zookeeper:3.5.7  Etcd:3.4.7  Elasticsearch 6.8.8  Kinaba:6.8.8

Go module
tail:github.com/hpcloud/tail
kafka:github.com/Shopify/sarama
etcd:go.etcd.io/etcd/clientv3
elasticsearch:github.com/olivere/elastic
