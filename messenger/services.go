package messenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"pikachu-bot/model"
)

const (
	// Load environment variables
	FACEBOOK_ENDPOINT = os.Getenv("FACEBOOK_ENDPOINT")
	ACCESS_TOKEN      = os.Getenv("ACCESS_TOKEN")
	MODE              = os.Getenv("MODE")
	VERIFY_TOKEN      = os.Getenv("VERIFY_TOKEN")
	ACCEPTED_PERCENT  = 50
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

type GroupFeed struct {
	Data []Feed `json:"data"`
}

type Feed struct {
	Message      string `json:"message"`
	PermalinkURL string `json:"permalink_url"`
}

type RequestBody struct {
	Recipient model.Recipient `json:"recipient"`
	Message   model.Message   `json:"message"`
}

func (s Service) verifyToken() string {
	log.Println("On VerifyToken")

	var challenge string

	// Verify FB messenger request token
	if s.queries.mode == MODE && s.queries.verifyToken == VERIFY_TOKEN {
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

func searchGroupFeed() {
	groupFeed := getGroupFeed()

}

func getGroupFeed() (groupFeed GroupFeed) {
	// Create GET request
	req, err := http.NewRequest("GET", FACEBOOK_ENDPOINT, nil)
	if err != nil {
		log.Println(err.Error())
	}

	// Add access token query parameter
	values := url.Values
	values.Add("access_token", ACCESS_TOKEN)
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	// Issue request
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	// Parse resposne body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	groupFeed = json.Unmarshal(body, &groupFeed)

	return groupFeed
}

func analyzeFeed(eventMessage string, groupFeed GroupFeed) {
	for _, feed := range groupFeed.Data {

	}
}

func calculateFrequency(eventMessage string, feed Feed) float64 {
	words := strings.Split(eventMessage, " ")
	count := 0
	for _, w := range words {
		if strings.Contains(feed.Message, w) {
			count++
		}
	}

	return count / len(strings.Split(feed.Message, " "))
}

func callSendAPI(senderID, text string) {
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

	// Add access token query parameter
	values := url.Values{}
	values.Add("access_token", ACCESS_TOKEN)
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	// Issue request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	// var result map[string]interface{}
	// resBody, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// if err := json.Unmarshal(resBody, &result); err != nil {
	// 	log.Println(err.Error())
	// }
	// log.Print(result)
}
