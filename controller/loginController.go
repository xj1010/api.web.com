package controller

import (
	"admin/models"
	"admin/utils"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type LoginController struct {
	baseController
}


func (lcl *LoginController) Routers(r *gin.RouterGroup) {
	g := r.Group("/login")
	{
		g.POST("/index", lcl.index )
		g.POST("/logout", lcl.loginOut )
	}
}

func (lcl *LoginController) index(c *gin.Context) {
	userInfo := lcl.GetUserInfo(c)
	if userInfo.Id > 0 {
		lcl.ErrorResponse(c, "用户已经登入", nil)
		return
	}

	var loginWrap models.LoginWrap
	if err := c.ShouldBindJSON(&loginWrap); err != nil {
		lcl.ErrorResponse(c, "参数错误", nil)
		return
	}

	var username = strings.TrimSpace(loginWrap.Username)
	var password = strings.TrimSpace(loginWrap.Password)

	valid := validation.Validation{}
	valid.Required(username, "username").Message("账号不能为空")
	valid.Required(password, "password").Message("密码不能为空")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			lcl.ErrorResponse(c, err.Message, nil)
			return
		}
	}

	var admin models.Admins
	adminInfo := admin.GetAdminInfoByUsername(username)
	if adminInfo.Id == 0  || adminInfo.Password != utils.Md5String(password) {
		lcl.ErrorResponse(c, "账号或者密码错误", nil)
		return
	}
	if adminInfo.Status == 2 {
		lcl.ErrorResponse(c, "用户账号已经被锁定", nil)
		return
	}

	updateMap := map[string]interface{}{
		"last_login_time" : time.Now(),
		"login_num" : gorm.Expr("login_num + ?", 1),
	}
	admin.UpdateAdminInfoByUid(adminInfo.Id, updateMap)

	//得到token
	var jwt utils.JsonWebToken
	token := jwt.CreateToken(adminInfo.Id, adminInfo.Name, 30000)
	//得到用户权限信息
	var adminNode models.AdminNode
	nodeList := adminNode.GetUserNodeAuthList(adminInfo.Id)

	var userAuthInfo models.UserAuthInfo
	userAuthInfo.AdminInfo = adminInfo
	userAuthInfo.NodeInfo = nodeList
	models.UserAuthSession[token] = userAuthInfo

	lcl.SuccessResponse(c, "用户登入成功", map[string]interface{}{
		"token" : token,
	})
}

func (lcl *LoginController) loginOut(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")
	if token == "" || len(token) == 0 {
		token = c.Query("token")
	}
	if token == "" {
		lcl.ErrorResponse(c, "token非法", nil)
		return
	}

	var jwt utils.JsonWebToken
	_, isVerify := jwt.VerifyToken(token)
	if !isVerify {
		lcl.ErrorResponse(c, "token非法", nil)
		return
	}

	if _, ok := models.UserAuthSession[token]; !ok {
		lcl.ErrorResponse(c, "用户未登入", nil)
		return
	}

	delete(models.UserAuthSession, token)

	lcl.SuccessResponse(c, "用户退出成功", nil)
	return
}


