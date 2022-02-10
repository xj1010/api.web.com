package controller

import (
	"admin/db"
	"admin/middleware"
	"admin/models"
	"admin/setting"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NodeController struct {
	baseController
}

func (ncl *NodeController) Routers(r *gin.RouterGroup) {
	g := r.Group("/node")
	g.Use(middleware.JWT())
	{
		g.GET("/list", ncl.list )
		g.POST("/update", ncl.update )
		g.POST("/create", ncl.create )
		g.GET("/delete", ncl.delete )
	}
}

func (ncl *NodeController) delete(c *gin.Context) {
	id :=  c.DefaultQuery("id", "")
	level :=  c.DefaultQuery("level", "")
	if id == ""  || level == ""{
		ncl.ErrorResponse(c, "参数非法", nil)
		return
	}

	nId, err := strconv.Atoi(id)
	if err != nil {
		ncl.ErrorResponse(c, "参数非法", nil)
		return
	}

	nLevel, err := strconv.Atoi(level)
	if err != nil {
		ncl.ErrorResponse(c, "参数非法", nil)
		return
	}

	var pidNodes []int
	pidNodes = append(pidNodes, nId)

	dbCon := db.GetInstance().GetMysqlDb()
	if nLevel == 1 {
		var secondNodes []models.AdminNode
		dbCon.Model(models.AdminNode{}).Where("pid = ?", nId).Select("id").Find(&secondNodes)

		if len(secondNodes) > 0 {
			for _, v := range secondNodes {
				pidNodes = append(pidNodes, v.Id)
			}

			var threeNodes []models.AdminNode
			dbCon.Model(models.AdminNode{}).Where("pid in (?) ", secondNodes).Select("id").Find(&threeNodes)
			if len(threeNodes) > 0 {
				for _, v := range threeNodes {
					pidNodes = append(pidNodes, v.Id)
				}
			}
		}
	}

	if nLevel == 2 {
		var threeNodes []models.AdminNode
		dbCon.Model(models.AdminNode{}).Where("pid = ? ", nId).Select("id").Find(&threeNodes)
		if len(threeNodes) > 0 {
			for _, v := range threeNodes {
				pidNodes = append(pidNodes, v.Id)
			}
		}
	}

	if err := dbCon.Delete(models.AdminNode{}, "id in (?)", pidNodes).Error; err != nil {
		ncl.ErrorResponse(c, "删除失败", nil)
		return
	}

	ncl.SuccessResponse(c, "删除成功", nil)
}

func (ncl *NodeController) list(c *gin.Context) {
	searchMap := make(map[string]string)
	searchMap["pid"] = c.DefaultQuery("pid", "-1")
	searchMap["name"] = c.DefaultQuery("name", "")
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(setting.PageSize)))
	isExport, _ := strconv.Atoi(c.DefaultQuery("export", "0"))
	searchMap["sort"] = c.DefaultQuery("sort", "")

	if isExport == 1 {
		page = 1
		setting.PageSize = 1000
	}

	var adminNode models.AdminNode
	total, nodeList := adminNode.GetList(searchMap, page, limit)

	//获取所有节点信息
	treeNodeList := make([]models.TreeNode, 0)
	if isExport != 1 {
		var allNodeList []models.AdminNode
		db.GetInstance().GetMysqlDb().Order("created_at desc").Find(&allNodeList)

		if len(allNodeList) > 0 {
			for _, n := range allNodeList {
				var node models.TreeNode
				node.Id = n.Id
				node.Pid = n.Pid
				node.Name = n.Name
				treeNodeList = append(treeNodeList, node)
			}
		}
	}

    ncl.SuccessResponse(c, "获取成功", map[string]interface{}{
		"total" : total,
		"items" : nodeList,
		"tree" : treeNodeList,
	})
}

func (ncl *NodeController) create(c *gin.Context) {
	var adminNodeWrap models.AdminNode
	if err := c.ShouldBindJSON(&adminNodeWrap); err != nil {
		fmt.Println(err)
		ncl.ErrorResponse(c, err.Error(), nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(adminNodeWrap.Name, "name").Message("节点名称不能为空")
	valid.Required(adminNodeWrap.Path, "path").Message("节点路径不能为空")
	valid.Required(adminNodeWrap.PermissionName, "permission_name").Message("权限标识不能为空")
	valid.Range(adminNodeWrap.Level, 1, 3, "level").Message("节点类型非法")
	valid.Range(adminNodeWrap.Sort, 1, 255, "sort").Message("权重非法")
	valid.Range(adminNodeWrap.Status, 1, 2, "status").Message("状态非法")

	if adminNodeWrap.Level != 1 {
		valid.Required(adminNodeWrap.Pid, "pid").Message("父节点不能为空")
	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			ncl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	if adminNodeWrap.Level < 0 || (adminNodeWrap.Level == 1 && adminNodeWrap.Pid != 0) {
		ncl.ErrorResponse(c, "父节点非法", nil)
		return
	}

	dbCon := db.GetInstance().GetMysqlDb()
	dbCon.Where("path=?", adminNodeWrap.Path).Find(&adminNodeWrap)
	if adminNodeWrap.Id > 0 {
		ncl.ErrorResponse(c, "该路径已经存在", nil)
		return
	}

	userInfo := ncl.GetUserInfo(c)

	adminNodeWrap.AdminName = userInfo.Name
	if err := dbCon.Create(&adminNodeWrap).Error; err != nil {
		ncl.ErrorResponse(c, "添加失败", nil)
		return
	}

	ncl.SuccessResponse(c, "添加成功", nil)
}

func (ncl *NodeController) update(c *gin.Context) {
	var adminNode models.AdminNode
	if err := c.ShouldBindJSON(&adminNode); err != nil {
		ncl.ErrorResponse(c, err.Error(), nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(adminNode.Name, "name").Message("节点名称不能为空")
	valid.Required(adminNode.Id, "id").Message("参数非法")
	valid.Required(adminNode.Path, "path").Message("节点路径不能为空")
	valid.Range(adminNode.Sort, 0, 255, "status").Message("权重非法")
	valid.Range(adminNode.Status, 1, 2, "status").Message("状态非法")

	if adminNode.Level != 1 {
		valid.Required(adminNode.Pid, "pid").Message("父节点不能为空")
	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			ncl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	if adminNode.Level < 0 || (adminNode.Level == 1 && adminNode.Pid != 0) {
		ncl.ErrorResponse(c, "父节点非法", nil)
		return
	}

	var dbAdminNode models.AdminNode
	dbCon := db.GetInstance().GetMysqlDb()
	dbCon.Where("path=?", adminNode.Path).Find(&dbAdminNode)
	if dbAdminNode.Id > 0 && (adminNode.Id != dbAdminNode.Id )  {
		ncl.ErrorResponse(c, "该路径已经存在", nil)
		return
	}

	nodeMap := make(map[string]interface{})
	nodeMap["name"] = adminNode.Name
	nodeMap["path"] = adminNode.Path
	nodeMap["icon"] = adminNode.Icon
	nodeMap["status"] = adminNode.Status
	nodeMap["sort"] = adminNode.Sort
	nodeMap["permission_name"] = adminNode.PermissionName

	if err := dbCon.Model(models.AdminNode{}).Where("id=?", adminNode.Id).Updates(nodeMap).Error; err != nil {
		ncl.ErrorResponse(c, "更新失败", nil)
		return
	}

	ncl.SuccessResponse(c, "更新成功", nil)
}






