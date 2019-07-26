package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/wuxiaobai24/gin-blog/config"
)

var db *gorm.DB
var err error

func InitDB() error {
	conf := config.Conf.Database
	var args string
	if conf.Type == "sqlite3" {
		args = conf.DatabaseName
	} else if conf.Type == "mysql" {
		args = fmt.Sprintf("%s:%s@/%s?charset=utf8?parseTime=True&loc=Local", conf.User, conf.Password, conf.DatabaseName)
	}

	db, err = gorm.Open(conf.Type, args)
	if err != nil {
		return err
	}
	db.AutoMigrate(&Post{}, &Tag{})
	db = db.Set("gorm:auto_preload", true)
	return nil
}

func CloseDB() {
	defer db.Close()
}
