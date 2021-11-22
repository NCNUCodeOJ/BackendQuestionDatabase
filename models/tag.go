package models

import "gorm.io/gorm"

// Tag Database
type Tag struct {
	gorm.Model
	Name string `gorm:"type:varchar(20) NOT NULL;"`
}

//AddTag 創建 tag
func addTag(tag *Tag) (err error) {
	err = DB.Create(&tag).Error
	return
}

// GetTagByName 用 tag name 查 tag
func GetTagByName(name string) (tag Tag, err error) {
	if err = DB.Where(&Tag{Name: name}).First(&tag).Error; err != nil {
		return Tag{}, err
	}
	return
}

func getTagByID(id uint) (tag Tag, err error) {
	err = DB.First(&tag, id).Error
	return
}
