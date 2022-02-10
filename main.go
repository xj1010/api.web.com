package main

import (
	"admin/console"
	"admin/controller"
	"admin/db"
	"admin/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建一个不包含中间件的路由器
	r := gin.Default()
	// 全局中间件
	// 使用 Logger 中间件
	r.Use(gin.Logger())
	// 使用 Recovery 中间件
	r.Use(gin.Recovery())
	db.GetInstance().InitMysqlPool()
	db.GetInstance().GetMysqlDb().LogMode(true)

	//注册各个路由
	controller.GinRouter(r)
	defer db.GetInstance().Close()

	r.StaticFS("/image", http.Dir("./image"))


	//开启计划任务
	//
	go func() {
	 	console.Start()
	}()

	// Listen and serve on 0.0.0.0:8080
	//r.Run(":8080")
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
