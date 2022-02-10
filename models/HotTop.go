package models

import (
	"admin/db"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type  HotTop struct {
	Id int `gorm:"primary_key;column:id;" json:"id"    form:"id" `
	Site string  `gorm:"column:site;type:varchar(100);"  json:"site"    form:"site" `
	Title string `gorm:"column:title;type:varchar(255)" json:"title"   form:"title"`
	Link string  `gorm:"column:link;type:varchar(255)" json:"link"   form:"link"`
	Num int64  `gorm:"column:num;type:bigint(20);"  json:"num"    form:"num" `
	Sort int  `gorm:"column:sort;type:int(10);"  json:"sort"    form:"sort" `
	Base
}
func init() {
}

func (hot *HotTop) AddOrUpdateHotTopNew(hotNew HotTop) (updateNum, addNum int) {
	dbCon := db.GetInstance().GetMysqlDb()
	var news  HotTop

	flag := dbCon.Where("title=?",  hotNew.Title ).First(&news).RecordNotFound()
	if !flag {
		updateMap := make(map[string]interface{}, 2)
		if news.Sort != hotNew.Sort {
			updateMap["sort"] = hotNew.Sort
		}
		if news.Num != hotNew.Num {
			updateMap["num"] = hotNew.Num
		}

		if len(updateMap) > 0 {
			if err := dbCon.Model(&HotTop{}).Where("id=?", news.Id).Update(updateMap).Error; err != nil {
				return
			}
			updateNum = 1
		}
	} else {
		if dbCon.Create(&hotNew).Error != nil {
			return
		}
		addNum = 1
	}

	return updateNum, addNum
}

func (hot *HotTop) GetList(searchMap map[string]string, page, limit int) (count int, hotTops []HotTop) {
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
		dbCon = dbCon.Where("created_at >= ?", searchMap["startDate"])
	}

	if len(searchMap["startDate"]) > 0 && len(endDate) > 0 {
		dbCon = dbCon.Where("created_at >= ? and created_at <= ?", searchMap["startDate"], endDate)
	}

	if len(searchMap["startDate"]) == 0 && len(endDate) > 0 {
		dbCon = dbCon.Where("created_at <= ?", endDate)
	}

	if len(searchMap["site"]) > 0  {
		dbCon = dbCon.Where("site=?", searchMap["site"])
	}

	keyword := searchMap["title"]
	if keyword != "" {
		dbCon = dbCon.Where("title like ?", "%"+strings.TrimSpace(keyword)+"%")
	}

	if page <= 0 {
		page = 1
	}

	//总数
	dbCon.Model(HotTop{}).Count(&count)
	fPage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(count) / float64(limit)), 64)
	maxPage := int(math.Ceil(fPage))
	if page > maxPage {
		page = maxPage
	}

	//数据
	dbCon = dbCon.Order("updated_at desc").Order("created_at desc").Limit(limit).Offset((page - 1) * limit)
	if err := dbCon.Find(&hotTops).Error; err != nil {
		return 0, hotTops
	}
	return count, hotTops
}


