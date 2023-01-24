package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://localhost:8080/assignment1")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response Data
	json.Unmarshal(body, &response)
	fmt.Printf("%+v\n", response)
}

type Data struct {
	Page        string   `json:"page"`
	Words       []string `json:"words"`
	Percentages struct {
		One   float64 `json:"one"`
		Three int     `json:"three"`
		Two   float64 `json:"two"`
	} `json:"percentages"`
	Special      []interface{} `json:"special"`
	ExtraSpecial []interface{} `json:"extraSpecial"`
}
