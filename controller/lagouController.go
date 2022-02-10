package controller

import (
	"admin/middleware"
	"admin/models"
	"admin/setting"
	"github.com/gin-gonic/gin"
	"strconv"
)

type LagouController struct {
	baseController
}

func (h *LagouController) Routers(r *gin.RouterGroup) {
	g := r.Group("/lagou")
	g.Use(middleware.JWT())
	{
		g.GET("/list", h.list )
	}
}

func (h *LagouController) list(c *gin.Context) {
	site := c.DefaultQuery("site", "")
	searchType := c.DefaultQuery("searchType", "")
	keyWord := c.DefaultQuery("keyword", "")
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(setting.PageSize)))

	searchMap := make(map[string]string)
	searchMap["site"] = site
	searchMap["searchType"] = searchType
	searchMap["keyword"] = keyWord
	searchMap["startDate"] = startDate
	searchMap["endDate"] = endDate

	laGouInfo := &models.Lagou{}
	total, topList := laGouInfo.GetList(searchMap, page, limit)

	h.SuccessResponse(c, "获取成功", map[string]interface{}{
		"total": total,
		"items": topList,
	})
}








