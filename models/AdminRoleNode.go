package models

import (
	"admin/db"
)

type  AdminRoleNode struct {
	Id int `gorm:"primary_key"`
	RoleId int `gorm:"column:role_id;type:int(10)" json:"role_id, string"   form:"role_id"`
	NodeId int  `gorm:"column:node_id;type:int(10);"  json:"node_id,string"    form:"node_id" `
	Status int `gorm:"column:status;type:int(10)" json:"status, string"   form:"status"`
	AdminNode AdminNode `gorm:"ForeignKey:NodeId;AssociationForeignKey:Id" json:"node,omitempty" form:"node,omitempty"`
	Base
}

func (adminRoleNode *AdminRoleNode) GetNodeListByRoleId(roleId int)  []AdminRoleNode{
	dbCon := db.GetInstance().GetMysqlDb()
	var roleNode  []AdminRoleNode
	if roleId > 0 {
		dbCon.Preload("AdminNode").Where("status=1 and role_id=?",  roleId ).Find(&roleNode)
	}

	return roleNode
}


