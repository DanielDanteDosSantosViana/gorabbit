package repository

import (
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/broker"
	"gopkg.in/mgo.v2/bson"
)

type BrokerRepository interface {
	Store(a *models.Broker) (*models.Broker, error)
	Delete(id bson.ObjectId) error
	List() ([]*models.Broker, error)
	Get(id bson.ObjectId) (*models.Broker, error)
}
