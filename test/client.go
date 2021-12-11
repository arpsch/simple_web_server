package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"ws/model"
)

var (
	port   *string
	server string
)

func init() {

	port = flag.String("port", "8080", "the server port for the client to connect, or defaults to :8080")
	flag.Parse()

	server = "http://localhost:" + *port
}

func main() {

	postUser()
	getUser()
	postUser()
}

func postUser() {
	client := &http.Client{}
	user := model.User{
		ID:         "ID1",
		Name:       "ID1Name",
		SignupTime: time.Now(),
	}
	userJson, _ := json.Marshal(&user)

	URL := server + "/users"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(userJson))
	req.SetBasicAuth("idt", "idt123")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	fmt.Printf("Response %v \n", s)
}

func getUser() {
	client := &http.Client{}
	URL := server + "/users/ID1"
	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth("idt", "idt123")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	fmt.Printf("Response %v \n", s)
}
