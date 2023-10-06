package models

type Message struct {
	Data Data `json:"DATA"`
}

type Data struct {
	Rows []interface{} `json:"rows"`
}
