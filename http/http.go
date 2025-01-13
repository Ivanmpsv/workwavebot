package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Output - структура для представления (читабельного) ответа сервера
type Output struct {
	JSON struct {
		Name string `json:"name"`
		ID   uint32 `json:"id"`
	} `json:"json"`
	URL string `json:"url"`
}

// структура-проводник-хранилище для компаний-клиентов
type Customer struct {
	Name string `json:"name"`
}

func GetCustomer() (out *Output, cus *Customer, err error) {
	resp, err := http.Get("http") //отправляем http с методом GET
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 { // какой статус ответа, если НЕ 200
		fmt.Println(resp.StatusCode)
		return nil, nil, err
	}

	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	body, err := io.ReadAll(resp.Body) // читаем ответ
	if err != nil {
		return nil, nil, err
	}

	// Декодируeм данные в формате JSON и заполняем структур, иначе ответ в виде JSON будет
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("%+v\n", out) // печатаем ответ в виде структуры
	fmt.Println(out.URL)     // печатаем конкретное поле структуры

	fmt.Println(cus.Name)

	return out, cus, nil
}

func PutCustomer() {

}

func DeleteCustomer() {

}
