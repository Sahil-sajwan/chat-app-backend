package model

type Message struct {
	Type     int    `json:"type"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Room     string `json:"room"`
}
