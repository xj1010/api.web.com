package console

import (
	"admin/console/service"
	"context"
	"unsafe"
)

func Start() {
	initServiceMap()
 }

var baseContext context.Context
var serviceMap  map[string]*service.BaseService

func initServiceMap() {
	baseContext, _ = context.WithCancel(context.Background())
	serviceMap  = make(map[string]*service.BaseService)
	serviceMap["BaiDu"] = (*service.BaseService)(unsafe.Pointer(service.NewBaiDuService(baseContext)))
	serviceMap["WeiBo"] = (*service.BaseService)(unsafe.Pointer(service.NewWeiBoService(baseContext)))
	serviceMap["ZhiHu"] = (*service.BaseService)(unsafe.Pointer(service.NewZhiHuService(baseContext)))
	serviceMap["LaGou"] = (*service.BaseService)(unsafe.Pointer(service.NewLaGouService(baseContext)))
}

func GetServiceMap() map[string]*service.BaseService {
	return serviceMap
}