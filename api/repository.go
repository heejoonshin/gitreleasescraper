package api

import (
	"github.com/gin-gonic/gin"
	"github.com/heejoonshin/gitbot/model"
	"github.com/heejoonshin/gitbot/response"
	"github.com/heejoonshin/gitbot/service"
	"github.com/heejoonshin/gitbot/validator"
	"github.com/heejoonshin/gitbot/watcher"
	"strconv"
)

type WatchRepo struct {
	*response.ErrorResponse
	WatchList []model.Repo `json:"data"`
	response.Pagination
}

func GetWatchRepo(c *gin.Context) {

	var page int
	var pagesize int
	var err error

	result := new(WatchRepo)
	defer func() {

		if r := recover(); r != nil {
			e := new(response.ErrorResponse)
			e.Message = r
			c.JSON(200, e)

		} else {
			if result.ErrorResponse == nil {
				c.JSON(200, result)

			} else {
				c.JSON(200, result.ErrorResponse)
			}

		}
	}()
	pagesize = 10
	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			panic("page 파싱 오류")
		}
	}

	if c.Query("pagesize") != "" {
		pagesize, err = strconv.Atoi(c.Query("pagesize"))

		if err != nil {
			panic("pagesize 파싱 오류")
		}
	}
	total, watchList := service.GetWatchRepo(page, pagesize)
	if len(watchList) == 0 {
		result.ErrorResponse = new(response.ErrorResponse)
		result.ErrorResponse.Message = "empty list"
		return
	}
	result.WatchList = watchList
	result.Pagination.Page = page
	result.Pagination.PageSize = pagesize
	result.Pagination.Total = total

}

func SetWatchRepo(c *gin.Context) {
	data := new(model.Repo)
	defer func() {
		if r := recover(); r != nil {
			e := new(response.ErrorResponse)
			e.Message = r
			c.JSON(200, e)

		}

	}()
	err := c.ShouldBindJSON(data)
	if err != nil {
		panic("Post json 데이터 파싱 오류")
	}
	if data.Target == "" {
		data.Target = "github"
	}
	validator.UrlCheck(data)

	err = data.CreateData()
	if err != nil {
		panic(err.Error())
	}
	watcher.Commandchan <- watcher.Repoinfo{Target: data.Target, Owner: data.Owner, Repo: data.Repo}

	c.JSON(200,
		map[string]string{
			"result": "ok",
		})
}

func DelWatchRepo(c *gin.Context) {
	data := new(model.Repo)
	defer func() {
		if r := recover(); r != nil {
			e := new(response.ErrorResponse)
			e.Message = r
			c.JSON(200, e)
		}
	}()
	err := c.ShouldBindJSON(data)
	if err != nil {
		panic("Post json 데이터 파싱 오류")
	}
	if data.Target == "" {
		data.Target = "github"
	}
	watcher.DelChan <- watcher.Repoinfo{Target: data.Target, Owner: data.Owner, Repo: data.Repo}
	err = data.DeleteData()
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200,
		map[string]string{
			"result": "ok",
		})
}
