package model

import "time"

type Message struct {
	Type      int       `json:"type" bson:"type"`
	Username  string    `json:"username" bson:"username"`
	Message   string    `json:"message" bson:"message"`
	Room      string    `json:"room" bson:"room"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
