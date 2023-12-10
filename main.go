package main

import (
	"crypto/sha256"
	"fmt"
	"main/jsonFile"
	"os"
)

const folderData string = "data/"
const folderQueue string = "queue/"
const configTelegramFile string = "config_telegram.json"

func hash_sha256(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func saveAllFlats(flats []FlatJson) {
	for _, flat := range flats {
		jsonFile.Save[FlatJson](fmt.Sprintf("%s%s.json", folderQueue, flat.Hash), flat)
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

func pathCreate(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func pathFiles(path string) []string {
	result := make([]string, 0, 1000)

	data, _ := os.ReadDir(path)

	for _, file := range data {
		if !file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return result
}

func indexOf[T comparable](list []T, value T) int {
	for i, v := range list {
		if v == value {
			return i
		}
	}

	return -1
}

func processDiffs() {
	queue := pathFiles(folderQueue)

	data := pathFiles(folderData)

	config := jsonFile.Read[ConfigTelegramJson](configTelegramFile)

	for _, file_queue := range queue {
		queue_data := jsonFile.Read[FlatJson](fmt.Sprintf("%s%s", folderQueue, file_queue))

		if indexOf(data, file_queue) == -1 {
			text := fmt.Sprintf("Новая кладовка %s\nОбьект: %s\nРазмер: %.1f\nЦена: %d", queue_data.Name, queue_data.Object, queue_data.Size, queue_data.Price)

			send_message(TelegramJsonInput{
				Sender: config.Sender,
				Text:   text,
			}, config)

			prices := make([]PriceInfo, 0, 1)
			prices = append(prices, PriceInfo{Price: queue_data.Price, Time: queue_data.Time})

			jsonFile.Save[FlatDataJson](
				fmt.Sprintf("%s%s", folderData, file_queue),
				FlatDataJson{
					Name:   queue_data.Name,
					Object: queue_data.Object,
					Hash:   queue_data.Hash,
					Prices: prices,
					Price:  queue_data.Price,
					Size:   queue_data.Size,
					Time:   queue_data.Time,
				},
			)
		}

		os.Remove(fmt.Sprintf("%s%s", folderQueue, file_queue))
	}

	fmt.Println(queue)
	fmt.Println(data)
}

func main() {
	data := getAllFlats()

	saveAllFlats(data)

	processDiffs()
}
