package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Username string
	Password string
}

func main() {

	steve := &User{
		Username: "steve",
		Password: "abc123",
	}

	j, err := json.Marshal(steve)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/account/signup", bytes.NewBuffer(j))
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
