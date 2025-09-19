package database

import (
	"github.com/CPtung/mattercontroller/pkg/model"
	"gorm.io/gorm"
)

func Load(id string) (*model.MatterDevice, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	rt := model.MatterDevice{}
	if err := dbClient().Where("id = ?", id).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, err
}

func Store(device *model.MatterDevice) (*model.MatterDevice, error) {
	mutex.Lock()
	defer mutex.Unlock()

	f := model.MatterDevice{}
	result := dbClient().Model(f).Where("id = ?", device.ID).First(&f)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		err := dbClient().Create(device).Error
		return device, err
	}
	return device, dbClient().Model(f).Updates(device).Error
}

func Delete(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	f := model.MatterDevice{}
	return dbClient().Delete(&f, id).Error
}
