package service

import (
	"admin/db"
	"testing"
)

func Test_weibo(t *testing.T) {
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)
	wb :=  WeiBoService{}
	wb.IsRunning = true

	wb.Work()

}
