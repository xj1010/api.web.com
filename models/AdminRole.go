package models

import (
	"admin/db"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type  AdminRole struct {
	Id int `gorm:"primary_key" json:"id"   form:"id"`
	Name string `gorm:"column:name;type:varchar(150)" json:"name"   form:"name"`
	Status int  `gorm:"column:status;type:int(10);"  json:"status"    form:"status" `
	Remark string `gorm:"column:remark;type:varchar(150)" json:"remark"   form:"remark"`
	Base
}
func init() {
}

func (adminRole *AdminRole) GetRoleInfoByUsername(username string)  AdminRole {
	dbCon := db.GetInstance().GetMysqlDb()
	var role  AdminRole
	dbCon.Where("name=?",  username ).First(&role)
	return role
}

func (adminRole *AdminRole) UpdateRoleInfoByUid(uid int, data map[string]interface{}) bool {
	dbCon := db.GetInstance().GetMysqlDb()
	err := dbCon.Model(&adminRole).Where("id=?", uid).Update(data).Error
	if err != nil {
		return true
	}
	return false
}

func (adminRole *AdminRole) GetList(username , status string, page, limit int) (count int, roleInfo []AdminRole) {
	dbCon := db.GetInstance().GetMysqlDb()

	username = strings.TrimSpace(username)
	if len(username) > 0  {
		dbCon = dbCon.Where("name=?", username)
	}

	status = strings.TrimSpace(status)
	if len(status) > 0 {
		status, err := strconv.Atoi(status)
		if err == nil {
			dbCon = dbCon.Where("status=?", status)
		} else {
			dbCon = dbCon.Where("status=?", -1)
		}
	}

	if page <= 0 {
		page = 1
	}

	//总数
	dbCon.Model(&adminRole).Count(&count)
	fPage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(count) / float64(limit)), 64)
	maxPage := int(math.Ceil(fPage))
	if page > maxPage {
		page = maxPage
	}

	//数据
	dbCon = dbCon.Order("created_at desc").Limit(limit).Offset((page - 1) * limit)

	if err := dbCon.Find(&roleInfo).Error; err != nil {
		return 0, roleInfo
	}
	return count, roleInfo
}
