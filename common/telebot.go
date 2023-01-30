package common

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI
var ChatID = map[int64]interface{}{}
