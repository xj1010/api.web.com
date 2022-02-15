package service

import (
	"admin/models"
	"admin/utils"
	"context"
	"fmt"
	"github.com/siddontang/go/log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type WeiBoService struct {
	BaseService
}

func NewWeiBoService(BaseCxt context.Context) *WeiBoService {
	s := &WeiBoService{}
	s.IsRunning = false
	s.Desc = "微博"
	s.ServiceID = "WeiBo"
	s.Task = s.Work
	s.baseCxt = BaseCxt
	s.SetContext()
	return s
}

func(s *WeiBoService) Work() {
	link := "https://s.weibo.com/top/summary"
	req := utils.NewBrowser()
	cookiesStr := `login_sid_t=01863878c95b8752a100a977f5a244b3; cross_origin_proto=SSL; _s_tentry=weibo.com; Apache=381481134740.3935.1633676797861; SINAGLOBAL=381481134740.3935.1633676797861; ULV=1633676797870:1:1:1:381481134740.3935.1633676797861:; SUB=_2AkMWPHUof8NxqwJRmPgXymLqaYVxwg_EieKgYITzJRMxHRl-yT9jqlMLtRB6PbxbxzzqZI9Bt9XA5iyIdD3_p-Aur3z5; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WWl-0zBC5s.nORGVdwwK1qG; WBStorage=6ff1c79b|undefined`
	cookiesList := strings.Split(cookiesStr, ";")
	cookieMap := make(map[string]string)
	for _, v := range cookiesList {
		cvArr := strings.Split(v, "=")
		k := strings.TrimSpace(cvArr[0])
		cookieMap[k] = url.QueryEscape(cvArr[1])
	}

	req.AddCookie(cookieMap)
	str, status  := req.Get(link)
	if len(str) == 0 || status != 200 {
		fmt.Println(string(str))
		log.Info("微博网站采集失败" + strconv.Itoa(status))
		return
	}

	htmlStr := string(str)
	r, _:= regexp.Compile("(?ism)<td\\s*class=\"td-01\\s*ranktop\">(\\d+)</td>\\s*<td class=\"td-02\">\\s*<a href=\"(.*?)\" target=\"_blank\">(.*?)</a>\\s*<span>\\s*(\\d+)</span>")
	matches := r.FindAllStringSubmatch(htmlStr, -1)
	if len(matches) == 0 {
		log.Info("微博网站没有匹配到信息")
		return
	}

	var updateTotal, addTotal int

	for _, match := range matches {
		var hotTop models.HotTop
		hotTop.Site = "新浪微博"
		hotTop.Title = match[3]
		hotTop.Link = "https://s.weibo.com" + match[2]

		num, _ := strconv.ParseInt(match[4], 10, 64)
		hotTop.Num = num

		sortNum, _ := strconv.Atoi(match[1])
		hotTop.Sort = sortNum
		updateNum, addNum :=  hotTop.AddOrUpdateHotTopNew(hotTop)

		updateTotal += updateNum
		addTotal += addNum

	}

	log.Infof(s.Desc + "采集成功,更新%d, 新增%d", updateTotal, addTotal)


}





