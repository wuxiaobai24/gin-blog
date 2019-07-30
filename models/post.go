package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(100);unique;not null"`
	Content string `gorm:'type:text;`
	Tags    []*Tag `gorm:"many2many:post_tag"`
}

func GetPosts(pageNum int, pageSize int) ([]*Post, error) {
	var posts []*Post
	err := db.Offset(pageNum).Limit(pageSize).Find(&posts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return posts, nil
}

func GetPost(id string) (*Post, error) {
	var post Post
	err := db.Where("id = ?", id).First(&post).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &post, err
}

func PostCount() (int, error) {
	var count int
	err := db.Model(&Post{}).Count(&count).Error
	return count, err
}

func AddPost(data map[string]interface{}) error {
	content := strings.Replace(data["Content"].(string), "\n\r", "\n", -1)
	post := Post{
		Title:   data["Title"].(string),
		Content: content,
		Tags:    data["Tags"].([]*Tag),
	}
	if err := db.Set("gorm:association_autocreate", true).Create(&post).Error; err != nil {
		return err
	}
	return nil
}

func (p *Post) UpdatePost() error {
	return db.Save(p).Error
}

func (p *Post) DeletePost() error {
	return db.Delete(p).Error
}
