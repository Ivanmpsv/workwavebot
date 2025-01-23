package startbot

import (
	"fmt"
	"strings"
	"workwavebot/bonushr"
	"workwavebot/server"

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
		// в разработке

	case userState == "waiting_for_salary_alfa" || userState == "waiting_for_salary_x5":
		// в разработке
	}

	actionClients(bot, message)

}

// кейсы по созданию/обновлению клиентов
func actionClients(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userState := userStates[chatID]

	switch {
	case message.Text == "/new":
		sendMessage(bot, chatID, "Добавить, удалить или обновить (формулу) клиента: отправьте 1 2 3 соответственно "+
			"\n"+"Посмотреть всех клиентов нажмите 4")
		userStates[chatID] = "1 2 3"

	case userState == "1 2 3" && message.Text == "1":
		sendMessage(bot, chatID, "Напишите имя клиента запятая пробел salary и формула расчётов \n"+
			"Пример: Альфа, payment * 12...")
		userStates[chatID] = "wait new client"

	case userStates[chatID] == "wait new client":
		handlePost(bot, chatID, message.Text)

	case userState == "1 2 3" && message.Text == "2":
		sendMessage(bot, chatID, "Обновить формулу клиента в разработке")

	case userState == "1 2 3" && message.Text == "3":
		sendMessage(bot, chatID, "напишите имя клиента для удаления")
		userStates[chatID] = "wait client to delete"

	case userStates[chatID] == "wait client to delete":
		handleDelete(bot, chatID, message.Text)

	case message.Text == "4":
		handleGet(bot, chatID)
	}
}

func handleGet(bot *tgbotapi.BotAPI, chatID int64) {
	clients, err := server.GetAllClients() // Получаем массив строк
	if err != nil {
		fmt.Println("ошибка на стадии handleGet")
	}

	if len(clients) == 0 {
		// Если клиентов нет, отправляем сообщение об этом
		msg := tgbotapi.NewMessage(chatID, "Клиенты отсутствуют.")
		bot.Send(msg)
		return
	}

	// Преобразуем массив строк в одну строку с переносами строки между элементами
	clientsInstring := strings.Join(clients, "\n")

	sendMessage(bot, chatID, clientsInstring)
}

func handlePost(bot *tgbotapi.BotAPI, chatID int64, userInput string) {
	parts := strings.SplitN(userInput, ",", 2)
	if len(parts) != 2 {
		sendMessage(bot, chatID, "Ошибка ввода. Пожалуйста, используйте формат: 'Клиент, формула'")
		return
	}

	clientName := strings.TrimSpace(parts[0])
	formula := strings.TrimSpace(parts[1])

	cl, err := server.Post(clientName, formula)
	if err != nil {
		sendMessage(bot, chatID, fmt.Sprintf("Ошибка при добавлении клиента: %v", err))
		return
	}

	sendMessage(bot, chatID, fmt.Sprintf("Клиент %s успешно добавлен!", cl))
}

func handleDelete(bot *tgbotapi.BotAPI, chatID int64, nameClient string) {

	server.Delete(nameClient)
	sendMessage(bot, chatID, fmt.Sprintf("Клиент %s удалён", nameClient))
}

// ввод клиента
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
