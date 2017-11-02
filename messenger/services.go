package messenger

import (
	"github.com/subosito/gotenv"
	"os"
)

type Queries struct {
	mode        string
	verifyToken string
	challenge   string
}

type Service struct {
	queries Queries
}

func (s Service) VerifyToken() {
	gotenv.Load()

	if s.queries.mode == os.Getenv("MODE") &&
		s.queries.verifyToken == os.Getenv("VERIFY_TOKEN") {
		io.WriteString(w, challenge)
	}

	fmt.Printf(challenge)
	w.WriteHeader(http.StatusOK)
}
