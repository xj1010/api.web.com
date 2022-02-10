package console

import (
	"admin/console/service"
	"sort"
	"unsafe"
)

var StructMap map[string]*service.BaseService

func Start() {
	initStructMap()
 }

func initStructMap() {
	StructMap = make(map[string]*service.BaseService)
	StructMap["WeiBoService"] = (*service.BaseService)(unsafe.Pointer(new(service.WeiBoService).Execute()))
	StructMap["BaiDuService"] = (*service.BaseService)(unsafe.Pointer(new(service.BaiDuService).Execute()))
	StructMap["ZhiHuService"] = (*service.BaseService)(unsafe.Pointer(new(service.ZhiHuService).Execute()))
	StructMap["LaGouService"] = (*service.BaseService)(unsafe.Pointer(new(service.LaGouService).Execute()))
}

func List() map[string]*service.BaseService {
	sortMap := make(map[string]*service.BaseService)
	var strs []string
	for key, _ := range StructMap {
		strs = append(strs, key)
	}

	sort.Strings(strs)
	for _, k := range strs {
		sortMap[k] = StructMap[k]
	}

	return StructMap
}