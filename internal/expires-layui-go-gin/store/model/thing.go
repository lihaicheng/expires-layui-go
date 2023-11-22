package model

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const (
	thingTable = "thing"
)

type Thing struct {
	gorm.Model
	Name           string    `gorm:"column:Name; type:string; size:36; not null; unique; <-:create;"`
	Type           string    `gorm:"column:type; type:string; size:10; not null; "`
	Count          int64     `gorm:"column:count; type:int; not null;"`
	ProductionDate time.Time `gorm:"column:production_date; not null;"`
	Life           int64     `gorm:"column:life; type:int; not null;"`
	ExpirationDate time.Time `gorm:"column:expiration_date; not null;"`
	Note           string    `gorm:"column:note; type:longtext; "`
}

// TableName returns table name
func (t Thing) TableName() string {
	return tableName(thingTable)
}

// BeforeSave is table operation hooks
func (t *Thing) BeforeSave(tx *gorm.DB) error {
	zap.L().Info("BeforeSave")
	return nil
}

// BeforeCreate is table operation hooks
func (t *Thing) BeforeCreate(tx *gorm.DB) error {
	zap.L().Info("BeforeCreate")
	return nil
}

// BeforeUpdate is table operation hooks
func (t *Thing) BeforeUpdate(tx *gorm.DB) error {
	zap.L().Info("BeforeUpdate")
	return nil
}

// BeforeDelete is table operation hooks
func (t *Thing) BeforeDelete(tx *gorm.DB) error {
	zap.L().Info("BeforeDelete")
	return nil
}
