package es

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivere/elastic"
)

var client *elastic.Client

// Init 初始化es 准备接收 kafka那边发来的数据
func Init(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(address))
	if err != nil {
		return
	}
	fmt.Println("conn to es success")
	return
}

// SendToES 发送数据到ES
func SendToES(indexStr, typeStr string, data interface{}) (err error) {
	fmt.Printf("prepare to send to es,index:%s, type:%s data:%v\n", indexStr, typeStr, data)
	put1, err := client.Index().
		Index(indexStr).
		Type(typeStr).
		BodyJson(data).
		Do(context.Background())
	if err != nil {
		// Handle error
		return
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	return
}
