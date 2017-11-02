package messenger

import (
	"github.com/subosito/gotenv"
	"os"

	"./model"
)

type Queries struct {
	mode        string
	verifyToken string
	challenge   string
}

type Service struct {
	queries        Queries
	receivedMesage model.ReceivedMessage
}

func (s Service) VerifyToken() string {
	var challenge string

	// Load environment variables
	gotenv.Load()

	// Verify FB messenger request token
	if s.queries.mode == os.Getenv("MODE") &&
		s.queries.verifyToken == os.Getenv("VERIFY_TOKEN") {
		challenge = s.queries.challenge
	}

	return challenge
}

func (s Service) handleMessage() {
	messagingEvents := s.receivedMessage.Entry[0].Messaging
	for _, event := range messagingEvents {
		senderID := event.Sender.ID
		if &event.Message != nil && event.Message.Text != "" {
			io.WriteString(w, senderID+event.Message.Text)
		}
	}
}

func callSendAPI() {
}
