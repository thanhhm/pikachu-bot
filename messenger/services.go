package messenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"pikachu-bot/model"
)

type Queries struct {
	mode        string
	verifyToken string
	challenge   string
}

type Service struct {
	queries         Queries
	receivedMessage model.ReceivedMessage
}

type RequestBody struct {
	Recipient model.Recipient `json:"recipient"`
	Message   model.Message   `json:"message"`
}

func (s Service) verifyToken() string {
	log.Println("On VerifyToken")

	var challenge string

	// Verify FB messenger request token
	if s.queries.mode == os.Getenv("MODE") &&
		s.queries.verifyToken == os.Getenv("VERIFY_TOKEN") {
		challenge = s.queries.challenge
	}

	return challenge
}

func (s Service) handleMessage() {
	log.Println("On handleMessage")

	messagingEvents := s.receivedMessage.Entry[0].Messaging
	for _, event := range messagingEvents {
		senderID := event.Sender.ID
		if &event.Message != nil && event.Message.Text != "" {
			callSendAPI(senderID, event.Message.Text)
		}
	}
}

func callSendAPI(senderID, text string) {
	// Load environment variables
	FACEBOOK_ENDPOINT := os.Getenv("FACEBOOK_ENDPOINT")
	ACCESS_TOKEN := os.Getenv("ACCESS_TOKEN")

	// Construct request body
	body := RequestBody{
		Recipient: model.Recipient{ID: senderID},
		Message: model.Message{
			Text: text},
	}
	log.Printf("%#v", body)

	// Construct HTTP POST request
	bodySerialize, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", FACEBOOK_ENDPOINT, bytes.NewReader(bodySerialize))
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	// Add access token on query parameter
	values := url.Values{}
	values.Add("access_token", ACCESS_TOKEN)
	req.URL.RawQuery = values.Encode()

	// Issue request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	defer res.Body.Close()
	var result map[string]interface{}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	if err := json.Unmarshal(resBody, &result); err != nil {
		log.Println(err.Error())
	}
	log.Print(result)
}
