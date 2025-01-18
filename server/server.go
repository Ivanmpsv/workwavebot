package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// структура-проводник-хранилище для компании-клиента
type Client struct {
	Name    string `json:"name"`
	Formula string `json:"formula"`
}

// статус ответа сервера + вложение структуры клиентов
type Output struct {
	Status  int       `json:"status"`
	Clients []*Client `json:"clients"`
}

const serverURL = "http://localhost:5000/get_all_clients" // константа хранит адрес сервера, для Get
const addClientURL = "http://localhost:5000/add_client"   // константа для Post и наверное Put запроса

func Get() (*Output, error) {
	resp, err := http.Get(serverURL) //отправляем http с методом GET
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

func Post(name string, formula string) (*Client, error) {
	// Создаём клиента
	client := Client{
		Name:    name,
		Formula: formula}

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
	resp, err := http.Post(addClientURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending POST request: %v", err)
	}

	defer resp.Body.Close() //закрываем тело ответа после работы с ним

	if resp.StatusCode == 409 {
		return nil, fmt.Errorf("Ошибка 409, клиент уже существует ёпта")
	}

	// Проверяем статус-код ответа
	if resp.StatusCode != 201 || resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Заполняем структуру Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response JSON: %v", err)
	}

	return &client, nil
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
