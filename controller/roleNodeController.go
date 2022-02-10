package controller

import (
	"admin/db"
	"admin/middleware"
	"admin/models"
	_ "encoding/json"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type RoleNodeController struct {
	baseController
}

func (rnc *RoleNodeController) Routers(r *gin.RouterGroup) {
	g := r.Group("/role_node")
	g.Use(middleware.JWT())
	{
		g.GET("/list", rnc.list )
		g.POST("/create", rnc.create )
	}
}

func (rnc *RoleNodeController) list(c *gin.Context) {
	roleId := 0
	rid := strings.TrimSpace(c.Query("role_id"))
	if len(rid) > 0 {
		rid, err := strconv.Atoi(rid)
		if err != nil {
			rnc.ErrorResponse(c, "参数错误", nil)
			return
		}
		roleId = rid
	}

	//得到所有节点
	var nodeList []models.AdminNode
	dbCon := db.GetInstance().GetMysqlDb()
	dbCon.Model(models.AdminNode{}).Where("status=?", 1).Find(&nodeList)
	if len(nodeList) == 0 {
		rnc.ErrorResponse(c, "没有添加节点", nil)
	}

	//得到所有角色
	var roleList []models.AdminRole
	roleOptionList := make([]interface{}, 0)
	dbCon.Model(models.AdminRole{}).Where("status=?", 1).Find(&roleList)
	if len(roleList) > 0 {
		for _, v := range roleList {
			roleOptionList = append(roleOptionList, map[string]interface{}{"key": v.Id, "label": v.Name})
		}
	}

	//得到有角色的节点
	var roleNode models.AdminRoleNode
	roleNodeList := roleNode.GetNodeListByRoleId(roleId)
	roleUserMap := make(map[int]models.AdminRoleNode)
	for _, roleNode := range roleNodeList {
		roleUserMap[roleNode.NodeId] = roleNode
	}

	nodeMapList := make([]map[string]interface{}, 0, len(nodeList))
	for _, node := range nodeList {
		var isChecked = false
		if _, ok := roleUserMap[node.Id]; ok {
			isChecked = true
		}
		nodeMap := make(map[string]interface{})
		nodeMap["id"]  = node.Id
		nodeMap["pId"] = node.Pid
		nodeMap["name"] = node.Name
		nodeMap["checked"] = isChecked
		nodeMap["open"] = true
		nodeMapList = append(nodeMapList, nodeMap)
	}

	rnc.SuccessResponse(c, "参数错误", map[string]interface{}{
		"tree": nodeMapList,
		"roleOptionList" : roleOptionList,
		"roleId": roleId,
	})
}

func (rnc *RoleNodeController) create(c *gin.Context) {
	var roleId, err = strconv.Atoi(c.PostForm("role_id"))
	if err != nil || roleId <= 0 {
		rnc.ErrorResponse(c, "请选择角色", nil)
		return
	}

	var nodeIds = strings.TrimSpace(c.PostForm("node_ids_str"))
	if len(nodeIds) == 0 {
		rnc.ErrorResponse(c, "请选择节点", nil)
		return
	}

	//删除原来的对应关系
	dbCon := db.GetInstance().GetMysqlDb()
	if err := dbCon.Where("role_id=?", roleId).Delete(&models.AdminRoleNode{}).Error; err != nil {
		rnc.ErrorResponse(c, "操作失败", nil)
		return
	}

	nodeIdList := strings.Split(nodeIds, ",")
	for _, nodeId := range nodeIdList {
		nodeId, _ := strconv.Atoi(nodeId)
		var roleNode models.AdminRoleNode
		roleNode.RoleId = roleId
		roleNode.NodeId = nodeId
		roleNode.Status = 1
		dbCon.Create(&roleNode)
	}

	rnc.SuccessResponse(c, "操作成功", nil)
	return
}
