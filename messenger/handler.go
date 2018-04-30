package messenger

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/thanhhm/pikachu-bot/model"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	s := newService(r)
	switch r.Method {
	case "GET":
		{
			challenge := s.verifyToken()
			io.WriteString(w, challenge)

		}
	case "POST":
		s.handleMessage()
	}

	w.WriteHeader(http.StatusOK)
}

func newService(r *http.Request) Service {
	s := Service{}
	switch r.Method {
	case "GET":
		{
			// Get query parameter
			mode := r.URL.Query().Get("hub.mode")
			verifyToken := r.URL.Query().Get("hub.verify_token")
			challenge := r.URL.Query().Get("hub.challenge")

			s = Service{
				queries: Queries{
					mode:        mode,
					verifyToken: verifyToken,
					challenge:   challenge},
			}
		}
	case "POST":
		{
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err.Error())
			}
			var receivedMessage model.ReceivedMessage
			if err = json.Unmarshal(body, &receivedMessage); err != nil {
				log.Println(err.Error()) // TODO log
			}

			s = Service{
				receivedMessage: receivedMessage,
			}
		}
	}

	return s
}
