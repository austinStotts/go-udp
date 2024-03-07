package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Username string
	Password string
}

func hash(input string) string {
	h := md5.Sum([]byte(input))
	return hex.EncodeToString(h[:])
}

func main() {

	fmt.Println("LISTENING ON :8000")

	mux := http.NewServeMux()

	mux.HandleFunc("/account/signup", signupHandler)
	mux.HandleFunc("/account/login", loginHandler)

	s := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	s.ListenAndServe()
}

func signupHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("REQUST ON /SIGNUP")

	// get data from request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(res, "can't read body", http.StatusBadRequest)
		return
	}

	// parse data into struct
	newuser := User{}
	json.Unmarshal([]byte(body), &newuser)
	fmt.Println(newuser.Password)

	// hash password
	hashedPassword := hash(newuser.Password)

	// check if username already exists

	data := []byte("signup")
	res.WriteHeader(200)
	res.Write(data)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("REQUST ON /LOGIN")
	data := []byte("login")
	res.WriteHeader(200)
	res.Write(data)
}
