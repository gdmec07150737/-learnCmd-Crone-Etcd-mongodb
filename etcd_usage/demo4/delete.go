package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config       clientv3.Config
		client       *clientv3.Client
		err          error
		kv           clientv3.KV
		delResp      *clientv3.DeleteResponse
		keyValuePair *mvccpb.KeyValue
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

	if delResp, err = kv.Delete(context.TODO(), "/etcd/job/data2", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
		return
	}
	if len(delResp.PrevKvs) != 0 {
		for _, keyValuePair = range delResp.PrevKvs{
			fmt.Println("删除了：", string(keyValuePair.Key), string(keyValuePair.Value))
		}
	}
}