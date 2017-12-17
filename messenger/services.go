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
	ACCEPTED_PERCENT = 50
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

type BatchRequest struct {
	Method      string `json:"method"`
	RelativeURL string `json:"relative_url"`
	Body        string `json:"body"`
}

func (s Service) verifyToken() string {
	log.Println("On VerifyToken")

	var challenge string

	// Verify FB messenger request token
	if s.queries.mode == os.Getenv("MODE") && s.queries.verifyToken == os.Getenv("VERIFY_TOKEN") {
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
			// Search feed relate with eventMessage
			groupFeed := getGroupFeed()
			data := analyzeFeed(event.Message.Text, groupFeed.Data)

			// Post to sender
			// callBatchRequest(senderID, data)
			log.Println(data)
			for _, v := range data {
				callSendAPI(senderID, v.Message+"Link:"+v.PermalinkURL)
			}
		}
	}
}

func getGroupFeed() (groupFeed GroupFeed) {
	log.Println("On getGroupFeed")
	// Create GET request
	req, err := http.NewRequest("GET", os.Getenv("GROUP_FEED_ENDPOINT"), nil)
	if err != nil {
		log.Println(err.Error())
	}

	// Add access token query parameter
	values := url.Values{}
	values.Add("access_token", os.Getenv("GROUP_ACCESS_TOKEN"))
	values.Add("fields", "message,permalink_url")
	values.Add("limit", "10")
	req.URL.RawQuery = values.Encode()

	// Issue request
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	// Parse resposne body
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &groupFeed)
	if err != nil {
		log.Println(err.Error())
	}

	return groupFeed
}

func analyzeFeed(eventMessage string, data []Feed) (d []Feed) {
	log.Println("On analyzeFeed")
	var per float64
	for _, feed := range data {
		per = calculatePercent(eventMessage, feed)

		// Choose relation feed
		if per >= ACCEPTED_PERCENT {
			d = append(d, feed)
		}
	}

	return d
}

func calculatePercent(eventMessage string, feed Feed) float64 {
	words := strings.Split(eventMessage, " ")
	var count float64
	for _, w := range words {
		if strings.Contains(feed.Message, w) {
			count++
		}
	}
	// Percent contain eventMessage words
	return (count / float64(len(words))) * 100
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
	req, err := http.NewRequest("POST", os.Getenv("MESSENGER_ENDPOINT"), bytes.NewReader(bodySerialize))
	if err != nil {
		log.Println(err.Error())
	}

	// Add access token query parameter
	values := url.Values{}
	values.Add("access_token", os.Getenv("PAGE_ACCESS_TOKEN"))
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

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

func callBatchRequest(senderID string, data []Feed) { // Construct request body
	// Construct batch request
	batchRequests := []BatchRequest{}
	for _, v := range data {
		// Construct body content
		body := RequestBody{
			Recipient: model.Recipient{ID: senderID},
			Message: model.Message{
				Text: v.Message + "; Link: " + v.PermalinkURL},
		}
		bstr, _ := json.Marshal(body)
		b := BatchRequest{
			Method:      "POST",
			RelativeURL: "me/messages",
			Body:        string(bstr),
		}
		batchRequests = append(batchRequests, b)

	}

	// Construct HTTP POST request
	bodySerialize, _ := json.Marshal(batchRequests)
	batch := "batch=" + string(bodySerialize)
	batch = `batch=[{"method":"POST", "relative_url":"me/messages",
		"body":{"recipient":"` + senderID + `", "message": {"text":"test"}}},
		{"method":"POST", "relative_url":"me/messages",
		"body":{"recipient":"` + senderID + `", "message": {"text":"test"}}}]`
	req, err := http.NewRequest("POST", "https://graph.facebook.com", bytes.NewReader([]byte(batch)))
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(batch)
	// Add access token query parameter
	values := url.Values{}
	values.Add("access_token", os.Getenv("PAGE_ACCESS_TOKEN"))
	req.URL.RawQuery = values.Encode()
	// req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	// Issue request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer res.Body.Close()

	// var result map[string]interface{}
	resBody, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// if err := json.Unmarshal(resBody, &result); err != nil {
	// 	log.Println(err.Error())
	// }
	log.Println(string(resBody))
}
