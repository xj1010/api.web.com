package controller

import (
	"admin/middleware"
	"admin/models"
	"admin/setting"
	"github.com/gin-gonic/gin"
	"strconv"
)

type HotController struct {
	baseController
}

func (h *HotController) Routers(r *gin.RouterGroup) {
	g := r.Group("/hot")
	g.Use(middleware.JWT())
	{
		g.GET("/list", h.list )
	}
}

func (h *HotController) list(c *gin.Context) {
	site := c.DefaultQuery("site", "")
	title := c.DefaultQuery("title", "")
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(setting.PageSize)))

	searchMap := make(map[string]string)
	searchMap["site"] = site
	searchMap["title"] = title
	searchMap["startDate"] = startDate
	searchMap["endDate"] = endDate

	hopInfo := &models.HotTop{}
	total, topList := hopInfo.GetList(searchMap, page, limit)

	h.SuccessResponse(c, "获取成功", map[string]interface{}{
		"total": total,
		"items": topList,
	})
}








