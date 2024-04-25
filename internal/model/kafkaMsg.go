package model

type Message struct {
	Topic   string  `json:"topic"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Email   string `json:"email"`
	Header  string `json:"header"`
	Message string `json:"message"`
}
