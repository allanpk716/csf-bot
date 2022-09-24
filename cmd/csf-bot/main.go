package main

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/settings"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/telegram_bot"
	"github.com/allanpk716/CSF-Telegram-Bot/internal/pkg/utils"
)

func main() {

	mySettings := settings.NewSettings()
	if utils.IsFile(mySettings.ConfigFPath()) == false {
		err := mySettings.Save()
		if err != nil {
			logger.Panic(err)
		}
	}
	err := mySettings.Read()
	if err != nil {
		logger.Panic(err)
	}
	myBot := telegram_bot.NewTelegramBot(mySettings)
	myBot.Start()
}
