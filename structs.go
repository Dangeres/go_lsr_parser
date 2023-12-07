package main

type ResponseJson struct {
	Html string `json:"html"`
}

type FlatJson struct {
	Name   string  `json:"name"`
	Object string  `json:"object"`
	Hash   string  `json:"hash"`
	Price  int     `json:"price"`
	Size   float64 `json:"size"`
}

type ConfigTelegramJson struct {
	Host                string `json:"host"`
	Sender              int    `json:"sender"`
	Token               string `json:"token"`
	AwaitTime           int    `json:"await_time"`
	SendTelegramMessage bool   `json:"send_telegram_message"`
}

type TelegramJsonInput struct {
	Sender int    `json:"sender"`
	Text   string `json:"text"`
	Silent bool   `json:"silent"`
}
