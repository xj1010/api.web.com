package middleware

import (
	"admin/models"
	"admin/utils"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		constantRoutes := map[string]bool{
			"/admin/info": true,
			"/admin/node": true,
		}
		//校验token
		 msg, code, userAuthInfo := verifyToken(c)
		 if code == 0 {
		 	//验证节点权限
		 	var flag bool
		 	path := c.Request.URL.Path
		 	if _, ok := constantRoutes[path]; !ok {
				for _, node := range userAuthInfo.NodeInfo {
					if path == node.Path {
						flag = true
						break
					}
				}
				if flag == false {
					code, msg  = -5, "您没有权限访问该节点"
				}
			}
		 }

		if code != 0 {
			c.JSON(200, gin.H{
				"code": code,
				"message": msg,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

func verifyToken(c  *gin.Context) (msg string, code int, userAuthInfo models.UserAuthInfo) {
	token := c.Request.Header.Get("X-Token")
	if token == "" {
		token = c.Query("token")
	}
	if token == "" {
		cookie, err := c.Request.Cookie("Admin-Token")
		if err != nil {
			token = ""
		} else {
			token = cookie.Value
		}
	}

	if token == "" {
		return "token为空，非法请求", -1, userAuthInfo
	}

	var jwt utils.JsonWebToken
	customClaims, isVerify := jwt.VerifyToken(token)
	if !isVerify {
		return "token鉴权失败", -2, userAuthInfo
	}

	if customClaims.UserId <= 0   || len(customClaims.Username) == 0 {
		return "用户非法", -3, userAuthInfo
	}

	if _, ok := models.UserAuthSession[token]; !ok {
		return "用户未登入", -4, userAuthInfo
	}

	return msg ,code, models.UserAuthSession[token]
}