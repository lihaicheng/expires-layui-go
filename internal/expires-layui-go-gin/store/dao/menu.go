package dao

import (
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/model"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/mysql"
	"gorm.io/gorm"
)

// ModelMenu define
type ModelMenu struct {
	tx *gorm.DB
}

// Menu returns ModelMenu
func Menu(tx ...*gorm.DB) *ModelMenu {
	m := new(ModelMenu)
	if len(tx) == 1 {
		m.tx = tx[0]
	} else {
		m.tx = mysql.DB
	}
	return m
}

// List returns []model.Menu and error types
func (m *ModelMenu) List(query map[string]interface{}) ([]*model.Menu, error) {
	var err error
	var rows []*model.Menu
	err = m.tx.Debug().Model(&model.Menu{}).Where(query).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Get returns model.Menu and error types
func (m *ModelMenu) Get(query map[string]interface{}) (*model.Menu, error) {
	var err error
	var row *model.Menu
	err = m.tx.Debug().Model(&model.Menu{}).Where(query).Find(&row).Error
	if err != nil {
		return nil, err
	}
	return row, nil
}
