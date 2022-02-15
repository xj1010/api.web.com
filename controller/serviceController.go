package controller

import (
	"admin/console"
	"admin/middleware"
	"admin/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
	"strconv"
)

type ServiceController struct {
	baseController
}

func (scl *ServiceController) Routers(r *gin.RouterGroup) {
	g := r.Group("/service")
	g.Use(middleware.JWT())
	{
		g.GET("/list", scl.list )
		g.GET("/update", scl.update )
	}
}

func (scl *ServiceController) update(c *gin.Context) {
	copType := c.Query("type")
	copTypeMap :=  map[string]func(c *gin.Context) {
		"start" : scl.start,
		"stop" : scl.stop,
		"allstart" : scl.allStart,
		"allstop" : scl.allStop,
		"restart" : scl.allStop,
	}

	if m, ok := copTypeMap[copType]; !ok {
		scl.ErrorResponse(c, "非法操作",  nil)
		return
	} else {
		m(c)
	}
}

func (scl *ServiceController) list(c *gin.Context) {
	serviceMap := console.GetServiceMap()
	fmt.Println(serviceMap)
	title := c.Query("title")
	status := c.Query("status")
	isStart, _ := strconv.ParseBool(status)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(setting.PageSize)))
	serviceList := make([]map[string]interface{}, 0, len(serviceMap))

	var serviceNames []string
	for key, _ := range serviceMap {
		serviceNames = append(serviceNames, key)
	}

	sort.Strings(serviceNames)
	for _, serviceName := range serviceNames {
		s := serviceMap[serviceName]
		if title != ""  && serviceName != title {
			continue
		}
		if status != "" && isStart != s.IsRunning{
			continue
		}

		serviceMap := make(map[string]interface{})
		serviceMap["name"] = serviceName
		serviceMap["status"] = s.IsRunning
		serviceMap["desc"] = s.Desc
		serviceList = append(serviceList, serviceMap)
	}

	total := len(serviceList)
	if total > 0 {
		fPage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(total) / float64(limit)), 64)
		maxPage := int(math.Ceil(fPage))
		if page > maxPage {
			page = maxPage
		}
		offsetStart := (page-1)*limit
		offsetEnd := offsetStart+ limit
		if offsetEnd > total {
			offsetEnd = total
		}
		serviceList =  serviceList[offsetStart:offsetEnd]
	}

	scl.SuccessResponse(c, "获取成功", map[string]interface{}{
		"total": total,
		"items": serviceList,
	})
}

func (scl *ServiceController) start(c *gin.Context) {
	serviceName := c.Query("s")
	serviceMap := console.GetServiceMap()
	if s, ok := serviceMap[serviceName]; ok {
		s.Start()
		scl.SuccessResponse(c, "操作成功",  nil)
		return
	}
	scl.ErrorResponse(c, "参数非法",  nil)
}

func (scl *ServiceController) allStart(c *gin.Context) {
	for _ , s := range console.GetServiceMap() {
		s.Start()
	}
	scl.SuccessResponse(c, "操作成功",  nil)
}

func (scl *ServiceController) allStop(c *gin.Context) {
	for _ , s := range console.GetServiceMap() {
		s.Stop()
	}
	scl.SuccessResponse(c, "操作成功",  nil)
}

func (scl *ServiceController) stop(c *gin.Context) {
	serviceName := c.Query("s")
	serviceMap := console.GetServiceMap()

	if s, ok := serviceMap[serviceName]; ok {
		s.Stop()
		scl.SuccessResponse(c, "操作成功",  nil)
		return
	}
	scl.ErrorResponse(c, "参数非法",  nil)
}

func (scl *ServiceController) restart(c *gin.Context) {
	serviceName := c.Query("s")
	serviceMap := console.GetServiceMap()

	if s, ok := serviceMap[serviceName]; ok {
		s.Stop()
		s.Start()
		scl.SuccessResponse(c, "操作成功",  nil)
		return
	}
	scl.ErrorResponse(c, "参数非法",  nil)
}
