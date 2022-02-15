package service

import (
	"context"
	"fmt"
	"github.com/siddontang/go/log"
	"math/rand"
	"time"
)

type BaseService struct {
	IsRunning   bool
	Desc string
	ServiceID string
	CancelFunc  context.CancelFunc
	Cxt context.Context
	baseCxt context.Context
	Task func()
}

func(b *BaseService) SetContext() {
	cxt, cancelFunc := context.WithCancel(b.baseCxt)
	b.Cxt = cxt
	b.CancelFunc = cancelFunc
}

func (b *BaseService)  Start() {
	 if b.IsRunning {
	 	return
	 }

	 b.IsRunning = true
	 b.Run()
}

func (b *BaseService)  Stop() {
	b.IsRunning = false
	b.CancelFunc()

}

func (b *BaseService)  Run()  {
	go func() {
		defer func() {
			if err := recover(); err !=nil {
				log.Info(err)
			}
		}()

		for {
			select {
				case <-b.Cxt.Done():
					//重置context
					b.SetContext()
					fmt.Println(b.ServiceID +"结束了")
					return
				default:
					b.Task()
			}

			rand.Seed(time.Now().UnixNano())
			sleepNum := rand.Intn(10)
			timer1 := time.NewTimer(time.Duration(sleepNum) * time.Second)
			<- timer1.C
			timer1.Stop()

		}

	}()

}



