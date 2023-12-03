package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

const folderFlats string = "flats/"
const folderQueue string = "queue/"

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

func hash_sha256(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func getPageFlats(pageId int) []FlatJson {
	req, err := http.NewRequest("GET", "https://www.lsr.ru/ajax/search_parking.php", nil)

	q := req.URL.Query()

	pageStr := strconv.Itoa(pageId)

	q.Add("object_code", "luchi")
	q.Add("subtype", "2")
	q.Add("building_id", "139")
	q.Add("floor_id", "7664")
	q.Add("page", pageStr)
	q.Add("sort", "")
	q.Add("order", "")

	req.URL.RawQuery = q.Encode()

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.2319 YaBrowser/23.9.0.2319 Yowser/2.5 Safari/537.36")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")

	// fmt.Println(req.URL.String())

	if err != nil {
		os.Exit(1)
	}

	resp, err := soup.Get(req.URL.String())

	var myAnswer ResponseJson

	json.Unmarshal([]byte(resp), &myAnswer)

	doc := soup.HTMLParse(myAnswer.Html)

	links := doc.FindAll("tr", "class", "b-list-parking__desktop-view")

	var flatJson FlatJson
	var flats []FlatJson

	for _, flat := range links {
		name := flat.Find("div", "class", "b-list-parking__number").Text()
		object := flat.Find("div", "class", "b-list-parking__object").Text()
		tds := flat.FindAll("td")

		hash := hash_sha256(name)

		name = strings.TrimSpace(name)
		object = strings.ReplaceAll(strings.TrimSpace(object), "\n", "")
		object = strings.ReplaceAll(object, "                                    ", " ")

		price, _ := strconv.Atoi(strings.ReplaceAll(tds[2].Find("b").Text(), " ", ""))
		size, _ := strconv.ParseFloat(tds[1].Text(), 64)

		flatJson = FlatJson{Name: name, Object: object, Price: price, Size: size, Hash: hash}

		flats = append(flats, flatJson)
	}

	return flats
}

func getAllFlats() []FlatJson {
	result := make([]FlatJson, 0, 300)

	for i := 1; i < 2; i++ {
		flats := getPageFlats(i)

		if len(flats) == 0 {
			break
		}

		result = append(result, flats...)

		fmt.Println(flats)
	}

	return result
}

func saveFlat(flatJson FlatJson) {
	data, _ := json.MarshalIndent(flatJson, "", "    ")

	folder := folderFlats

	if !pathExists(folder) {
		pathCreate(folder)
	}

	os.WriteFile(fmt.Sprintf("%s%s.json", folder, flatJson.Hash), data, 0777)
}

func saveAllFlats(flats []FlatJson) {
	for _, flat := range flats {
		saveFlat(flat)
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

func readFileFlat(fileName string) FlatJson {
	data, _ := os.ReadFile(fileName)

	var flatJson FlatJson

	json.Unmarshal([]byte(data), &flatJson)

	return flatJson
}

func main() {

	fmt.Println(pathFiles(folderFlats))

	// fmt.Println(readFileFlat("flats/7a946d889db4e24bcfba1d4e1a741bb29c1f8559dda00ac1599933af25b614b1.json"))
	data := getAllFlats()

	saveAllFlats(data)
}
