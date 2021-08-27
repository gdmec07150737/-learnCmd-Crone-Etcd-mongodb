package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		opResp clientv3.OpResponse
		putOp clientv3.Op
		getOp clientv3.Op
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	putOp = clientv3.OpPut("/demo7", "value")
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入Revision：", opResp.Put().Header.Revision)
	getOp = clientv3.OpGet("/demo7")
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("数据Revision：", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value：", string(opResp.Get().Kvs[0].Value))
}