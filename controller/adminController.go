package controller

import (
	"admin/db"
	"admin/middleware"
	"admin/models"
	"admin/setting"
	"admin/utils"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type AdminController struct {
	baseController
}

func (acl *AdminController) Routers(r *gin.RouterGroup) {
	g := r.Group("/admin")
	g.Use(middleware.JWT())
	{
		g.GET("/list", acl.list)
		g.POST("/update", acl.update)
		g.POST("/create", acl.create)
		g.GET("/info", acl.info)
		g.GET("/node", acl.node)
	}
}


func (acl *AdminController) node(c *gin.Context) {
	var adminNode models.AdminNode
	token := c.Request.Header.Get("X-Token")
	uid := models.UserAuthSession[token].AdminInfo.Id
	nodeList := adminNode.GetFormatUserNodeAuthList(uid)
	acl.SuccessResponse(c, "获取成功", map[string]interface{}{
		"node":     nodeList,
	})
}

func (acl *AdminController) info(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")
	adminInfo := models.UserAuthSession[token].AdminInfo
	avatar := "https://wpimg.wallstcn.com/57ed425a-c71e-4201-9428-68760c0537c4.jpg"
	if len(adminInfo.ImageUrl) > 0 {
		if setting.HTTPPort != 80 {
			avatar = setting.WebUrl + ":" + strconv.Itoa(setting.HTTPPort) +"/" + adminInfo.ImageUrl
		} else {
			avatar = setting.WebUrl +"/" + adminInfo.ImageUrl
		}
	}

	acl.SuccessResponse(c, "获取成功", map[string]interface{}{
		"avatar":       avatar,
		"introduction": "搬砖人",
		"roles":        []string{"edit"},
	})
}

func (acl *AdminController) list(c *gin.Context) {
	status := c.DefaultQuery("status", "-1")
	username := c.DefaultQuery("username", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(setting.PageSize)))


	searchMap := make(map[string]string)
	searchMap["status"] = status
	searchMap["username"] = strings.TrimSpace(username)

	adminInfo := &models.Admins{}
	total, adminList := adminInfo.GetList(searchMap, page, limit)

	adminWrapList := make([]models.AdminsWrap, 0)
	for _, admin := range adminList {
		var adminWrap models.AdminsWrap
		adminWrap.Id = admin.Id
		adminWrap.Name = admin.Name
		adminWrap.Status = admin.Status
		adminWrap.Email = admin.Email
		adminWrap.Mobile = admin.Mobile
		adminWrap.CreatedAt = admin.CreatedAt
		adminWrap.UpdatedAt = admin.UpdatedAt
		adminWrap.LoginNum = admin.LoginNum
		if len(admin.ImageUrl) > 0 {
			adminWrap.ImageUrl = "http://localhost:8089/" + admin.ImageUrl
			adminWrap.ImagePath = admin.ImageUrl
		}

		adminWrapList = append(adminWrapList, adminWrap)
	}

	acl.SuccessResponse(c, "获取成功", map[string]interface{}{
		"total": total,
		"items": adminWrapList,
	})
}

func (acl *AdminController) update(c *gin.Context) {
	var adminsWrap models.AdminsWrap
	if err := c.ShouldBindJSON(&adminsWrap); err != nil {
		acl.ErrorResponse(c, err.Error(), nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(adminsWrap.Email, "email").Message("邮箱不能为空")
	valid.Required(adminsWrap.Id, "id").Message("用户id非法")
	valid.Email(adminsWrap.Email, "email").Message("邮箱非法")
	valid.Required(adminsWrap.Mobile, "mobile").Message("手机号不能为空")
	valid.Mobile(adminsWrap.Mobile, "mobile").Message("手机号非法")
	valid.Range(adminsWrap.Status, 1, 2, "status").Message("状态非法")
	if len(adminsWrap.Password) > 0 {
		valid.MinSize(adminsWrap.Password, 6, "password").Message("密码最少6位数")
		valid.Required(adminsWrap.ConfirmPassword, "confirmPassword").Message("确认密码不能为空")
	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			acl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	if len(adminsWrap.Password) > 0 {
		if adminsWrap.Password != adminsWrap.ConfirmPassword {
			acl.ErrorResponse(c, "两次输入的密码不一致", nil)
			return
		}
	}

	var admins models.Admins
	adminInfo := admins.GetAdminInfoByUsername(adminsWrap.Name)
	if adminInfo.Id > 0 && adminsWrap.Id != adminInfo.Id {
		acl.ErrorResponse(c, "该用户名已经存在", nil)
		return
	}

	adminMap := make(map[string]interface{})
	adminMap["email"] = adminsWrap.Email
	adminMap["mobile"] = adminsWrap.Mobile
	adminMap["status"] = adminsWrap.Status
	if len(adminsWrap.ImagePath) > 0 {
		adminMap["image_url"] = adminsWrap.ImagePath
	}

	if len(adminsWrap.Password) > 0 {
		adminMap["password"] = utils.Md5String(adminsWrap.Password)
	}

	dbCon := db.GetInstance().GetMysqlDb()
	err := dbCon.Model(&models.Admins{}).Where("id=?", adminsWrap.Id).Updates(adminMap).Error

	if err != nil {
		acl.ErrorResponse(c, "更新失败", nil)
		return
	}
	acl.SuccessResponse(c, "更新成功", nil)
}

func (acl *AdminController) create(c *gin.Context) {
	var adminsWrap models.AdminsWrap
	if err := c.ShouldBindJSON(&adminsWrap); err != nil {
		acl.ErrorResponse(c, err.Error(), nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(adminsWrap.Email, "email").Message("邮箱不能为空")
	valid.Required(adminsWrap.Password, "password").Message("密码不能为空")
	valid.MinSize(adminsWrap.Password, 6, "password").Message("密码最少6位数")
	valid.Required(adminsWrap.ConfirmPassword, "confirmPassword").Message("确认密码不能为空")
	valid.Email(adminsWrap.Email, "email").Message("邮箱非法")
	valid.Required(adminsWrap.Mobile, "mobile").Message("手机号不能为空")
	valid.Mobile(adminsWrap.Mobile, "mobile").Message("手机号非法")
	valid.Range(adminsWrap.Status, 1, 2, "status").Message("状态非法")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			acl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	if adminsWrap.Password != adminsWrap.ConfirmPassword {
		acl.ErrorResponse(c, "两次输入的密码不一致", nil)
		return
	}

	var admins models.Admins
	adminInfo := admins.GetAdminInfoByUsername(adminsWrap.Name)
	if adminInfo.Id > 0 {
		acl.ErrorResponse(c, "该用户名已经存在", nil)
		return
	}

	admins.Name = adminsWrap.Name
	admins.Email = adminsWrap.Email
	admins.Mobile = adminsWrap.Mobile
	admins.Password = utils.Md5String(adminsWrap.Password)
	admins.Status = adminsWrap.Status
	if len(adminsWrap.ImagePath) > 0 {
		admins.ImageUrl = adminsWrap.ImagePath
	}

	dbCon := db.GetInstance().GetMysqlDb()
	err := dbCon.Create(&admins).Error
	if err != nil {
		fmt.Println(err)
		acl.ErrorResponse(c, "添加失败", nil)
		return
	}
	acl.SuccessResponse(c, "添加成功", nil)
}
