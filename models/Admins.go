package models

import (
	"admin/db"
	"fmt"
	"math"
	"strconv"
)

type Admins struct {
	Id        int     `gorm:"column:id" json:"id" form:"id"`
	Name      string    `gorm:"column:name" json:"name" form:"name"`
	Email     string    `gorm:"column:email" json:"email" form:"email"`
	Mobile    int64     `gorm:"column:mobile" json:"mobile,string" form:"mobile"`
	Password  string    `gorm:"column:password" json:"password" form:"password"`
	ImageUrl  string    `gorm:"column:image_url" json:"image_url" form:"image_url"`
	Status    int     `gorm:"column:status" json:"status" form:"status"`
	LoginNum  int     `gorm:"column:login_num" json:"login_num,string" form:"login_num"`
	Base
}

type AdminsWrap struct {
	ConfirmPassword  string    ` json:"confirmPassword" form:"confirmPassword"`
	ImagePath  string    ` json:"image_path" form:"image_path"`
	Admins
}

type LoginWrap struct {
	Username      string    `json:"username" form:"username"`
	Password      string    `json:"password" form:"password"`

}

func init() {
}

func (admins *Admins) GetAdminInfoByUsername(username string)  Admins {
	dbCon := db.GetInstance().GetMysqlDb()
	var admin  Admins
	dbCon.Where("name=?",  username ).First(&admin)
	return admin
}

func (admins *Admins) CheckAuth(username, password string) bool {
	dbCon := db.GetInstance().GetMysqlDb()
	var admin  Admins
	dbCon.Where("name=? and password=?",  username, password ).First(&admin)
	if admin.Id > 0 {
		return true
	}
	return false
}

func (admins *Admins) UpdateAdminInfoByUid(uid int, data map[string]interface{}) bool {
	dbCon := db.GetInstance().GetMysqlDb()
	err := dbCon.Model(&admins).Where("id=?", uid).Update(data).Error
	if err != nil {
		return true
	}
	return false
}

func (admins *Admins) GetList(searchMap map[string]string, page, limit int) (count int, adminInfo []Admins) {
	dbCon := db.GetInstance().GetMysqlDb()

	if _, ok := searchMap["status"]; ok {
		status, _ := strconv.Atoi(searchMap["status"])
		if status > 0 {
			dbCon = dbCon.Where("status=?", status)
		}
	}

	if _, ok := searchMap["username"]; ok {
		if len(searchMap["username"]) > 0 {
			dbCon = dbCon.Where("name=?", searchMap["username"])
		}
	}

	searchType, _ := strconv.Atoi(searchMap["searchType"])
	keyword := searchMap["keyword"]
	if searchType>0 && len(keyword)>0 {
		switch searchType {
		case 1:
			dbCon = dbCon.Where("name=?", keyword)
		case 2:
			uid, _ := strconv.Atoi(keyword)
			dbCon = dbCon.Where("id=?", uid)
		}
	}

	if page <= 0 {
		page = 1
	}

	//总数
	dbCon.Model(Admins{}).Count(&count)
	fPage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(count) / float64(limit)), 64)
	maxPage := int(math.Ceil(fPage))
	if page > maxPage {
		page = maxPage
	}

	//数据
	dbCon = dbCon.Order("created_at desc").Limit(limit).Offset((page - 1) * limit)

	if err := dbCon.Find(&adminInfo).Error; err != nil {
		return 0, adminInfo
	}
	return count, adminInfo
}
