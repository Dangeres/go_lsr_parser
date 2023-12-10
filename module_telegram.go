package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func send_message(data TelegramJsonInput, config ConfigTelegramJson) {
	url := fmt.Sprintf("%s/message", config.Host)

	dataJSON, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJSON))

	req.Header.Add("token", config.Token)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
}

func test(name string) {
	test := func(text string) bool {
		fmt.Println(text)
		return true
	}

	test(name)
}
