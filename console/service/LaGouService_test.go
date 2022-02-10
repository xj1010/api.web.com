package service

import (
	"admin/db"
	"testing"
)

func Test_findByPk(t *testing.T) {
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)
	lagou :=  LaGouService{}
	lagou.IsRunning = true

	lagou.Work()

}
