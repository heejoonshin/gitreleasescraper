package watcher

import (
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/heejoonshin/gitbot/gitinfo"
	"github.com/heejoonshin/gitbot/model"
	"log"
	"time"
)
type Repoinfo struct{
	Owner string
	Repo string
}

type WatchRepo struct{

	watchList hashmap.HashMap
	commandchan chan Repoinfo
	releases hashmap.HashMap

}
func New(_commandchan chan Repoinfo) *WatchRepo{

	h := hashmap.HashMap{}
	r := model.Repo{}
	All :=r.SelectAll()

	for _,val := range All{
		temp := gitinfo.Latest{}
		temp.Setter(val.Owner,val.Repo)
		h.Set(val.Repo,temp)
	}


	ret := &WatchRepo{commandchan:_commandchan,watchList: h}
	releases := model.Latest{}
	releasedata := releases.SeleteAll()
	for _ ,val := range releasedata{
		ret.releases.Set(val.Repo,val)
	}
	return ret

}
func (w *WatchRepo)Run(){
	timer := time.NewTimer(1* time.Second)
	defer timer.Stop()
	for {

		select {
		case repo := <-w.commandchan:

			w.add(repo)
		case <-timer.C:
			//fmt.Print(w.watchList.Len())

			for repo:= range w.watchList.Iter(){

				x ,ok := repo.Value.(gitinfo.Latest)
				if !ok{
					continue
				}
				x.RequestInfo()

				r := repo.Key.(string)
				fmt.Println(r)
				newdata := model.Latest{Tag:x.Tagname,PublishedAt:x.PublishedAt,Repo:r,Url:x.Url}
				fmt.Println("ttt:",newdata)

				if release,ok := w.releases.Get(r); !ok{

					newdata.CreateLatest()
					w.releases.Set(r,newdata)



				}else{
					data := release.(model.Latest)
					if data.Tag != newdata.Tag && newdata.Tag != "" {
						newdata.UpdateLatest()
						w.releases.Set(r,newdata)

					}
				}

			}
			timer = time.NewTimer(30*time.Minute)

		}
	}

}



func (w *WatchRepo) add(info Repoinfo) {

	if _,ok :=w.watchList.Get(info.Repo); !ok{

		repo := &model.Repo{Owner:info.Owner,Repo:info.Repo}
		temp := gitinfo.Latest{}
		temp.Setter(info.Owner,info.Repo)


		w.watchList.Set(info.Repo,&temp)
		temp.RequestInfo()
		newdata := model.Latest{Tag:temp.Tagname,PublishedAt:temp.PublishedAt,Repo:info.Repo,Url:temp.Url}
		newdata.CreateLatest()
		w.releases.Set(info.Repo,newdata)
		repo.CreateData()

	}else{
		log.Print("already watched")
	}
}

