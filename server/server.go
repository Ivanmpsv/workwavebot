package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// структура-проводник-хранилище для компании-клиента
type Client struct {
	Name    string `json:"name"`
	Formula string `json:"formula"`
}

// временная структура для обработки ответа от сервера
type OutputClients struct {
	Status  string          `json:"status"`
	Clients [][]interface{} `json:"clients"` // временный тип для декодирования двумерного массива
}

type DeleteClient struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UpdateFormula struct {
	Name       string `json:"name"`
	NewFormula string `json:"new_formula"`
}

func Get() ([]Client, error) {
	resp, err := http.Get("http://localhost:5000/get_all_clients") // отправляем HTTP-запрос с методом GET
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

	// декодируем временный формат данных
	var out OutputClients
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	// преобразуем массив массивов в массив структур Client
	clients := make([]Client, len(out.Clients))
	for i, clientData := range out.Clients {
		if len(clientData) >= 3 {
			name, okName := clientData[1].(string)
			formula, okFormula := clientData[2].(string)
			if okName && okFormula {
				clients[i] = Client{
					Name:    name,
					Formula: formula,
				}
			} else {
				return nil, fmt.Errorf("invalid client data: %+v", clientData)
			}
		}
	}

	return clients, nil
}

func Post(name string, formula string) (*Client, error) {
	// Создаём клиента
	client := Client{
		Name:    name,
		Formula: formula,
	}

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
	resp, err := http.Post("http://localhost:5000/add_client", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	// Проверяем статус-код ответа
	if resp.StatusCode == 409 {
		return nil, fmt.Errorf("error 409: client already exists")
	}
	if resp.StatusCode != 201 {
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

func Delete(name string) error {
	url := fmt.Sprintf("http://localhost:5000/delete_client/%s", name)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("client '%s' not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deleteResp DeleteClient
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		return fmt.Errorf("error decoding DELETE response: %v", err)
	}

	fmt.Println(deleteResp.Message)
	return nil
}

func Put() {

}

func GetAllClients() ([]string, error) {
	clients, err := Get()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	var varClients []string

	// Выводим список клиентов
	for _, el := range clients {
		varClients = append(varClients, el.Name, el.Formula)
	}

	return varClients, nil
}
