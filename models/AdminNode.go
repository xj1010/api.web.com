package models

import (
	"admin/db"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type  AdminNode struct {
	Id int `gorm:"primary_key"  json:"id"   form:"name"`
	Name string `gorm:"column:name;type:varchar(150)" json:"name"   form:"name"`
	Pid int  `gorm:"column:pid;type:int(10);"  json:"pid"    form:"pid" `
	Path string  `gorm:"column:path;type:varchar(20);" json:"path"   form:"path"`
	Icon string  `gorm:"column:icon;type:varchar(20);" json:"icon"   form:"icon"`
	PermissionName string  `gorm:"column:permission_name;type:varchar(20);" json:"permission_name"   form:"permission_name"`
	GuardName string  `gorm:"column:guard_name;type:varchar(20);" json:"guard_name"   form:"guard_name"`
	Sort int  `gorm:"column:sort;type:int(10);"  json:"sort"    form:"sort" `
	Level int  `gorm:"column:level;type:int(10);"  json:"level"    form:"level" `
	Status int  `gorm:"column:status;type:int(10);"  json:"status"    form:"status" `
	AdminName string  `gorm:"column:admin_name;type:varchar(200);"  json:"admin_name"    form:"admin_name" `
	Base
}

type TreeNode struct {
	Id int `json:"id"`
	Pid int  `json:"pId"`
	Name string  `json:"name"`
	Open bool  `json:"open"`
}

type UserAuthInfo struct {
	AdminInfo  Admins
	NodeInfo []AdminNode
}

var UserAuthSession map[string]UserAuthInfo
func init()  {
	UserAuthSession = make(map[string]UserAuthInfo)
}

func (adminNode *AdminNode) GetFormatUserNodeAuthList(uid int) []AuthNode {
	var authNodeList []AuthNode

	nodeInfoList := adminNode.GetUserNodeAuthList(uid)
	if len(nodeInfoList) == 0 {
		return authNodeList
	}

	for _, node := range nodeInfoList {
		if node.Level != 1 {
			continue
		}

		var authNode AuthNode
		var meta Meta

		authNode.Name =  node.Name
		authNode.Path = node.Path
		authNode.Component = "layout"
		authNode.AlwaysShow = true
		meta.Icon = node.Icon
		meta.Title = node.Name
		authNode.Meta = meta

		for _, cNode := range nodeInfoList {
			if cNode.Pid != node.Id {
				continue
			}

			var secondNode Children
			var meta Meta
			secondNode.Name = cNode.Name
			secondNode.Component = cNode.PermissionName
			secondNode.Path = cNode.Path
			meta.Title = cNode.Name
			meta.Icon = cNode.Icon
			secondNode.Meta = meta

			authNode.Children = append(authNode.Children, secondNode)
		}

		authNodeList = append(authNodeList, authNode)
	}

	return authNodeList
}

func (adminNode *AdminNode) GetUserNodeAuthList(uid int) []AdminNode {
	var nodeInfoList []AdminNode
	var sql string
	dbCon := db.GetInstance().GetMysqlDb()

	if uid == 1 {
		sql = "select id,pid,name,level,path,icon,permission_name from  admin_node where status"
		dbCon.Raw(sql).Scan(&nodeInfoList)
	} else {
		sql = "select admin_node.id,pid, admin_node.name,level,path, icon,permission_name" +
			" from admin_role_user" +
			" left join  admin_role_node on admin_role_user.role_id=admin_role_node.role_id" +
			" left join admin_node on admin_node.id=admin_role_node.node_id" +
			" where admin_role_user.status=1 and admin_role_node.status=1 and admin_node.status=1 and user_id =? group by admin_node.id"
		dbCon.Raw(sql, uid).Scan(&nodeInfoList)
	}

	return nodeInfoList
}

func (adminNode *AdminNode) GetList(searchMap map[string]string, page, limit int) (count int, nodeInfo []AdminNode) {
	dbCon := db.GetInstance().GetMysqlDb()

	orderBy := " id desc"
	sort := strings.TrimSpace(searchMap["sort"])
	if len(sort) > 0 {
		if sort == "-id" {
			orderBy = "id asc"
		}
	}

	pid, _ := strconv.Atoi(searchMap["pid"])
	if pid > 0 {
		dbCon = dbCon.Where("pid=?", pid).Or("id=?", pid)
		orderBy = "id asc"
	}

	name := strings.TrimSpace(searchMap["name"])
	if len(name) > 0 {
		dbCon = dbCon.Where("name=?",name)
	}


	if page <= 0 {
		page = 1
	}

	//总数
	dbCon.Model(AdminNode{}).Count(&count)
	fPage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(count) / float64(limit)), 64)
	maxPage := int(math.Ceil(fPage))
	if page > maxPage {
		page = maxPage
	}

	//数据
	dbCon = dbCon.Order(orderBy).Limit(limit).Offset((page - 1) * limit)

	if err := dbCon.Find(&nodeInfo).Error; err != nil {
		return 0, nodeInfo
	}
	return count, nodeInfo
}
