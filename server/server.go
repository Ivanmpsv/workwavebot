package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// для общения с API, хранилище клиентов
type client struct {
	Name    string `json:"name"`
	Formula string `json:"formula"`
}

// структура с полем статуса и массивом клиентов
type outputClients struct {
	Status  string          `json:"status"`
	Clients [][]interface{} `json:"clients"` // временный тип для декодирования двумерного массива
}

// для общения с API, удал
type deleteClient struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// для общения с API, изменение формулы
type updateFormula struct {
	Name       string `json:"name"`
	NewFormula string `json:"new_formula"`
}

// для общения с API, подсчёт формул
type calculate struct {
	Client_name string  `json:"client_name"`
	Payment     float64 `json:"payment"`
}

// для передачи id админа по API
type admin struct {
	ID int `json:"id"`
}

// для хранения всех админов
type DBadmins struct {
	ID []admin
}

func Get() ([]client, error) {
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
	var out outputClients
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	// преобразуем lдвумерный массив в массив структур Client
	clients := make([]client, len(out.Clients))
	for i, clientData := range out.Clients {
		if len(clientData) >= 3 {
			name, okName := clientData[1].(string)
			formula, okFormula := clientData[2].(string)
			if okName && okFormula {
				clients[i] = client{
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

func PostAddClient(name, formula *string) (*client, error) {
	// Создаём клиента
	client := client{
		Name:    *name,
		Formula: *formula,
	}

	jsonData, err := json.Marshal(client)
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

func PostCalculatePayment(name string, payment float64) (float64, error) {
	cl := calculate{
		Client_name: name,
		Payment:     payment,
	}

	//превращаем экземпляр структуры в json данные
	jsonData, err := json.Marshal(cl)
	if err != nil {
		return 0, err
	}

	// Отправляем POST-запрос на подсчёт премии
	resp, err := http.Post("http://localhost:5000/calculate_payment", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("не удалось отправить POST запрос: %v", err)
	}

	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	// читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать ответ")
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, &cl)
	if err != nil {
		return 0, fmt.Errorf("unmarshal не успешно, не удалось заполнить структуру: %v", err)
	}

	return cl.Payment, nil
}

func PostAddAdmin(id *int) (*admin, error) {

	ad := admin{
		ID: *id,
	}

	data, err := json.Marshal(ad)
	if err != nil {
		return nil, fmt.Errorf("ошибка, не удалось json.Marshal %v", err)
	}

	resp, err := http.Post("http://localhost:5000/add_admin", "application/json", bytes.NewBuffer(data))

	if err != nil {
		return nil, fmt.Errorf("не удалось сделать Post запрос %v", err)
	}

	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать тело ответа %v", err)
	}

	//заполняем структуру BDadmins
	var dba DBadmins

	err = json.Unmarshal(body, &dba)
	if err != nil {
		return nil, fmt.Errorf("не удалось json.Unmarshal %v", err)
	}

	return &ad, nil
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

	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deleteResp deleteClient
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		return fmt.Errorf("error decoding DELETE response: %v", err)
	}

	return nil
}

// /remove_admin/<id>
func DeleteAdmin(id string) error {
	url := fmt.Sprintf("http://localhost:5000/remove_admin/%s", id)

	// для DELETE не нужно body т.к удаление идёт по уникальному url
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("cant newRequest")
	}

	defer req.Body.Close() // закрываем тело ответа после выполнения

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error StatusCode: %v", resp.StatusCode)
	}

	return nil
}

func Put(name, formula string) error {
	// записываем имя и формулу клиента для обновления (по аналогии с Post)
	body := updateFormula{
		Name:       name,
		NewFormula: formula,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", "http://localhost:5000/update_client_formula", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error PUT request: %v", err)
	}

	//в случае с запросом PUT через http.NewRequest нужно дополнительно указать, что мы передаём JSON
	req.Header.Set("Content-Type", "application/json")

	// HTTP-запрос с использованием встроенного HTTP-клиента, норм если не нужна кастомизация
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error resp: %v", err)
	}

	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	// проверяем статус код
	if resp.StatusCode != 200 {
		fmt.Printf("Status code: %d", resp.StatusCode)
		return err
	}

	// заполнять структуру не нужно, т.к UpdateFormula нужна для передачи, а не хранения данных

	return nil

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
