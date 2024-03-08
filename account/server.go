package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Username     string
	Password     string
	SessionToken string
	TokenCreated string
	Created      string
}

func hash(input string) string {
	salt := "uranium"
	h := md5.Sum([]byte(salt + input + salt))
	return hex.EncodeToString(h[:])
}

func generateToken256() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		fmt.Println(err)
	}

	return hex.EncodeToString(token), nil
}

func main() {

	fmt.Println("LISTENING ON :8000")

	mux := http.NewServeMux()

	// define what routes do what
	mux.HandleFunc("/account/signup", signupHandler)
	mux.HandleFunc("/account/login", loginHandler)
	mux.HandleFunc("/account/validate-token", validateTokenHandler)

	s := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	s.ListenAndServe()
}

func getUser(username string) (User, error) {
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
		return user, fmt.Errorf("user not found")
	} else {
		err = dynamodbattribute.UnmarshalMap(result.Item, &user)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		fmt.Println("Found user:")
		// fmt.Println(user)
		return user, nil
	}
}

func isUserInTable(username string) bool {
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
		return false
	} else {
		err = dynamodbattribute.UnmarshalMap(result.Item, &user)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		fmt.Println("Found user:")
		// fmt.Println(user)
		return true
	}
}

func putUser(user User) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println(err)
	}

	svc := dynamodb.New(sess)

	tableName := "users"

	u, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      u,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("user put in table")
}

func updateToken(user User) (string, error) {

	// fmt.Println(user)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("could not connect to table")
	}

	svc := dynamodb.New(sess)

	tableName := "users"

	now := time.Now().Local().String()

	t, err := generateToken256()
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("could not generate token")
	} else {
		for i := 0; i < 2; i++ {
			fmt.Println(i)
			if i == 0 {
				inputToken := &dynamodb.UpdateItemInput{
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":t": {
							S: aws.String(t),
						},
					},
					TableName: aws.String(tableName),
					Key: map[string]*dynamodb.AttributeValue{
						"Username": {
							S: aws.String(user.Username),
						},
					},
					ReturnValues:     aws.String("UPDATED_NEW"),
					UpdateExpression: aws.String("set SessionToken = :t"),
				}

				_, err := svc.UpdateItem(inputToken)
				if err != nil {
					fmt.Println(err)
					return "", fmt.Errorf("could not update table")
				}

			} else {
				inputCreated := &dynamodb.UpdateItemInput{
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":r": {
							S: aws.String(now),
						},
					},
					TableName: aws.String(tableName),
					Key: map[string]*dynamodb.AttributeValue{
						"Username": {
							S: aws.String(user.Username),
						},
					},
					ReturnValues:     aws.String("UPDATED_NEW"),
					UpdateExpression: aws.String("set TokenCreated = :r"),
				}

				_, err := svc.UpdateItem(inputCreated)
				if err != nil {
					fmt.Println(err)
					return "", fmt.Errorf("could not update table")
				}
			}
		}

		return t, nil
	}

}

func signupHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("/signup request")

	// get data from request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("failed to read request data")
		data := []byte("signup failed")
		res.WriteHeader(500)
		res.Write(data)
	}

	// parse data into struct
	requestUser := User{}
	json.Unmarshal([]byte(body), &requestUser)

	// hash password
	hashedPassword := hash(requestUser.Password)
	requestUser.Password = hashedPassword

	// check if username already exists
	if !isUserInTable(requestUser.Username) {
		// add user to table
		fmt.Println("adding user to table")

		t, err := generateToken256()
		if err != nil {
			fmt.Println(err)
		}

		now := time.Now().Local().String()

		requestUser.SessionToken = t
		requestUser.TokenCreated = now
		requestUser.Created = now

		putUser(requestUser)
		data := []byte(t)
		res.WriteHeader(200)
		res.Write(data)
	} else {
		// return error
		fmt.Println("user already exists")
		data := []byte("signup failed")
		res.WriteHeader(500)
		res.Write(data)
	}
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("/login request")

	// get data from request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("failed to read request data")
		data := []byte("signup failed")
		res.WriteHeader(500)
		res.Write(data)
	}

	// parse data into struct
	requestUser := User{}
	json.Unmarshal([]byte(body), &requestUser)

	// hash password
	hashedPassword := hash(requestUser.Password)
	requestUser.Password = hashedPassword

	user, err := getUser(requestUser.Username)
	if err != nil {
		fmt.Println(err)
		data := []byte("login failed")
		res.WriteHeader(500)
		res.Write(data)
	}

	if user.Password == requestUser.Password {
		fmt.Println("hashes match")
		newToken, err := updateToken(user)
		if err != nil {
			fmt.Println("could not update token")
			fmt.Println(err)
		}
		data := []byte(newToken)
		res.WriteHeader(200)
		res.Write(data)
	} else {
		fmt.Println("hashes do not match")
		data := []byte("signup failed")
		res.WriteHeader(500)
		res.Write(data)
	}

	// fmt.Println(user)

}

func validateTokenHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("/validateToken request")

	// get data from request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("failed to read request data")
		data := []byte("login failed")
		res.WriteHeader(500)
		res.Write(data)
	}

	// parse data into struct
	requestUser := User{}
	json.Unmarshal([]byte(body), &requestUser)

	user, err := getUser(requestUser.Username)
	if err != nil {
		fmt.Println(err)
	}

	if user.SessionToken == requestUser.SessionToken {
		fmt.Println("tokens match")
		data := []byte(user.SessionToken)
		res.WriteHeader(200)
		res.Write(data)
	} else {
		fmt.Println("tokens do not match")
		data := []byte("login failed")
		res.WriteHeader(500)
		res.Write(data)
	}
}
