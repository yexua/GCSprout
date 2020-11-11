package db

import (
	"GCSprout/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

const (
	URL = "go:qwe123@tcp(123.56.112.92:3306)/go?charset=utf-8mb4&parseTime=True&loc=Local"
)

var (
	DB *gorm.DB
)

func InitMySQL()  {

	DB, err := gorm.Open("mysql", URL)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println("数据库链接初始化成功")
	gormDBTable(DB)
}

func gormDBTable(db *gorm.DB) {
	err := db.AutoMigrate(model.ToDo{})
	if err != nil{
		log.Fatal("register table failed", err)
	}
	log.Println("register table success")
}