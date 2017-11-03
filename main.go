package main

import (
	"github.com/subosito/gotenv"
	"net/http"

	"pikachu-bot/messenger"
)

func init() {
	gotenv.Load(".evn")
}

func main() {
	http.HandleFunc("/webhook", messenger.Webhook)
	http.ListenAndServe(":8080", nil)
}
