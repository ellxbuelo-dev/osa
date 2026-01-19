package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	if err := InitDB(); err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/health", Healthcheck).Methods("GET")
	r.HandleFunc("/hello", Hello).Methods("GET")
	r.HandleFunc("/user/{id}", GetUser).Methods("GET")
	r.HandleFunc("/user", CreateUser).Methods("POST")

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
