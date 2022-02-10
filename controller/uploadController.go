package controller

import (
	"admin/middleware"
	"admin/utils"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
)

type UploadController struct {
	baseController
}

func (uc *UploadController) Routers(r *gin.RouterGroup) {
	g := r.Group("/upload")
	g.Use(middleware.JWT())
	{
		g.POST("/file", uc.upload )
	}
}

func (uc *UploadController) upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		uc.ErrorResponse(c, "获取文件上传信息失败", nil)
		return
	}
	
	uploadPath := "image"
	isExist, _ := utils.PathExists(uploadPath)
	if !isExist {
		if !utils.CreateDir(uploadPath) {
			uc.ErrorResponse(c, "创建目录失败", nil)
			return
		}
	}

	filePath := path.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		uc.ErrorResponse(c, "文件上传错误", nil)
		return
	}

	 f, err := os.Open(filePath)
	 if err != nil {
		 uc.ErrorResponse(c, "文件上传失败", nil)
		 return
	 }
	 defer f.Close()

	 im, _, err := image.DecodeConfig(f)
	 if err != nil {
		 uc.ErrorResponse(c, "获取文件信息失败", err.Error())
		 return
	 }

	uc.SuccessResponse(c, "文件上传成功", map[string]interface{}{
		"imageUrl":  "http://localhost:8089/" + filePath,
		"fileName": file.Filename,
		"size":     file.Size,
		"width" : im.Width,
		"height" : im.Height,
		"filePath" : filePath,
	})
}









