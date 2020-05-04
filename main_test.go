package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	type getResponse struct {
		Args    map[string]string `json:"args"`
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
		URL     string            `json:"url"`
	}
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bodyBytes))
	var response getResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		panic(err)
	}
	// print struct
	fmt.Printf("%+v\n\n", response)

	// print json struct
	res2B, _ := json.Marshal(response)
	fmt.Println(string(res2B))
}
