package controller

import (
	"admin/db"
	"admin/middleware"
	"admin/models"
	"admin/setting"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RoleController struct {
	baseController
}

func (rl *RoleController) Routers(r *gin.RouterGroup) {
	g := r.Group("/role")
	g.Use(middleware.JWT())
	{
		g.GET("/list", rl.list )
		g.POST("/update", rl.update )
		g.POST("/create", rl.create )
	}
}

func (rl *RoleController) list(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	status := c.DefaultQuery("status", "")
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(setting.PageSize)))

	roleInfo := &models.AdminRole{}
	total, roleList := roleInfo.GetList(username, status, page, limit)

	rl.SuccessResponse(c, "获取成功", map[string]interface{}{
		"total" : total,
		"items" : roleList,
	})


}


func (rl *RoleController) update(c *gin.Context) {
	var roleWrap models.AdminRole
	if err := c.ShouldBindJSON(&roleWrap); err != nil {
		rl.ErrorResponse(c, err.Error(), nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(roleWrap.Id, "id").Message("角色id非法")
	valid.Required(roleWrap.Name, "name").Message("角色名称非法")
	valid.Required(roleWrap.Remark, "remark").Message("备注不能为空")
	valid.Range(roleWrap.Status, 1, 2, "status").Message("状态非法")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			rl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	var role models.AdminRole
	roleInfo := role.GetRoleInfoByUsername(roleWrap.Name)
	if roleInfo.Id > 0 && roleInfo.Id != roleWrap.Id {
		rl.ErrorResponse(c, "角色名称已经存在", nil)
		return
	}

	adminMap := make(map[string]interface{})
	adminMap["name"] = roleWrap.Name
	adminMap["remark"] = roleWrap.Remark
	adminMap["status"] = roleWrap.Status
	err := db.GetInstance().GetMysqlDb().Model(&role).Where("id=?", roleWrap.Id).Updates(adminMap).Error
	if err != nil {
		rl.ErrorResponse(c, "更新角色失败", nil)
		return
	}

	rl.SuccessResponse(c, "更新角色成功", nil)
}

func (rl *RoleController) create(c *gin.Context) {
	var roleWrap models.AdminRole
	if err := c.ShouldBindJSON(&roleWrap); err != nil {
		rl.ErrorResponse(c, err.Error(), nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(roleWrap.Name, "name").Message("角色名称非法")
	valid.Required(roleWrap.Remark, "remark").Message("备注不能为空")
	valid.Range(roleWrap.Status, 1, 2, "status").Message("状态非法")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			rl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	var role models.AdminRole
	roleInfo := role.GetRoleInfoByUsername(roleWrap.Name)
	if roleInfo.Id > 0 {
		rl.ErrorResponse(c, "角色名称已经存在", nil)
		return
	}

	role.Status = roleWrap.Status
	role.Remark = roleWrap.Remark
	role.Name = roleWrap.Name

	err := db.GetInstance().GetMysqlDb().Create(&role).Error
	if err != nil {
		rl.ErrorResponse(c, "添加角色失败", nil)
		return
	}

	rl.SuccessResponse(c, "添加角色成功", nil)
}







