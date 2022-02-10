package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Base struct {
	CreatedAt MyTime `gorm:"column(created_at);" json:"created_at" form:"-"`
	UpdatedAt MyTime `gorm:"column(updated_at);" json:"updated_at" form:"-"`
}

func (v Base) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("create_time", time.Now())
	scope.SetColumn("update_time", time.Now())
	return nil
}

func (v Base) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("update_time", time.Now())
	return nil
}



