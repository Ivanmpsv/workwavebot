package bonushr

import (
	"strconv"
	"strings"
)

var AlfaCustomer = [8]string{"альфа", "альфа-банк", "альфа банк", "альфабанк", "красный банк", "alfa", "alfa-bank", "alfa bank"}
var X5Customer = [6]string{"x5", "x5group", "x5 group", "пятёрочка", "х5", "х5груп"}

// приводим ввод в нижний регистр, проверяем соответвует ли одному из массивов (клиентов)
func CheckNameCustomer(messageText string) string {
	cutomerInLowerCase := strings.ToLower(messageText)
	for _, el := range AlfaCustomer {
		if strings.Contains(cutomerInLowerCase, el) { // Contains - содержание подстроки в строке
			return "alfa"
		}

	}

	cutomerInLowerCase = strings.ToLower(messageText)
	for _, el := range X5Customer {
		if strings.Contains(cutomerInLowerCase, el) {
			return "x5"
		}

	}

	return ""
}

func CountBonusAlfa(salary string) float64 {
	bonus, _ := strconv.ParseFloat(salary, 64)

	bonus = bonus * 12 * 0.12 * 0.3

	return bonus
}

func CountBonusX5(salary string) float64 {
	bonus, _ := strconv.ParseFloat(salary, 64)

	bonus = bonus * 12 * 0.18 * 0.7 * 0.3

	return bonus
}
