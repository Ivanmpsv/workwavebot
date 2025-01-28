package startbot

import (
	"fmt"
	"strconv"
	"strings"
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
		sendMessage(bot, chatID, "Локтар огар!"+
			"Напишите имя клиента запятая и ЗП кандидата чтобы рассчитать премию рекрутера"+
			"\n"+"Пример: Альфа, 300 000")
		userStates[chatID] = "waiting_for_client"

	case userState == "waiting_for_client":
		BonusRecruiter(bot, chatID, message.Text)
	}

	actionClients(bot, message)

}

// кейсы по созданию/обновлению клиентов
func actionClients(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userState := userStates[chatID]

	switch {
	case message.Text == "/new":
		sendMessage(bot, chatID, "Добавить, обновить (формулу), удалить клиента: отправьте 1 2 3 соответственно "+
			"\n"+"Посмотреть всех клиентов нажмите 4")
		userStates[chatID] = "1 2 3"

	case userState == "1 2 3" && message.Text == "1":
		sendMessage(bot, chatID, "Напишите: имя клиента запятая пробел salary и формула расчётов \n"+
			"Пример: Альфа, payment * 12...")
		userStates[chatID] = "wait new client"

	case userStates[chatID] == "wait new client":
		handlePost(bot, chatID, message.Text)
		userStates[chatID] = "1 2 3"

	case userState == "1 2 3" && message.Text == "2":
		sendMessage(bot, chatID, "Чтобы обновить формулу уже существуюшего клиента напишите: "+
			"имя клиента запятая пробел salary и формула расчётов \n"+
			"Пример: Альфа, payment * 12...")
		userStates[chatID] = "wait new formula"

	case userStates[chatID] == "wait new formula":
		handlePut(bot, chatID, message.Text)
		userStates[chatID] = "1 2 3"

	case userState == "1 2 3" && message.Text == "3":
		sendMessage(bot, chatID, "напишите имя клиента для удаления")
		userStates[chatID] = "wait client to delete"

	case userStates[chatID] == "wait client to delete":
		handleDelete(bot, chatID, message.Text)
		userStates[chatID] = "1 2 3"

	case message.Text == "4":
		handleGet(bot, chatID)
		userStates[chatID] = "1 2 3"

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

	cl, err := server.PostAddClient(&clientName, &formula)
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

func handlePut(bot *tgbotapi.BotAPI, chatID int64, updateFormula string) {
	parts := strings.SplitN(updateFormula, ",", 2)
	if len(parts) != 2 {
		sendMessage(bot, chatID, "Ошибка ввода. Пожалуйста, используйте формат: 'Клиент, формула'")
		return
	}

	clientName := strings.TrimSpace(parts[0])
	Newformula := strings.TrimSpace(parts[1])

	err := server.Put(clientName, Newformula)
	if err != nil {
		sendMessage(bot, chatID, "Что-то пошло не так")
		return
	}

	sendMessage(bot, chatID, fmt.Sprintf("Формула килента %s обнавлена", clientName))

}

// ввод клиента
func BonusRecruiter(bot *tgbotapi.BotAPI, chatID int64, calculateSalary string) {
	cs := strings.SplitN(calculateSalary, ",", 2)
	if len(cs) != 2 {
		sendMessage(bot, chatID, "Ошибка ввода. Пожалуйста, используйте формат: 'Клиент, формула'")
		return
	}

	clientName := strings.TrimSpace(cs[0])
	salary := strings.TrimSpace(cs[1])

	salaryFloat, err := strconv.ParseFloat(salary, 64)
	if err != nil {
		sendMessage(bot, chatID, "Ошибка ввода. Попробуйте заново")
		return
	}

	bonusFloat, err := server.PostCalculatePayment(clientName, salaryFloat)

	if err != nil {
		sendMessage(bot, chatID, "Не получилось рассчитать бонус")
		return
	}

	bonusString := strconv.FormatFloat(bonusFloat, 'f', 2, 64)

	sendMessage(bot, chatID, bonusString)
}
