package service

import (
	"admin/models"
	"admin/utils"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/siddontang/go/log"
	"strconv"
	"strings"
)

type BaiDuService struct {
	BaseService
}


func NewBaiDuService(BaseCxt context.Context) *BaiDuService {
	s := &BaiDuService{}
	s.IsRunning = false
	s.Desc = "百度"
	s.ServiceID = "BaiDu"
	s.Task = s.Work
	s.baseCxt = BaseCxt
	s.SetContext()
	return s
}

func(s *BaiDuService) Work() {
	link := "https://top.baidu.com/board?tab=realtime"
	req := utils.NewBrowser()
	str, status := req.Get(link)

	if len(str) == 0 || status != 200 {
		fmt.Println(string(str))
		log.Info("百度网站采集失败:" + strconv.Itoa(status))
		return
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(str)))
	if err != nil {
		log.Info(err)
		return
	}

	var updateTotal, addTotal int
	dom.Find(".horizontal_1eKyQ").Each(func(i int, selection *goquery.Selection) {
		sortNum , _ := selection.Find(".index_1Ew5p").Html()
		selection.Find(".content_1YWBm .title_dIF3B .hot-tag_1G080 ").Remove()
		hrefSelect := selection.Find(".content_1YWBm .title_dIF3B")
		link , _ := hrefSelect.Attr("href")
		title , _ := hrefSelect.Find(".c-single-text-ellipsis").Html()
		hotNum , _ := selection.Find(".hot-index_1Bl1a").Html()

		var hotTop models.HotTop
		hotTop.Site = "百度"
		hotTop.Title = strings.TrimSpace(title)
		hotTop.Link = strings.TrimSpace(link)

		num, _ := strconv.ParseInt(strings.TrimSpace(hotNum), 10, 64)
		hotTop.Num = num

		sort, _ := strconv.Atoi(strings.TrimSpace(sortNum))
		hotTop.Sort = sort

		updateNum, addNum :=  hotTop.AddOrUpdateHotTopNew(hotTop)

		updateTotal += updateNum
		addTotal += addNum
	})

	log.Infof(s.Desc + "采集成功,更新%d, 新增%d", updateTotal, addTotal)

}


