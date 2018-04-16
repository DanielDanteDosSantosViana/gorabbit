package broker

import "time"

type Broker struct {
	ID        int64     `json:"id"`
	Name      string    `json:"title" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
