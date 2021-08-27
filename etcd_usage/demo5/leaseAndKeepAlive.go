package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config       clientv3.Config
		client       *clientv3.Client
		err          error
		kv           clientv3.KV
		lease		 clientv3.Lease
		leaseGrantResp	*clientv3.LeaseGrantResponse
		leaseId	clientv3.LeaseID
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		keepRespchan <- chan *clientv3.LeaseKeepAliveResponse
		keepResp *clientv3.LeaseKeepAliveResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	//申请一个租约
	lease = clientv3.NewLease(client)
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	//获取租约的id
	leaseId = leaseGrantResp.ID
	//自动续租
	if keepRespchan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}
	//处理续约应答协程
	go func() {
		for {
			select {
			case keepResp = <- keepRespchan:
				if keepResp == nil {
					fmt.Println("租约已经失效")
					goto END
				} else {
					fmt.Println("续租成功：", keepResp.ID)
				}
			}
		}
		END:
	}()
	//申请一个kv
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "dada", "ddd", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("写入成功：", putResp.PrevKv.Version)
	fmt.Println("写入成功：", putResp.Header.Revision)
	for {
		if getResp, err = kv.Get(context.TODO(), "dada"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}