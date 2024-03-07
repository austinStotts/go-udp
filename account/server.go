package main

import "net/http"

func main() {

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
	data := []byte("signup")
	res.WriteHeader(200)
	res.Write(data)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	data := []byte("login")
	res.WriteHeader(200)
	res.Write(data)
}
