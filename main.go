package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	Name string `json:"name"`
}

type Site struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type Device struct {
	Name string    `json:"name"`
	Tags []Tag     `json:"tags"`
	ID   uuid.UUID `json:"id"`
	Site Site      `json:"site"`
}

type graphQLResponse struct {
	Data struct {
		Devices []Device `json:"devices"`
	} `json:"data"`
}

func main() {
	query := map[string]string{
		"query": `{
			devices(role: "spine") {
				name
			tags {
				name
			}
				id
				site {
					name
					id
				}
			}
		}
    `,
	}
	queryJson, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", "https://demo.nautobot.com/api/graphql/", bytes.NewBuffer(queryJson))
	request.Header.Add("Authorization", "Token aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	httpResponse, err := client.Do(request)
	defer httpResponse.Body.Close()
	if err != nil {
		log.Fatalf("The HTTP request failed with error %s\n", err)
	}
	decoder := json.NewDecoder(httpResponse.Body)
	response := &graphQLResponse{}
	err = decoder.Decode(response)
	if err != nil {
		log.Fatalf("Failed to decode reponse json, got error %s\n", err)
	}

	for _, device := range response.Data.Devices {
		fmt.Printf("Device: %+v\n", device)
	}
}
