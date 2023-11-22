package dao

import (
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/model"
	"github.com/lihaicheng/expires-layui-go/internal/expires-layui-go-gin/store/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ModelThing define
type ModelThing struct {
	tx *gorm.DB
}

// Thing returns ModelThing
func Thing(tx ...*gorm.DB) *ModelThing {
	m := new(ModelThing)
	if len(tx) == 1 {
		m.tx = tx[0]
	} else {
		m.tx = mysql.DB
	}
	return m
}

// Create returns string and error types
func (m *ModelThing) Create(data model.Thing) (uint, error) {
	var err error
	err = m.tx.Create(&data).Error
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

// Delete returns int64 and error types
func (m *ModelThing) Delete(id uint) (int64, error) {
	result := m.tx.Where(&model.Thing{Model: gorm.Model{ID: id}}).Delete(&model.Thing{})
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, errors.Wrapf(errors.New("record not found"), "Thing %d .", id)
	}
	return result.RowsAffected, nil
}

// Update returns int64 and error types
func (m *ModelThing) Update(id uint, data map[string]interface{}) (int64, error) {
	result := m.tx.Model(&model.Thing{}).Where(&model.Thing{Model: gorm.Model{ID: id}}).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// List returns []model.Thing and error types
func (m *ModelThing) List(query map[string]interface{}) ([]*model.Thing, error) {
	var err error
	var rows []*model.Thing
	err = m.tx.Debug().Model(&model.Thing{}).Where(query).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Get returns model.Thing and error types
func (m *ModelThing) Get(query map[string]interface{}) (*model.Thing, error) {
	var err error
	var row *model.Thing
	err = m.tx.Debug().Model(&model.Thing{}).Where(query).Find(&row).Error
	if err != nil {
		return nil, err
	}
	return row, nil
}
