package model

import (
	"errors"
	"github.com/heejoonshin/gitbot/common"
)

type Repo struct {
	Id     uint64 `grom:"AUTO_INCREMENT;primary_key" json:"id,omitempty"`
	Owner  string `grom:"type:varchar(1020)" json:"owner" binding:"required"`
	Repo   string `gorm:"type:varchar(1020)" json:"repo" binding:"required"`
	Target string `gorm:"type:varchar(1020)" json:"target"`
}

func (r *Repo) SelectAll() []Repo {
	db := common.DB
	ret := []Repo{}
	db.Find(&ret)

	return ret
}

func SelectRepo(page, pageSize int) []Repo {
	db := common.DB
	ret := []Repo{}
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&ret)
	return ret
}

func CountAllWatchRepo() int {

	db := common.DB
	var count int
	db.Model(&Repo{}).Count(&count)

	return count
}

func (r *Repo) CreateData() error {
	db := common.GetDB()
	checkrepo := new(Repo)
	db = db.Where("owner = ? AND repo = ?", r.Owner, r.Repo).Find(&checkrepo)
	if checkrepo.Id != 0 {
		return errors.New("중복된 데이터")
	}
	db.Create(&r)
	return nil
}
func (r *Repo) DeleteData() error {
	db := common.GetDB()
	checkrepo := new(Repo)
	db = db.Where("owner = ? AND repo = ?", r.Owner, r.Repo).Find(&checkrepo)
	if checkrepo.Id == 0 {
		return errors.New("존재하지 않는 데이터")
	}
	db.Where("owner = ? AND repo = ?", r.Owner, r.Repo).Delete(&Repo{})

	return nil
}
