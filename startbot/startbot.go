package startbot

import (
	"fmt"
	"workwavebot/bonushr"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO - сделать, чтобы бот работал и в чатах
func StartBot(bot *tgbotapi.BotAPI) {
	Createbot()

	u := tgbotapi.NewUpdate(0) //бот будет начинать с последнего непрочитанного обновления.
	// ждём 60sec чтобы получить новые обновления от Telegram API. Это снижает нагрузку на сервер tg
	u.Timeout = 60

	//канал который будет получать все входящие сообщения и другие события (обновления) от Telegram
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			// Приветствие пользователя при команде /start
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Локтар огар! Напишите имя клиента чтобы рассчитать вашу премию [тест версия, только альфа и X5]")
				bot.Send(msg)
			} else if bonushr.CheckNameCustomer(update.Message.Text) == "alfa" {
				// Если название клиента найдено, ожидаем сумму от пользователя
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите месячную ЗП в гросс, на которую наняли кандидата: ")
				bot.Send(msg)
			} else if bonushr.CheckNameCustomer(update.Message.Text) == "x5" {
				// Если название клиента найдено, ожидаем сумму от пользователя
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите месячную ЗП в гросс, на которую наняли кандидата: ")
				bot.Send(msg)

			} else {
				// Если сумма уже введена, считаем бонус
				bonus := bonushr.CountBonusAlfa(update.Message.Text)

				// Формируем и отправляем ответ с рассчитанным бонусом
				response := fmt.Sprintf("Ваш бонус составляет: %.2f гросс", bonus)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
				bot.Send(msg)
			}
		}
	}

}
