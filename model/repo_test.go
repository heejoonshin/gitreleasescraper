package model

import (
	"fmt"
	"github.com/heejoonshin/gitbot/common"
	"strconv"
	"testing"
)

func init() {
	db := common.TestDBInit()
	db.AutoMigrate(&Repo{})

}

func RepoDumyData() {
	for i := 0; i < 100; i++ {
		owner := "owner " + strconv.Itoa(i)
		repo := "repo " + strconv.Itoa(i)
		inputdata := Repo{
			Owner: owner,
			Repo:  repo,
		}
		inputdata.CreateData()
	}

}

func TestSelectRepo(t *testing.T) {

	//befer
	RepoDumyData()

	result := SelectRepo(1, 10)

	for _, repo := range result {

		fmt.Printf("id : %d owner : %s repo : %s\n", repo.Id, repo.Owner, repo.Repo)

	}

}
