package service

import (
	"admin/db"
	"context"
	"testing"
	"time"
)

func Test_baidu(t *testing.T) {
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)
	BaseContext, _ := context.WithCancel(context.Background())

	service := NewBaiDuService(BaseContext)
	service.Start()
	for {
		time.Sleep(1000)
	}

}
