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

type RoleUserController struct {
	baseController
}

func (ruc *RoleUserController) Routers(r *gin.RouterGroup) {
	g := r.Group("/role_user")
	g.Use(middleware.JWT())
	{
		g.GET("/list", ruc.list )
		g.POST("/create", ruc.create )
	}
}

func (ruc *RoleUserController) list(c *gin.Context) {
	roleId := 0
	rid := strings.TrimSpace(c.Query("role_id"))
	if len(rid) > 0 {
		rid, err := strconv.Atoi(rid)
		if err != nil {
			ruc.ErrorResponse(c, "参数错误", nil)
			return
		}
		roleId = rid
	}

	//得到所有用户
	var adminList []models.Admins
	dbCon := db.GetInstance().GetMysqlDb()
	dbCon.Model(models.Admins{}).Where("status=?", 1).Find(&adminList)

	//得到所有角色
	var roleList []models.AdminRole
	roleOptionList := make([]interface{}, 0)
	dbCon.Model(models.AdminRole{}).Where("status=?", 1).Find(&roleList)
	if len(roleList) > 0 {
		for _, v := range roleList {
			roleOptionList = append(roleOptionList, map[string]interface{}{"key": v.Id, "label": v.Name})
		}
	}

	//得到有权限的角色
	var roleUser models.AdminRoleUser
	roleUserList := roleUser.GetUserInfoListByRoleId(roleId)
	roleUserMap := make(map[int]models.AdminRoleUser)
	for _, roleUser := range roleUserList {
		roleUserMap[roleUser.UserId] = roleUser
	}

	userMapList := make([]map[string]interface{}, 0, len(roleUserList))
	for _, admin := range adminList {
		var isChecked = false
		if _, ok := roleUserMap[admin.Id]; ok {
			isChecked = true
		}
		userMap := make(map[string]interface{})
		userMap["id"]  = admin.Id
		userMap["pId"] = 0
		userMap["name"] = admin.Name
		userMap["checked"] = isChecked
		userMapList = append(userMapList, userMap)
	}

	ruc.SuccessResponse(c, "参数错误", map[string]interface{}{
		"tree": userMapList,
		"roleOptionList" : roleOptionList,
	})
}

func (ruc *RoleUserController) create(c *gin.Context) {
	var roleId, err = strconv.Atoi(c.PostForm("role_id"))
	if err != nil || roleId <= 0 {
		ruc.ErrorResponse(c, "请选择角色", nil)
		return
	}

	var uidsStr = strings.TrimSpace(c.PostForm("uids_str"))
	if len(uidsStr) == 0 {
		ruc.ErrorResponse(c, "请选择用户", nil)
		return
	}

	userIdList := strings.Split(uidsStr, ",")

	//删除原来的对应关系
	dbCon := db.GetInstance().GetMysqlDb()
	if err := dbCon.Where("role_id=?", roleId).Delete(&models.AdminRoleUser{}).Error; err != nil {
		ruc.ErrorResponse(c, "操作失败", nil)
		return
	}

	for _, uid := range userIdList {
		uid, _ := strconv.Atoi(uid)
		var roleUser models.AdminRoleUser
		roleUser.RoleId = roleId
		roleUser.UserId = uid
		roleUser.Status = 1
		dbCon.Create(&roleUser)
	}

	ruc.SuccessResponse(c, "操作成功", nil)
	return
}
