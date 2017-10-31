package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	TOKEN        = ""
	VERIFY_TOKEN = "verify_token"
)

type ReceivedMessage struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	// ID        int64       `json:"id"`
	// Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	// Timestamp int64     `json:"timestamp"`
	Message Message `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MID  string `json:"mid"`
	Seq  int64  `json:"seq"`
	Text string `json:"text"`
}

type Payload struct {
	TemplateType string  `json:"template_type"`
	Text         string  `json:"text"`
	Buttons      Buttons `json:"buttons"`
}

type Buttons struct {
	Type  string `json:"type"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type ButtonMessageBody struct {
	Attachment Attachment `json:"attachment"`
}

type ButtonMessage struct {
	Recipient         Recipient         `json:"recipient"`
	ButtonMessageBody ButtonMessageBody `json:"message"`
}

type SendMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Text string `json:"text"`
	} `json:"message"`
}

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			mode := r.URL.Query().Get("hub.mode")
			verifyToken := r.URL.Query().Get("hub.verify_token")
			challenge := r.URL.Query().Get("hub.challenge")

			if mode == "subscribe" && verifyToken == VERIFY_TOKEN {
				io.WriteString(w, challenge)
			}

			fmt.Printf(challenge)
			w.WriteHeader(http.StatusOK)
		}
		if r.Method == "POST" {

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err.Error())
			}

			var receivedMessage ReceivedMessage
			if err = json.Unmarshal(body, &receivedMessage); err != nil {
				fmt.Println(err.Error())
			}

			messagingEvents := receivedMessage.Entry[0].Messaging
			for _, event := range messagingEvents {
				senderID := event.Sender.ID
				if &event.Message != nil && event.Message.Text != "" {
					io.WriteString(w, senderID+event.Message.Text)
				}
			}

			fmt.Println("event.Message.Text")
			w.WriteHeader(http.StatusOK)

		}
	})
	http.ListenAndServe(":8080", nil)
}
