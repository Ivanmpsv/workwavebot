package startbot

import (
	"fmt"
	"strconv"
	"strings"
	"workwavebot/api"
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
	actionAdmins(bot, message)

}

// кейсы по созданию/обновлению клиентов и добавить админа
func actionClients(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userState := userStates[chatID]

	switch {
	case message.Text == "/new":
		sendMessage(bot, chatID, "Добавить, обновить формулу, удалить клиента: отправьте 1 2 3 соответственно "+
			"\n"+"Посмотреть всех клиентов 4")
		userStates[chatID] = "1 2 3"

	case userState == "1 2 3" && message.Text == "1":
		sendMessage(bot, chatID, "Напишите: имя клиента запятая пробел salary и формула расчётов \n"+
			"Пример: Альфа, payment * 12...")
		userStates[chatID] = "wait new client"

	case userStates[chatID] == "wait new client":
		handlePost(bot, chatID, message.Text)
		userStates[chatID] = "1 2 3"

	case userState == "1 2 3" && message.Text == "2":
		sendMessage(bot, chatID, "Чтобы обновить формулу уже существующего клиента напишите: "+
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

// кейсы по добавлению/удалению админов
func actionAdmins(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userState := userStates[chatID]

	switch {
	case message.Text == "/admin":
		sendMessage(bot, chatID, "Панель управления админами. Нажмите:"+
			"\n"+"1 - Чтобы добавить админа"+
			"\n"+"2 - Чтобы удалить админа"+
			"\n"+"3 - Чтобы посмотреть всех админов")
		userStates[chatID] = "123"

	case userState == "123" && message.Text == "1":
		sendMessage(bot, chatID, "напишите id нового админа")
		userStates[chatID] = "wait id to add"

	case userState == "wait id to add":
		handleAdmin(bot, chatID, message.Text)
		userStates[chatID] = "123"

	case userState == "123" && message.Text == "2":
		sendMessage(bot, chatID, "напишите id админа для удаления")
		userStates[chatID] = "wait id to dealete"

	case userState == "wait id to deleate":
		api.HandleDeleteAdmin(message.Text)
		userStates[chatID] = "123"

	case userState == "123" && message.Text == "3":
		sendMessage(bot, chatID, "Посмотреть всех клиентов в разработке")
		userStates[chatID] = "123"
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

	//приводим к нижнему регистру и удаляем пробелы по бокам
	clientName = strings.ToLower(clientName)
	clientName = strings.TrimSpace(clientName)
	formula = strings.ToLower(formula)
	formula = strings.TrimSpace(clientName)

	cl, err := server.PostAddClient(&clientName, &formula)
	if err != nil {
		sendMessage(bot, chatID, fmt.Sprintf("Ошибка при добавлении клиента: %v", err))
		return
	}

	sendMessage(bot, chatID, fmt.Sprintf("Клиент %s успешно добавлен!", cl))
}

func handleDelete(bot *tgbotapi.BotAPI, chatID int64, nameClient string) {
	//приводим к нижнему регистру и удаляем пробелы по бокам
	nameClient = strings.ToLower(nameClient)
	nameClient = strings.TrimSpace(nameClient)

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

	sendMessage(bot, chatID, fmt.Sprintf("Формула клиента %s обновлена", clientName))

}

func handleAdmin(bot *tgbotapi.BotAPI, chatID int64, id string) {
	// пока что пропускаю ошибку
	num, _ := strconv.Atoi(id)

	//хотел выводить список всех админов, но пока просто сообщение, что успешно добавлен
	_, err := server.PostAddAdmin(&num)
	if err != nil {
		err = fmt.Errorf("произошла ошибка: %v", err)
		strErr := err.Error()
		sendMessage(bot, chatID, strErr)
		return
	}

	id = fmt.Sprintf("администратор с id %s добавлен, будьте аккуратны, бот не умеет проверять id на корректность", id)

	sendMessage(bot, chatID, id)

}

// ввод клиента, рассчёт премии
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
