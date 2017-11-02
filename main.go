package main

import (
	"encoding/json"
	"fmt"

	"io"
	"io/ioutil"
	"net/http"

	"messenger"
)

const (
	TOKEN        = ""
	VERIFY_TOKEN = "verify_token"
)

func main() {
	http.HandleFunc("/webhook", messenger.Webhook)
	http.ListenAndServe(":8080", nil)
}
