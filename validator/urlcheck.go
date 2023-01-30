package validator

import (
	"github.com/heejoonshin/gitbot/gitinfo"
	"github.com/heejoonshin/gitbot/model"
)

func UrlCheck(repo *model.Repo) {
	url := new(gitinfo.Latest)
	url.Setter(repo.Owner, repo.Repo)

	if url.RequestInfo() != nil {
		panic("url이 존재 하지 않습니다.")
	}

}
