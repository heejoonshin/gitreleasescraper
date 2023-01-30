package watcher

import (
	"fmt"
	"log"

	"github.com/heejoonshin/gitbot/gitinfo"
	"github.com/heejoonshin/gitbot/model"

	"time"
)

var Commandchan = make(chan Repoinfo)
var DelChan = make(chan Repoinfo)

type Repoinfo struct {
	Target string
	Owner  string
	Repo   string
}

func (r *Repoinfo) String() string {
	return r.Owner + ":" + r.Repo
}

type WatchRepo struct {
	watchList   map[string]interface{}
	commandchan chan Repoinfo
	releases    map[string]interface{}
}

func New(_commandchan chan Repoinfo) *WatchRepo {

	h := map[string]interface{}{}
	r := model.Repo{}
	All := r.SelectAll()

	for _, val := range All {
		temp := gitinfo.Latest{}
		temp.Setter(val.Owner, val.Repo)
		h[temp.String()] = temp
	}

	ret := &WatchRepo{commandchan: _commandchan, watchList: h, releases: map[string]interface{}{}}
	releases := model.Latest{}
	releasedata := releases.SeleteAll()

	for _, val := range releasedata {
		ret.releases[val.Repo] = val
		fmt.Println(val)
	}
	return ret

}
func (w *WatchRepo) Run() {
	timer := time.NewTicker(30 * time.Minute)
	defer timer.Stop()
	for {

		select {
		case repo := <-w.commandchan:

			w.add(repo)
		case <-timer.C:

			//timer = time.NewTimer(1 * time.Second)
			//fmt.Println("run")
			//fmt.Print(w.watchList.Len())
			timer.Reset(30 * time.Minute)

			for key, repo := range w.watchList {

				x, ok := repo.(gitinfo.Latest)
				if !ok {
					continue
				}
				x.RequestInfo()

				r := key
				fmt.Println(r)
				newdata := model.Latest{Tag: x.Tagname, PublishedAt: x.PublishedAt, Repo: r, Url: x.Url, Id: x.String()}
				fmt.Println("ttt:", newdata)

				if release, ok := w.releases[r]; !ok {

					newdata.CreateLatest()
					w.releases[r] = newdata

				} else {
					data := release.(model.Latest)
					if data.Tag != newdata.Tag && newdata.Tag != "" {
						newdata.UpdateLatest()
						w.releases[r] = newdata

					}
				}

			}
		case repo := <-DelChan:
			w.del(repo)

		}
	}

}

func (w *WatchRepo) add(info Repoinfo) {

	if _, ok := w.watchList[info.String()]; !ok {

		repo := &model.Repo{Target: info.Target, Owner: info.Owner, Repo: info.Repo}
		temp := gitinfo.Latest{}
		temp.Setter(info.Owner, info.Repo)

		w.watchList[temp.String()] = temp
		temp.RequestInfo()
		newdata := model.Latest{Tag: temp.Tagname, PublishedAt: temp.PublishedAt, Repo: info.Repo, Url: temp.Url, Id: temp.String()}
		newdata.CreateLatest()
		w.releases[newdata.Id] = newdata
		repo.CreateData()

	} else {
		log.Print("already watched")
	}
}

func (w *WatchRepo) del(info Repoinfo) {

	if _, ok := w.watchList[info.String()]; ok {
		delete(w.watchList, info.String())
	}

}
