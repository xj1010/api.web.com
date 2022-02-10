package service

import (
	"github.com/siddontang/go/log"
	"math/rand"
	"os"
	"time"
)

type BaseService struct {
	IsRunning   bool
	Desc string
}

func (b *BaseService)  Start() {
	 b.IsRunning = true
}

func (b *BaseService)  Stop() {
	b.IsRunning = false
}

func (b *BaseService)  GetRunningStatus() bool {
	return b.IsRunning
}


func (b *BaseService)  Run(f func()) {
	go func() {
		defer func() {
			if err := recover(); err !=nil {
				log.Info(err)
			}
		}()
		for {
			if b.IsRunning {
				f()
			}

			rand.Seed(time.Now().UnixNano())
			sleepNum := rand.Intn(120)
			timer1 := time.NewTimer(time.Duration(sleepNum) * time.Second)
			<- timer1.C
			timer1.Stop()
		}


	}()
}

func (b *BaseService) DirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	// 判断报错是否是因为文件不存在引起的
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 新建文件夹，必要时创建文件的父级目录
func (b *BaseService) MkDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}


