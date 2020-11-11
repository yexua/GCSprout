package model

type ToDo struct {
	// 封装了一些基本字段
	//gorm.Model
	Id int `json:"id"`
	Title string `json:"title"`
	Status string `json:"status"`
}


