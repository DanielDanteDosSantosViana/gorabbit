package queue

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Queue struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name      string    `json:"name" validate:"required"`
	BrokerID  bson.ObjectId `bson:"broker_id,omitempty" json:"broker_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
