package service

import (
	"admin/db"
	"testing"
)

func Test_baidu(t *testing.T) {
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)
	lagou :=  BaiDuService{}
	lagou.IsRunning = true

	lagou.Work()

}
