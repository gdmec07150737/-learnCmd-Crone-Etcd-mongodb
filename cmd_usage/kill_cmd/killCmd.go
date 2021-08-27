package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	output []byte
	err error
}

func main() {
	var (
		cmd *exec.Cmd
		ctx context.Context
		cancelFunc context.CancelFunc
		resultChan chan *result
		res *result
	)
	resultChan = make(chan *result, 1000)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	//ctx, cancelFunc = context.WithCancel(context.Background())
	go func() {
		var (
			output []byte
			err error
		)
		cmd = exec.CommandContext(ctx, "bash", "-c", "sleep 2;ls -l")
		output, err = cmd.CombinedOutput()
		resultChan <- &result{
			output: output,
			err:    err,
		}
	}()
	time.Sleep(time.Second * 1)
	cancelFunc()
	res = <- resultChan
	fmt.Println(string(res.output))
	fmt.Println(res.err)
}