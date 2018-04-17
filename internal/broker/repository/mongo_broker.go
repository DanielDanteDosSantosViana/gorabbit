package repository

import (
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/broker"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/enviroment"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var broker_collection = "broker"

type mongoBrokerRepository struct {
	SessionDB db.Session
}

func NewMongoBrokerRepository(session db.Session) BrokerRepository {

	return &mongoBrokerRepository{session}
}

func (m *mongoBrokerRepository) Store(broker *models.Broker) (*models.Broker, error) {
	session := m.SessionDB.Clone()
	defer session.Close()

	broker.ID = bson.NewObjectId()
	collection := m.getCollection()
	err := collection.Insert(broker)

	return broker, err
}

func (m *mongoBrokerRepository) Delete(id bson.ObjectId) error {
	session := m.SessionDB.Clone()
	defer session.Close()

	collection := m.getCollection()
	err := collection.Remove(bson.M{"_id": id})

	return  err
}

func (m *mongoBrokerRepository) List() ([]*models.Broker, error) {
	session := m.SessionDB.Clone()
	defer session.Close()

	brokers := make([]*models.Broker, 0)

	collection := m.getCollection()

	err := collection.Find(nil).All(&brokers)

	return brokers, err
}

func (m *mongoBrokerRepository) Get(id bson.ObjectId) (*models.Broker, error) {
	session := m.SessionDB.Clone()
	defer session.Close()

	broker := &models.Broker{}

	collection := m.getCollection()

	err := collection.Find(bson.M{"_id":id}).One(&broker)
	log.Println(broker)
	return broker, err
}
func (r *mongoBrokerRepository) getCollection() db.Collection {
	return r.SessionDB.DB(enviroment.Conf.Db.Name).C(broker_collection)
}
