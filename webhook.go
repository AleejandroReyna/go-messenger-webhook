package main

import (
	"fmt"
	"log"
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

func main() {
	token := getEnvVariable("TOKEN")
	fmt.Println(token)
}
