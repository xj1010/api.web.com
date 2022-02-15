package service

import (
	"admin/db"
	"context"
	"testing"
	"time"
)

func Test_lagou(t *testing.T) {
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)
	BaseContext, _ := context.WithCancel(context.Background())

	service := NewLaGouService(BaseContext)
	service.Start()
	for {
		time.Sleep(1000)
	}

}
