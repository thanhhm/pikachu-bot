package messenger

import (
	"model"
	"net/http"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	s := newService(r)
	switch r.Method {
	case "GET":
		s.VerifyToken()
	case "POST":

	}
	if r.Method == "GET" {

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
}

func newService(r *http.Request) Service {
	switch r.Method {
	case "GET":
		{
			// Get query parameter
			mode := r.URL.Query().Get("hub.mode")
			verifyToken := r.URL.Query().Get("hub.verify_token")
			challenge := r.URL.Query().Get("hub.challenge")

			return Service{
				queries: {
					mode:        mode,
					verifyToken: verifyToken,
					challenge:   challenge},
			}
		}
	case "POST":
		{

		}
	}
}
