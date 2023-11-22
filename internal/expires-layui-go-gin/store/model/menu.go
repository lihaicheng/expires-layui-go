package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	SplitFlag4Child = "|"
	menuTable       = "menu"
)

type Menu struct {
	gorm.Model
	Code    string `gorm:"column:code; type:string; size:36; not null; unique; <-:create;"`
	Type    string `gorm:"column:type; type:string; size:10; "` // home、logo、head、directory、menu
	Title   string `gorm:"column:title; type:string; size:50; "`
	Image   string `gorm:"column:image; type:string; size:50; "`
	Icon    string `gorm:"column:icon; type:string; size:50; "`
	Href    string `gorm:"column:href; type:string; size:50; "`
	Target  string `gorm:"column:target; type:string; size:50; "`
	IsShown int64  `mapstructure:"is_shown" gorm:"column:is_shown; type:tinyint; not null;"` // 0表示不展示， 1表示展示
	Child   string `gorm:"column:child; type:longtext; "`
}

// TableName returns table name
func (t Menu) TableName() string {
	return tableName(menuTable)
}

// BeforeSave is table operation hooks
func (t *Menu) BeforeSave(tx *gorm.DB) error {
	zap.L().Info("BeforeSave")
	return nil
}

// BeforeCreate is table operation hooks
func (t *Menu) BeforeCreate(tx *gorm.DB) error {
	zap.L().Info("BeforeCreate")
	return nil
}

// BeforeUpdate is table operation hooks
func (t *Menu) BeforeUpdate(tx *gorm.DB) error {
	zap.L().Info("BeforeUpdate")
	return nil
}

// BeforeDelete is table operation hooks
func (t *Menu) BeforeDelete(tx *gorm.DB) error {
	zap.L().Info("BeforeDelete")
	return nil
}
