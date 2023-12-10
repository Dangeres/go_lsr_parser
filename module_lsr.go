package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

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

		flatJson = FlatJson{Name: name, Object: object, Price: price, Size: size, Hash: hash, Time: time.Now().Unix()}

		flats = append(flats, flatJson)
	}

	return flats
}

func getAllFlats() []FlatJson {
	result := make([]FlatJson, 0, 300)

	for i := 1; true; i++ {
		flats := getPageFlats(i)

		if len(flats) == 0 {
			break
		}

		if i == 2 && false {
			break
		}

		result = append(result, flats...)

		fmt.Println(flats)
	}

	return result
}
