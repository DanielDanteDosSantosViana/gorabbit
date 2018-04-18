package repository

import (
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/queue"
	"gopkg.in/mgo.v2/bson"
)

type QueueRepository interface {
	Store(a *models.Queue) (*models.Queue, error)
	Delete(id bson.ObjectId) error
	ListByBrokerID(id bson.ObjectId) ([]*models.Queue, error)
	DeleteByBrokerID(id bson.ObjectId) error
}
