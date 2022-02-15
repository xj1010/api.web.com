package service

import (
	"admin/models"
	_ "admin/models"
	"admin/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	uuid "github.com/satori/go.uuid"
	"github.com/siddontang/go/log"
	"math/rand"
	"net/url"
	_ "regexp"
	"strconv"
	_ "strconv"
	"strings"
	"sync"
	"time"
	_ "time"
)

type LaGouService struct {
	BaseService
}

func NewLaGouService(BaseCxt context.Context) *LaGouService {
	s := &LaGouService{}
	s.IsRunning = false
	s.Desc = "拉勾网"
	s.ServiceID = "LaGou"
	s.Task = s.Work
	s.baseCxt = BaseCxt
	s.SetContext()
	return s
}

func(s *LaGouService) Work() {
	 subjectChan := make(chan string ,30)
	 var wg sync.WaitGroup
	 //得到科目
	 wg.Add(1)
	 go s.getSubjectList(subjectChan, &wg)

	 for i := 0; i < 3; i++ {
	 	 wg.Add(1)
		 go s.getSubjectDetail(subjectChan, &wg)
	 }

	 wg.Wait()
}

func(s *LaGouService) getSubjectDetail(subjectChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if subjectName, ok := <- subjectChan; ok && s.IsRunning  {
			var addTotal, updateTotal int
			for page := 1; page <= 1000; page++ {
				if !s.IsRunning {
					break
				}
			    addNum , updateNum := s.getSubjectDataByName(subjectName, page)
				rand.Seed(time.Now().Unix())
				time.Sleep(time.Duration(rand.Intn(2)+1)*time.Second)
			    if addNum == -1 && updateNum == -1 {
					break
			    }
			    addTotal += addNum
			    updateTotal += updateNum
			}

			fmt.Printf("岗位:%s成功添加数据%d条，更新数据%d条\n", subjectName, addTotal, updateTotal)
		} else {
			break
		}
	}
}

func(s *LaGouService) getSubjectList(subjectChan chan string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(subjectChan)
	}()
	link := "https://www.lagou.com"
	req := utils.NewBrowser()
	str, status := req.Get(link)
	if len(str) == 0 || status != 200 {
		log.Info("拉钩网站采集失败")
		return
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(str)))
	if err != nil {
		log.Info(err)
		return
	}

	subjectMap := make(map[string]string, 10);
	dom.Find(".menu_box:first-child a").Each(func(i int, selection *goquery.Selection) {
		name , err := selection.Find("h3").Html()
		href , flag := selection.Attr("href")
		if err == nil && flag {
			subjectMap[name] = href
		}
	})

	if len(subjectMap) == 0 {
		log.Info("采集语言失败")
		return
	}

	for name, _ := range subjectMap {
		subjectChan <- name
	}
}

func(s *LaGouService) getSubjectDataByName(subjectName string, page int, ) (addTotal, updateTotal int) {
	link := "https://www.lagou.com/wn/jobs?cl=false&fromSearch=true&kd=" + subjectName
	req := utils.NewBrowser()
	str, status, redirectMap := req.GetRedirect(link)
	if len(str) == 0 || status != 200 {
		log.Info("拉钩网站采集失败")
		return addTotal, updateTotal
	}

	cookiesStr := `JSESSIONID=ABAAAECABFAACEA2005649916AD1471E6593E7166381BBF; WEBTJ-ID=20220122152129-17e80a95b2d2f0-012fc7738f09b8-131a6a50-1049088-17e80a95b2e12c; RECOMMEND_TIP=true; PRE_UTM=; PRE_HOST=; PRE_LAND=https://www.lagou.com/; user_trace_token=20220122152134-b196f5d7-85b0-4673-a6c4-e12fdd2ebb0c; LGSID=20220122152134-fd868c17-868d-47c2-abdc-5212e5f593cd; PRE_SITE=; LGUID=20220122152134-3660c4fa-091e-4113-8ff6-abbd77f22c2d; _ga=GA1.2.535061443.1642836091; _gat=1; Hm_lvt_4233e74dff0ae5bd0a3d81c6ccf756e6=1642836091; privacyPolicyPopup=false; _gid=GA1.2.88997650.1642836091; sajssdk_2015_cross_new_user=1; sensorsdata2015session={}; gate_login_token=ebab73f2632b4e50502e40140f7a485eb601981b29bf26b2f274103169d5ed79; LG_LOGIN_USER_ID=c3824c21bcfd381aae778a1f7f65ad37cc16595ae5aa1062322f36377751caee; LG_HAS_LOGIN=1; _putrc=696BFBB6640D2423123F89F2B170EADC; login=true; unick=熊金; showExpriedIndex=1; showExpriedCompanyHome=1; showExpriedMyPublish=1; hasDeliver=8; index_location_city=武汉; X_HTTP_TOKEN=d29a8af91c63cceb931638246181d9c2687bd94fb7; __SAFETY_CLOSE_TIME__11585367=1; Hm_lpvt_4233e74dff0ae5bd0a3d81c6ccf756e6=1642836137; sensorsdata2015jssdkcross={"distinct_id":"11585367","first_id":"17e80a96178a9-02a2e8b14b55af-131a6a50-1049088-17e80a9617923d","props":{"$latest_traffic_source_type":"直接流量","$latest_search_keyword":"未取到值_直接打开","$latest_referrer":"","$os":"Windows","$browser":"Chrome","$browser_version":"97.0.4692.71","lagou_company_id":""},"$device_id":"17e80a96178a9-02a2e8b14b55af-131a6a50-1049088-17e80a9617923d"}; LGRID=20220122152228-48dd63da-facd-449c-8632-849a0258f7c5`
	cookiesList := strings.Split(cookiesStr, ";")
	cookieMap := make(map[string]string)
	for _, v := range cookiesList {
		cvArr := strings.Split(v, "=")
		k := strings.TrimSpace(cvArr[0])
		timeNow := time.Now()
		if k == "LGRID" {
			cookieMap[k] = fmt.Sprintf("%04d%02d%02d%02d%02d%02d-%s", timeNow.Year(), timeNow.Month(), timeNow.Day(),timeNow.Hour(),timeNow.Minute(), timeNow.Second(), uuid.NewV4())
			continue
		}

		if k == "Hm_lpvt_4233e74dff0ae5bd0a3d81c6ccf756e6" {
			cookieMap[k] = strconv.FormatInt(timeNow.Unix(), 10)
			continue
		}
		if k == "SEARCH_ID" {
			cookieMap[k] = utils.Md5String(subjectName)
			continue
		}
		cookieMap[k] = url.QueryEscape(cvArr[1])
	}
	for k, v := range redirectMap {
		k = strings.TrimSpace(k)
		cookieMap[k] = url.QueryEscape(v)
	}

	userAgent := " Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36"
	link = strings.Replace(link, subjectName, url.QueryEscape(subjectName), -1)
	var headerMap = map[string]string{
		"Referer" : link,
		"user-agent" : userAgent,
		"x-requested-with" : "XMLHttpRequest",
		"sec-ch-ua": `"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"`,
	}

	req.AddHeader(headerMap)
	req.AddCookie(cookieMap)

	subjectUrl := "https://www.lagou.com/jobs/v2/positionAjax.json"
	first := "false"
	if page == 1 {
		first = "true"
	}
	paramsMap := map[string]string{
		"first" : first,
		"needAddtionalResult" : "false",
		"city" : "武汉",
		"cl" : "false",
		"pn" : strconv.Itoa(page),
		"kd" : subjectName,
		"fromSearch" : "true",
		"px" : "new",
	}

	str = req.Post(subjectUrl, paramsMap)

	if len(str) == 0 {
		fmt.Printf("请求接口失败")
		return addTotal, updateTotal
	}

	var positionMap models.PositionMap
	err := json.Unmarshal(str, &positionMap)
	if err != nil {
		fmt.Println(err)
		return addTotal, updateTotal
	}

	fmt.Println(page, "==",string(str))
	//如果已经获取不到数据了
	if len(positionMap.Content.Hrinfomap) == 0 {
		return -1, -1
	}

	if positionMap.Content.Positionresult.Resultsize > 0 {
		for _, v := range positionMap.Content.Positionresult.Result {
			lagouObj := models.Lagou{}
			lagouObj.CompanySize = v.Companysize
			lagouObj.CompanyId = v.Companyid
			lagouObj.CompanyName = v.Companyfullname
			lagouObj.PositionId = v.Positionid
			lagouObj.PositionName = v.Positionname
			lagouObj.FirstType = v.Firsttype
			lagouObj.SecondType = v.Secondtype
			lagouObj.ThirdType = v.Thirdtype
			lagouObj.City = v.City
			lagouObj.SalaryMin = 0
			lagouObj.SalaryMax = 0

			if len(v.Salary) > 0 && strings.Contains(v.Salary, "-") {
				salaryArr := strings.Split(v.Salary, "-")
				salaryMin, _ := strconv.Atoi(strings.Replace(salaryArr[0], "k", "000", -1))
				salaryMax, _ := strconv.Atoi(strings.Replace(salaryArr[1], "k", "000", -1))
				lagouObj.SalaryMin = salaryMin
				lagouObj.SalaryMax = salaryMax
			}
			salaryMonth, _ := strconv.Atoi(v.Salarymonth)
			lagouObj.SalaryMonth = salaryMonth

			lagouObj.WorkYear = v.Workyear
			lagouObj.Education = v.Education
			createTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Createtime, time.Local)
			lagouObj.Createtime = createTime

			updateNum , addNum := lagouObj.AddOrUpdateRecord(lagouObj)
			if updateNum > 0 {
				updateTotal += updateNum
			}
			if addNum > 0 {
				addTotal += addNum
			}
		}
	}

	return addTotal, updateTotal

}





