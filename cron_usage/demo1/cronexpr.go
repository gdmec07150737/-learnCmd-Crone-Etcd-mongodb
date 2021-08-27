package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expression *cronexpr.Expression
		err error
		nowTime time.Time
		nextTime time.Time
	)
	nowTime = time.Now()
	if expression, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}
	nextTime = expression.Next(nowTime)
	fmt.Println(nowTime, nextTime)
	time.AfterFunc(nextTime.Sub(nowTime), func() {
		fmt.Println("被调用了：", nextTime)
	})
	time.Sleep(time.Second * 6)
}