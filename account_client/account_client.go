package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Username     string
	Password     string
	SessionToken string
	TokenCreated string
	Created      string
}

func main() {

	steve := &User{
		Username: "steve",
		Password: "abc123",
	}

	_, err := json.Marshal(steve)
	if err != nil {
		fmt.Println(err)
		return
	}

	tokenUser := &User{
		Username:     "steve",
		SessionToken: "1e5fa495b52c62fd7bc0456629821cf5c87dfbc036ad7847b84f0eb824266727",
	}

	t, err := json.Marshal(tokenUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/account/tokenlogin", bytes.NewBuffer(t))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Println("sending request")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
