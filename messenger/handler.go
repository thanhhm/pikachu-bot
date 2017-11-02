package messenger

import (
	"io"
	"io/ioutil"
	"net/http"

	"./model"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	s := newService(r)
	switch r.Method {
	case "GET":
		{
			challenge := s.VerifyToken()
			io.WriteString(w, challenge)
		}
	case "POST":

	}

	if r.Method == "POST" {

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
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err.Error())
			}
			var receivedMessage model.ReceivedMessage
			if err = json.Unmarshal(body, &receivedMessage); err != nil {
				fmt.Println(err.Error()) // TODO log
			}

			return Service{
				receivedMesage: receivedMessage,
			}
		}
	}
}
