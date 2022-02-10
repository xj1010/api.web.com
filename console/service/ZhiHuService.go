package service

import (
	"admin/models"
	"admin/utils"
	"github.com/siddontang/go/log"
	"regexp"
	"strconv"
	"strings"
)

type ZhiHuService struct {
	BaseService
}

func (s *ZhiHuService) Execute() *ZhiHuService {
	s.Desc = "知乎热榜"
	s.Run(s.Work)
	return s
}

func(s *ZhiHuService) Work() {
	link := "https://www.zhihu.com/billboard"
	req := utils.NewBrowser()
	str, status := req.Get(link)

	if len(str) == 0 || status != 200 {
		log.Info("知乎网站采集失败")
		return
	}

	htmlStr := string(str)
	r, _:= regexp.Compile(`"(https:(?:\\u002F\\u002Fwww.zhihu.com\\u002F)(?:question|campaign|special).*?)"`)
	pageLinkMatch := r.FindAllStringSubmatch(htmlStr, -1)
	if len(pageLinkMatch) == 0 {
		log.Info("知乎网站link没有匹配到信息")
		return
	}

	r, _ = regexp.Compile(`<div class="HotList-itemIndex\s*(?:HotList-itemIndexHot)?">(\d+)</div><div class="HotList-itemLabel" style="color:(?:.*?)"></div></div><div class="HotList-itemBody"><div class="HotList-itemTitle">(.*?)</div><div class="HotList-itemMetrics">(\d+)\s*万热度</div>`)
	urlTitleMatches := r.FindAllStringSubmatch(htmlStr, -1)
	if len(urlTitleMatches) == 0 {
		log.Info("知乎网站没有匹配到信息 ")
		return
	}

	var updateTotal, addTotal int

	for j, urlTitle := range urlTitleMatches {
		var hotTop models.HotTop
		link := strings.Replace(pageLinkMatch[j][1], "\\u002F", "/",  -1)
		hotTop.Site = "知乎"
		hotTop.Title = urlTitle[2]
		hotTop.Link = link

		num, _ := strconv.ParseInt(urlTitle[3], 10, 64)
		hotTop.Num = num*10000

		sortNum, _ := strconv.Atoi(urlTitle[1])
		hotTop.Sort = sortNum

		updateNum, addNum :=  hotTop.AddOrUpdateHotTopNew(hotTop)

		updateTotal += updateNum
		addTotal += addNum
	}

	log.Infof(s.Desc + "采集成功,更新%d, 新增%d", updateTotal, addTotal)
}





