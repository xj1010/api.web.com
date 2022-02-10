package models

import (
	"admin/db"
)

type  AdminRoleUser struct {
	Id int `gorm:"primary_key"`
	RoleId int `gorm:"column:role_id;type:int(10)" json:"role_id, string"   form:"role_id"`
	UserId int  `gorm:"column:user_id;type:int(10);"  json:"user_id,string"    form:"user_id" `
	Status int `gorm:"column:status;type:int(10)" json:"status, string"   form:"status"`
	Admins Admins `gorm:"ForeignKey:UserId;AssociationForeignKey:Id" json:"admin,omitempty" form:"admin,omitempty"`
	Role AdminRole `gorm:"ForeignKey:RoleId;AssociationForeignKey:Id" json:"role,omitempty" form:"role,omitempty"`
	Base
}

func (adminRoleUser *AdminRoleUser) GetUserInfoListByRoleId(roleId int)  []AdminRoleUser {
	dbCon := db.GetInstance().GetMysqlDb()
	var roleUser  []AdminRoleUser
	if roleId > 0 {
		dbCon.Preload("Admins").Where("role_id=? and status=1",  roleId ).Find(&roleUser)
	}

	return roleUser
}


