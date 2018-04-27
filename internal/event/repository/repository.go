package repository

import (
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/event"
	"gopkg.in/mgo.v2/bson"
)

type EventRepository interface {
	Store(a *models.Event) (*models.Event, error)
	Delete(id bson.ObjectId) error
	List() ([]*models.Event, error)
	Get(id bson.ObjectId) (*models.Event, error)
}
