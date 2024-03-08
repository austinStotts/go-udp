package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	newuser := User{
		Username: "bobby",
		Password: "aaabbbccc",
	}

	session_token, err := login(newuser)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(session_token)

}

func login(user User) (string, error) {
	// marshal user to json
	str, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// send request
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/account/login", bytes.NewBuffer(str))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// parse response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	session_token := string(body[:])
	return session_token, nil

}

func signup(user User) (string, error) {
	// marshal user to json
	str, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// send request
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/account/signup", bytes.NewBuffer(str))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// parse response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	session_token := string(body[:])
	return session_token, nil
}

func validateToken(user User) (string, error) {
	// marshal user to json
	str, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// send request
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/account/validate-token", bytes.NewBuffer(str))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// parse response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	session_token := string(body[:])
	return session_token, nil

}
