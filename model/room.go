package model

import "time"

type Room struct {
	Rname     string    `json:"rname" bson:"rname"`
	Rpass     string    `json:"rpass" bson:"rpass"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
