package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		kv             clientv3.KV
		leaseGrantResp *clientv3.LeaseGrantResponse
		lease          clientv3.Lease
		leaseId        clientv3.LeaseID
		keepRespChan   <- chan *clientv3.LeaseKeepAliveResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		ctx context.Context
		cancelFunc context.CancelFunc
		txn clientv3.Txn
		txnResp *clientv3.TxnResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//（创建租约，自动续租，拿着租约去抢占一个key）

	//1，创建租约
	lease = clientv3.NewLease(client)
	//申请一个5秒租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}
	//获取租约的id
	leaseId = leaseGrantResp.ID
	ctx, cancelFunc = context.WithCancel(context.TODO())

	//7，确保函数退出后，停止自动续租和取消租约
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)
	//defer cancelFunc()

	//2，自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	//3，处理续约应答协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {
					fmt.Println("收到自动续租：", keepResp.ID)
				}
			}
		}
	END:
	}()
	kv = clientv3.NewKV(client)

	//4，创建事务
	txn = kv.Txn(context.TODO())

	//5，事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision("/con/lock/pgc8"), "=", 0)).
		Then(clientv3.OpPut("/con/lock/pgc8", "pgc8", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/con/lock/pgc8"))//抢锁失败
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	//判断是否抢到锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	//6，处理业务
	fmt.Println("处理业务")
	time.Sleep(5 * time.Second)

	//7,释放锁（取消自动续约，释放租约）
	//defer会把租约释放掉，关联的KV就被删除了
}