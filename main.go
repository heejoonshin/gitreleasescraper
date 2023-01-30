package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/heejoonshin/gitbot/api"
	"github.com/heejoonshin/gitbot/common"
	"github.com/heejoonshin/gitbot/model"
	"github.com/heejoonshin/gitbot/watcher"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"strings"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Latest{})
	db.AutoMigrate(&model.Repo{})
	db.AutoMigrate(&model.Chatroom{})
}
func telebotStarter(arg string) {

	chatroom := model.Chatroom{}
	chatrooms := chatroom.SelectAll()
	for _, val := range chatrooms {
		//common.ChatID.Set(val.ChatId, val)

		common.ChatID[val.ChatId] = &val
	}

	w := watcher.New(watcher.Commandchan)
	go w.Run()
	bot, err := tgbotapi.NewBotAPI(arg)
	common.Bot = bot
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		parse := strings.Split(msg.Text, " ")
		if parse[0] == "watchrepo" {
			fmt.Println(parse)
			if len(parse) == 3 {
				watcher.Commandchan <- watcher.Repoinfo{Target: "github", Owner: parse[1], Repo: parse[2]}
			} else if len(parse) == 4 {
				watcher.Commandchan <- watcher.Repoinfo{Target: parse[1], Owner: parse[2], Repo: parse[3]}
			}

		} else if parse[0] == "help" {
			//t := tgbotapi.N{msg.BaseChat}
			x, _ := bot.GetMe()
			fmt.Print(msg)
			remsg := tgbotapi.NewMessage(msg.ChatID, "watchrepo git오너명 git저장소명")
			//remsg.ReplyToMessageID = 787098406
			fmt.Print(x)
			common.Bot.Send(remsg)
		}
		if _, ok := common.ChatID[msg.ChatID]; !ok {
			newroom := &model.Chatroom{ChatId: msg.ChatID}
			fmt.Print(newroom)
			common.ChatID[msg.ChatID] = newroom
			newroom.CreateData()

		}

		//bot.Send(msg)
	}

}

type Users struct {
	Name string

	Phone string

	Age int

	Id int
}

func main() {

	arg := os.Args[1]
	fmt.Print(arg)
	db := common.Init()
	Migrate(db)
	defer db.Close()
	go telebotStarter(arg)
	r := gin.Default()

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")

	r.GET("/index", func(c *gin.Context) {

		repoData := model.Repo{}
		resData := repoData.SelectAll()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "Main website",
			"repodata": resData,
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/test", api.GetWatchRepo)
	r.POST("test", api.SetWatchRepo)
	r.DELETE("/test", api.DelWatchRepo)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
