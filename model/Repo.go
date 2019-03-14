package model

import (
	"github.com/heejoonshin/gitbot/common"
)

type Repo struct{
	Id uint64 `grom:"AUTO_INCREMENT;primary_key"`
	Owner string `grom:"type:varchar(1020)"`
	Repo string `gorm:"type:varchar(1020)"`

}

func (r *Repo)SelectAll() []Repo {
	db := common.DB
	ret := []Repo {}
	db.Find(&ret)

	return ret
}

func (r *Repo)CreateData() error{
	db := common.GetDB()

	db.Create(&r)

	return nil
}


