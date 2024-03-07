package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func getUser(username string) User {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println(err)
	}

	svc := dynamodb.New(sess)

	tableName := "users"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	user := User{}
	if result.Item == nil {
		fmt.Println("result is empty")
		return user
	} else {
		err = dynamodbattribute.UnmarshalMap(result.Item, &user)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		fmt.Println("Found user:")
		fmt.Println(user)
		return user
	}
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
	requestUser := User{}
	json.Unmarshal([]byte(body), &requestUser)

	// hash password
	hashedPassword := hash(requestUser.Password)
	fmt.Println(hashedPassword)

	// check if username already exists
	foundUser := getUser(requestUser.Username)
	fmt.Println(foundUser)

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
