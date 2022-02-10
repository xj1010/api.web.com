package models

import (
	"admin/db"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Lagou struct {
	ID           int     `gorm:"column:id" json:"id"`
	CompanySize  string    `gorm:"column:company_size" json:"company_size"`
	CompanyId    int     `gorm:"column:company_id" json:"company_id"`
	CompanyName  string    `gorm:"column:company_name" json:"company_name"`
	PositionId   int     `gorm:"column:position_id" json:"position_id"`
	PositionName string    `gorm:"column:position_name" json:"position_name"`
	FirstType    string    `gorm:"column:first_type" json:"first_type"`
	SecondType   string    `gorm:"column:second_type" json:"second_type"`
	ThirdType    string    `gorm:"column:third_type" json:"third_type"`
	City         string    `gorm:"column:city" json:"city"`
	SalaryMin    int     `gorm:"column:salary_min" json:"salary_min"`
	SalaryMax    int     `gorm:"column:salary_max" json:"salary_max"`
	SalaryMonth  int     `gorm:"column:salary_month" json:"salary_month"`
	WorkYear     string    `gorm:"column:work_year" json:"work_year"`
	Education    string    `gorm:"column:education" json:"education"`
	Createtime   JsonTime `gorm:"column:create_time" json:"create_time"`
}

type JsonTime time.Time

//实现它的json序列化方法
func (this JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (lagou *Lagou) AddOrUpdateRecord(laGou Lagou) (updateNum, addNum int) {
	dbCon := db.GetInstance().GetMysqlDb()
	var lg Lagou

	flag := dbCon.Where("position_id=?", laGou.PositionId).First(&lg).RecordNotFound()
	if !flag {
		updateMap := make(map[string]interface{}, 2)
		if laGou.SalaryMin != lg.SalaryMin {
			updateMap["SalaryMin"] = laGou.SalaryMin
		}
		if laGou.SalaryMax != lg.SalaryMax {
			updateMap["SalaryMax"] = laGou.SalaryMax
		}

		if len(updateMap) > 0 {
			if err := dbCon.Model(&Lagou{}).Where("id=?", lg.ID).Update(updateMap).Error; err != nil {
				return
			}
			updateNum = 1
		}
	} else {
		if dbCon.Create(&laGou).Error != nil {
			return
		}
		addNum = 1
	}

	return updateNum, addNum
}


func (lagou *Lagou) GetList(searchMap map[string]string, page, limit int) (count int, LagouList []Lagou) {
	dbCon := db.GetInstance().GetMysqlDb()

	var endDate string
	if len(searchMap["endDate"]) > 0 {
		d, err := time.ParseInLocation("2006-01-02", searchMap["endDate"], time.Local)
		if err != nil {
			panic("日期非法")
		}
		hour, _ := time.ParseDuration("24h")
		endDate = d.Add(hour).Format("2006-01-02")
	}

	if len(searchMap["startDate"]) > 0 && len(endDate) == 0 {
		dbCon = dbCon.Where("create_time >= ?", searchMap["startDate"])
	}

	if len(searchMap["startDate"]) > 0 && len(endDate) > 0 {
		dbCon = dbCon.Where("create_time >= ? and create_time <= ?", searchMap["startDate"], endDate)
	}

	if len(searchMap["startDate"]) == 0 && len(endDate) > 0 {
		dbCon = dbCon.Where("create_time <= ?", endDate)
	}

	if len(searchMap["searchType"]) > 0 && len(searchMap["keyword"]) > 0 {
		switch searchMap["searchType"] {
			case "company_name":
				dbCon = dbCon.Where("company_name like ?", "%"+strings.TrimSpace(searchMap["keyword"])+"%")
			case "position_name":
				dbCon = dbCon.Where("position_name like ?", "%"+strings.TrimSpace(searchMap["keyword"])+"%")
		    case "city":
			    dbCon = dbCon.Where("city = ?", strings.TrimSpace(searchMap["keyword"]))
		}
	}

	if page <= 0 {
		page = 1
	}

	//总数
	dbCon.Model(Lagou{}).Count(&count)
	fPage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(count) / float64(limit)), 64)
	maxPage := int(math.Ceil(fPage))
	if page > maxPage {
		page = maxPage
	}

	//数据
	dbCon = dbCon.Order("create_time desc").Order("id desc").Limit(limit).Offset((page - 1) * limit)
	if err := dbCon.Find(&LagouList).Error; err != nil {
		return 0, LagouList
	}
	return count, LagouList
}