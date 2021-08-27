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
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		watcher clientv3.Watcher
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watchRespChan <-chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
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
	//模拟变化
	go func() {
		for {
			kv.Put(context.TODO(), "/etc/bin/test", "bash")
			kv.Delete(context.TODO(), "/etc/bin/test")
			time.Sleep(1 * time.Second)
		}
	}()
	//先get到当前值，并监听到后续变化
	if getResp, err = kv.Get(context.TODO(), "/etc/bin/test"); err != nil {
		fmt.Println(err)
		return
	}
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：", string(getResp.Kvs[0].Value))
	}
	//当前etcd的事务ID，单调递增
	watchStartRevision = getResp.Header.Revision + 1
	//创建一个watcher
	watcher = clientv3.NewWatcher(client)
	fmt.Println("从该版本后开始监听：", watchStartRevision)
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5 * time.Second, func() {
		cancelFunc()
	})
	//开始监听
	//watchRespChan = watcher.Watch(context.TODO(), "/etc/bin/test", clientv3.WithRev(watchStartRevision))
	watchRespChan = watcher.Watch(ctx, "/etc/bin/test", clientv3.WithRev(watchStartRevision))
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}