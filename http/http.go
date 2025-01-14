package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// структура-проводник-хранилище для компании-клиента
type Client struct {
	Name string `json:"name"`
}

// статус ответа сервера + вложение структуры клиентов
type Output struct {
	Status  int       `json:"status"`
	Clients []*Client `json:"clients"`
	Payment float64   `json:"payment"`
	Formula float64   `json:"formula"`
}

const server = "http://localhost:5000/get_all_clients" // константа хранит адрес сервера

func Get() (*Output, error) {
	resp, err := http.Get(server) //отправляем http с методом GET
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) // читаем ответ
	if err != nil {
		return nil, err
	}

	//заполняем структуру
	var out Output
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func Post(name string, formula float64) (*Output, error) {
	// Создаём клиента
	client := Client{Name: name}

	// Создаём JSON-данные
	data := map[string]interface{}{
		"name":    client.Name,
		"formula": formula,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Отправляем POST-запрос
	resp, err := http.Post(server, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending POST request: %v", err)
	}

	defer resp.Body.Close() //закрываем тело ответа после работы с ним

	// Проверяем статус-код ответа
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Заполняем структуру Output
	var out Output
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response JSON: %v", err)
	}

	return &out, nil
}

func Put() {

}

func Delete() {

}

func premain() {
	out, err := Get()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var clients []Client

	for _, el := range out.Clients {
		clients = append(clients, *el)
	}

	fmt.Println(out.Status, "\n", clients)
}
