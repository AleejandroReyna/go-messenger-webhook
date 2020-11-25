package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

type HTTPHandler struct{}

func (h HTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestToken := req.URL.Query().Get("hub.verify_token")
	requestChallenge := req.URL.Query().Get("hub.challenge")
	var data []byte

	if requestToken == getEnvVariable("WEBHOOK_TOKEN") {
		data = []byte(requestChallenge)
		fmt.Println("Accepted")
	} else {
		data = []byte("No Permissions")
		fmt.Println("Denegate")
	}

	res.Write(data)
}

func main() {
	handler := HTTPHandler{}
	http.ListenAndServe(":9000", handler)
}
