package startbot

import (
	"fmt"
	"workwavebot/bonushr"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userStates = make(map[int64]string) // Карта для отслеживания состояний пользователей

// StartBot запускает бота
func StartBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0) // бот будет начинать с последнего непрочитанного обновления.
	u.Timeout = 60             // ждём 60sec чтобы получить новые обновления от Telegram API, снижает нагрузку на сервер tg

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}

// sendMessage отправляет сообщение пользователю
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

// handleMessage обрабатывает входящие сообщения
func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userState := userStates[chatID]

	switch {
	case message.Text == "/start" || message.Text == "старт бот" || message.Text == "бот проснись":
		sendMessage(bot, chatID, "Локтар огар! Напишите имя клиента чтобы рассчитать вашу премию")
		userStates[chatID] = "waiting_for_client"

	case userState == "waiting_for_client":
		handleClient(bot, chatID, message.Text)

	case userState == "waiting_for_salary_alfa" || userState == "waiting_for_salary_x5":
		handleSalary(bot, chatID, message.Text, userState)
	}

	switch {
	case message.Text == "/new":
		sendMessage(bot, chatID, "Добавить, удалить или обновить (формулу) клиента: нажмите 1 2 3 соответственно")
		userStates[chatID] = "1 2 3"
	case userState == "1 2 3" && message.Text == "1":
		sendMessage(bot, chatID, "Добавить клиента в разработке")

	case userState == "1 2 3" && message.Text == "2":
		sendMessage(bot, chatID, "Обновить формулу клиента в разработке")

	case userState == "1 2 3" && message.Text == "3":
		sendMessage(bot, chatID, "Удалить клиента в разработке")
	}

}

// handleClient обрабатывает ввод пользователя
func handleClient(bot *tgbotapi.BotAPI, chatID int64, clientName string) {
	client := bonushr.CheckNameCustomer(clientName)
	if client == "alfa" || client == "x5" {
		sendMessage(bot, chatID, "Введите месячную ЗП в гросс, на которую наняли кандидата:")
		userStates[chatID] = "waiting_for_salary_" + client
	} else {
		sendMessage(bot, chatID, "Клиент не найден. Попробуйте снова.")
	}
}

// handleSalary обрабатывает ввод зарплаты и считает бонус
func handleSalary(bot *tgbotapi.BotAPI, chatID int64, salary, state string) {
	var bonus float64

	// TODO: добавить обработку ошибок
	if state == "waiting_for_salary_alfa" {
		bonus = bonushr.CountBonusAlfa(salary)
	} else if state == "waiting_for_salary_x5" {
		bonus = bonushr.CountBonusX5(salary)
	}

	response := fmt.Sprintf("Ваш бонус составляет: %.2f гросс", bonus)
	sendMessage(bot, chatID, response)

	// Сброс состояния
	userStates[chatID] = "waiting_for_client"
}
