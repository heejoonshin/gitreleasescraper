package service

import "github.com/heejoonshin/gitbot/model"

func GetWatchRepo(page, pageSize int) (total int, result []model.Repo) {
	count := model.CountAllWatchRepo()
	data := model.SelectRepo(page, pageSize)

	return count, data

}
