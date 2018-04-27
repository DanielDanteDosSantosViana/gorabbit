package repository

import (
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/event"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/enviroment"
	"gopkg.in/mgo.v2/bson"
)

var event_collection = "event"

type mongoEventRepository struct {
	SessionDB db.Session
}

func NewMongoEventRepository(session db.Session) EventRepository {

	return &mongoEventRepository{session}
}

func (m *mongoEventRepository) Store(event *models.Event) (*models.Event, error) {
	session := m.SessionDB.Clone()
	defer session.Close()

	event.ID = bson.NewObjectId()
	collection := m.getCollection()
	err := collection.Insert(event)

	return event, err
}

func (m *mongoEventRepository) Delete(id bson.ObjectId) error {
	session := m.SessionDB.Clone()
	defer session.Close()

	collection := m.getCollection()
	err := collection.Remove(bson.M{"_id": id})

	return err
}

func (m *mongoEventRepository) List() ([]*models.Event, error) {
	session := m.SessionDB.Clone()
	defer session.Close()

	events := make([]*models.Event, 0)

	collection := m.getCollection()

	err := collection.Find(nil).All(&events)

	return events, err
}

func (m *mongoEventRepository) Get(id bson.ObjectId) (*models.Event, error) {
	session := m.SessionDB.Clone()
	defer session.Close()

	event := &models.Event{}

	collection := m.getCollection()

	err := collection.Find(bson.M{"_id": id}).One(&event)
	return event, err
}
func (r *mongoEventRepository) getCollection() db.Collection {
	return r.SessionDB.DB(enviroment.Conf.Db.Name).C(event_collection)
}
