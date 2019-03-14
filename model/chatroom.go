package model

import "github.com/heejoonshin/gitbot/common"

type Chatroom struct{
	ChatId int64 `gorm:"primary_key"`
}

func (c *Chatroom)SelectAll() []Chatroom {
	db := common.DB
	ret := []Chatroom {}
	db.Find(&ret)
	return ret
}

func (c *Chatroom)CreateData() error{
	db := common.GetDB()

	db.Create(&c)
	return nil

}