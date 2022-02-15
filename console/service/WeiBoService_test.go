package service

import (
	"admin/db"
	"context"
	"testing"
	"time"
)

func Test_weibo(t *testing.T) {
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)
	BaseContext, _ := context.WithCancel(context.Background())

	service := NewWeiBoService(BaseContext)
	service.Start()
	for {
		time.Sleep(1000)
	}

}
