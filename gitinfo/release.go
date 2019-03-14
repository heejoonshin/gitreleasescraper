package gitinfo

import (
	"github.com/heejoonshin/gitbot/common"
	"time"
)

type Latest struct{
	Tagname string `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Url string `json:"html_url"`
	owner string
	repo string


}
func (l *Latest)Setter(owner ,repo string){
	l.owner = owner
	l.repo = repo

}
func (l *Latest) RequestInfo() error{

	u := "https://api.github.com/repos/"+l.owner + "/" + l.repo + "/releases/latest"

	common.GetJSON(u,&l)


	//fmt.Println(l)
	return nil

}
