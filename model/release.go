package model

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/heejoonshin/gitbot/common"
	"github.com/jinzhu/gorm"
	"time"
)

type Latest struct {
	Tag         string `gorm:"type:varchar(1020)"`
	CreatedAt   time.Time
	PublishedAt time.Time
	Repo        string `gorm:"type:varchar(1020)"`
	Url         string `gorm:"type:varchar(2020)"`
	Id          string `gorm:"type:varchar(2020);primary_key"`
}

func (L *Latest) SeleteAll() []Latest {
	db := common.DB
	var ret []Latest
	db.Find(&ret)
	return ret

}
func (L *Latest) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("created_at", time.Now())
	return nil
}
func (L *Latest) AfterCreate(scope *gorm.Scope) error {
	fmt.Println(L.Tag)
	for key, _ := range common.ChatID {
		id := key
		remsg := tgbotapi.NewMessage(id, L.Url)
		common.Bot.Send(remsg)

	}
	return nil
}
func (L *Latest) AfterUpdate(scope *gorm.Scope) error {
	fmt.Println(L.Tag)
	for key, _ := range common.ChatID {
		id := key
		remsg := tgbotapi.NewMessage(id, L.Url)
		common.Bot.Send(remsg)

	}
	return nil
}

func (L *Latest) CreateLatest() error {
	db := common.GetDB()
	if err := db.Create(&L); err != nil {
		return nil

	}
	return nil
}
func (L *Latest) UpdateLatest() error {
	db := common.GetDB()
	if err := db.Model(&Latest{}).Update(&L); err != nil {
		return nil
	}
	return nil
}
