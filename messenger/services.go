package messenger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
	recipient model.Recipient
	message   model.Message
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

	// Construct request body
	body := RequestBody{
		recipient: model.Recipient{ID: senderID},
		message: model.Message{
			Text: text},
	}
	log.Printf("%#v", body)

	// Construct HTTP POST request
	bodySerialize, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", FACEBOOK_ENDPOINT, bytes.NewReader(bodySerialize))
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	// Issue request
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
}
