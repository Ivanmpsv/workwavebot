package main

import (
	"log"
	"workwavebot/startbot"
)

func main() {
	// Создаём экземпляр бота
	bot, err := startbot.Createbot()
	if err != nil {
		log.Panic(err)
	}

	// Передаём бота в StartBot
	startbot.StartBot(bot)
}
