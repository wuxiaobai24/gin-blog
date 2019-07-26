package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(100);unique;not null"`
	Posts []*Post `gorm:"many2many:post_tag"`
}

func GetTags() ([]*Tag, error) {
	var tags []*Tag
	err := db.Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

func GetTag(id interface{}) (*Tag, error) {
	var tag Tag
	err := db.Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &tag, err
}
