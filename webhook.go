package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// InputMessage we get from Messenger
type InputMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Time    int64  `json:"time"`
		Payload struct {
			Title   string `json:"title"`
			payload string `json:"payload"`
		} `json:"payload"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
				Nlp  struct {
					Entities struct {
						Sentiment []struct {
							Confidence float64 `json:"confidence"`
							Value      string  `json:"value"`
						} `json:"sentiment"`
						Greetings []struct {
							Confidence float64 `json:"confidence"`
							Value      string  `json:"value"`
						} `json:"greetings"`
					} `json:"entities"`
					DetectedLocales []struct {
						Locale     string  `json:"locale"`
						Confidence float64 `json:"confidence"`
					} `json:"detected_locales"`
				} `json:"nlp"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

// HTTPHandler exported
type HTTPHandler struct{}

func (h HTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var data []byte
	if req.Method == "GET" {
		requestToken := req.URL.Query().Get("hub.verify_token")
		requestChallenge := req.URL.Query().Get("hub.challenge")

		if requestToken == getEnvVariable("WEBHOOK_TOKEN") {
			data = []byte(requestChallenge)
			fmt.Println("Accepted")
		} else {
			data = []byte("No Permissions")
			fmt.Println("Denegate")
		}

		res.Write(data)
	}

	if req.Method == "POST" {
		data := []byte("Response here!")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed parsing body: %s", err)
			res.WriteHeader(400)
			res.Write([]byte("An error occurred"))
			return
		}

		// Parse message into the Message struct
		var message InputMessage
		err = json.Unmarshal(body, &message)
		if err != nil {
			log.Printf("Failed unmarshalling message: %s", err)
			res.WriteHeader(400)
			res.Write([]byte("An error occurred"))
			return
		}
		webhookEvent := message.Entry[0]
		if webhookEvent.Messaging != nil {
			for _, v := range webhookEvent.Messaging {
				fmt.Println(v.Message.Text)
			}
		}
		res.WriteHeader(200)
		res.Write(data)
	}
}

func main() {
	handler := HTTPHandler{}
	http.ListenAndServe(":9000", handler)
}
