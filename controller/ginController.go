package controller

import "github.com/gin-gonic/gin"

func GinRouter(r *gin.Engine)  {
	rr := r.Group("")
	new(AdminController).Routers(rr)
	new(NodeController).Routers(rr)
	new(LoginController).Routers(rr)
	new(RoleController).Routers(rr)
	new(RoleNodeController).Routers(rr)
	new(RoleUserController).Routers(rr)
	new(UploadController).Routers(rr)
	new(HotController).Routers(rr)
	new(ServiceController).Routers(rr)
	new(LagouController).Routers(rr)

}