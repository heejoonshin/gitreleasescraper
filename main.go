package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/heejoonshin/gitbot/common"
	"github.com/heejoonshin/gitbot/model"
	"github.com/heejoonshin/gitbot/watcher"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"strings"
)
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Latest{})
	db.AutoMigrate(&model.Repo{})
	db.AutoMigrate(&model.Chatroom{})
}


func main() {



	arg := os.Args[1]
	fmt.Print(arg)
	db := common.Init()
	Migrate(db)
	chatroom := model.Chatroom{}
	chatrooms := chatroom.SelectAll()
	for _,val := range chatrooms{
		common.ChatID.Set(val.ChatId,val)
	}


	defer db.Close()
	commandchan := make(chan watcher.Repoinfo)
	w := watcher.New(commandchan)
	go w.Run()
	bot, err := tgbotapi.NewBotAPI(arg)
	common.Bot = bot
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout=60


	updates, err := bot.GetUpdatesChan(u)


	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		parse := strings.Split(msg.Text," ")
		if parse[0] == "watchrepo"{
			fmt.Println(parse)
			commandchan <- watcher.Repoinfo{Owner:parse[1],Repo:parse[2]}
		}else if(parse[0] == "help"){
			//t := tgbotapi.N{msg.BaseChat}
			x, _ :=bot.GetMe()
			fmt.Print(msg)
			remsg := tgbotapi.NewMessage(msg.ChatID,"watchrepo git오너명 git저장소명")
			//remsg.ReplyToMessageID = 787098406
			fmt.Print(x)
			common.Bot.Send(remsg)
		}
		if _,ok := common.ChatID.Get(msg.ChatID); !ok{
			newroom := &model.Chatroom{ChatId:msg.ChatID}
			fmt.Print(newroom)
			common.ChatID.Set(msg.ChatID,newroom)
			newroom.CreateData()

		}


		//bot.Send(msg)
	}



}
