package broker

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Broker struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name      string    `json:"name" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
