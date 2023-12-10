package jsonFile

import (
	"encoding/json"
	"os"
)

func Read[T any](filePath string) T {
	var result T

	data, _ := os.ReadFile(filePath)

	json.Unmarshal(data, &result)

	return result
}

func Save[T any](filePath string, fileData T) {
	data, _ := json.MarshalIndent(fileData, "", "    ")

	os.WriteFile(filePath, data, 0777)
}
