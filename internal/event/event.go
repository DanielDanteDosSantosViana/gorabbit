package event

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Event struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Message      string        `json:"message" validate:"required"`
	Error    string		`json:"error" validate:"required"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedAt time.Time     `json:"created_at"`
}