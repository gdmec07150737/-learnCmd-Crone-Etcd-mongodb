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
		getResp *clientv3.GetResponse
		//putResp *clientv3.PutResponse
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

	/*//写etcd的键值对
	if putResp, err = kv.Put(context.TODO(), "/etcd/job/data1", "update2", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PreValue:", string(putResp.PrevKv.Value))
			fmt.Println(putResp.PrevKv)
		}
	}*/

	//读取etcd的键值对
	if getResp, err = kv.Get(context.TODO(), "/etcd/job/data1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs)
	}
}