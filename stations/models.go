package stations

import (
	_ "fmt"

	"12306.com/12306/common"
)

// Station Station model
type Station struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Spell string `json:"spell"`
	Name  string `json:"name"`
}

//FindStations Find all staions from database
func FindStations() ([]Station, error) {
	db := common.GetDB()
	var models []Station
	err := db.Order("spell").Find(&models).Error
	return models, err
}

//FindHotStations Find all hot staions from database
func FindHotStations() ([]Station, error) {
	db := common.GetDB()
	var models []Station
	err := db.Order("spell").Where("is_hot = ?", 1).Find(&models).Error
	return models, err
}
