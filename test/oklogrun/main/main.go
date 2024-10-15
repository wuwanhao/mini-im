package main

import (
	"context"
	"fmt"
	"github.com/oklog/run"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// goroutine 编排开始
func main1() {
	var g run.Group
	// 创建一个带有超时的上下文
	ctxAll, cancelAll := context.WithCancel(context.Background())
	fmt.Println(ctxAll)
	{
		// 处理信号退出的 channel
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancelC := make(chan struct{})

		// 注册 goroutine
		g.Add(func() error {
			select {
			case <-term:
				fmt.Println("Receive SIGTERM, exiting gracefully...")
				cancelAll()
				return nil
			case <-cancelC:
				fmt.Println("other cancel exiting")
				return nil
			}

		}, func(err error) {
			// 退出时做的清理操作
			close(cancelC)
		})
	}
	g.Run()
}

func main() {
	var g run.Group
	// 创建一个带有超时的上下文
	ctxAll, cancelAll := context.WithCancel(context.Background())
	{

		// 处理信号退出的 channel
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancelC := make(chan struct{})
		// 注册 goroutine
		g.Add(func() error {
			select {
			case <-term:
				fmt.Println("收到了终止信号，平滑关闭.....")
				cancelAll()
				return nil
			case <-cancelC:
				fmt.Println("other cancel exiting")
				return nil
			}
		}, func(err error) {
			// 退出时做的清理操作
			close(cancelC)
		})
	}

	{
		g.Add(func() error {
			for {
				ticker := time.NewTicker(time.Second)
				select {
				case <-ticker.C:
					fmt.Println("打工人01 正在工作")
				case <-ctxAll.Done():
					fmt.Println("打工人01，接收到了cancelAll的退出指令")
					return nil
				}

			}
		}, func(err error) {
			fmt.Println("clean up")
		})
	}

	g.Run()
}
