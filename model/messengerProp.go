package model

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
