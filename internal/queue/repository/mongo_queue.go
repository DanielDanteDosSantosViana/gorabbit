package repository

import (
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/enviroment"
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/queue"
	"gopkg.in/mgo.v2/bson"
)

var queue_collection = "queue"

type mongoQueueRepository struct {
	SessionDB db.Session
}

func NewMongoQueueRepository(session db.Session) QueueRepository {

	return &mongoQueueRepository{session}
}

func (q *mongoQueueRepository) Store(queue *models.Queue) (*models.Queue, error) {
	session := q.SessionDB.Clone()
	defer session.Close()

	queue.ID = bson.NewObjectId()
	collection := q.getCollection()
	err := collection.Insert(queue)

	return queue, err
}

func (q *mongoQueueRepository) Delete(id bson.ObjectId) error {
	session := q.SessionDB.Clone()
	defer session.Close()

	collection := q.getCollection()
	err := collection.Remove(bson.M{"_id": id})

	return err
}

func (q *mongoQueueRepository) ListByBrokerID(id bson.ObjectId) ([]*models.Queue, error) {
	session := q.SessionDB.Clone()
	defer session.Close()

	queues := make([]*models.Queue, 0)

	collection := q.getCollection()

	err := collection.Find(bson.M{"broker_id": id}).All(&queues)

	return queues, err
}

func (q *mongoQueueRepository) DeleteByBrokerID(id bson.ObjectId) error {
	session := q.SessionDB.Clone()
	defer session.Close()

	collection := q.getCollection()

	err := collection.Remove(bson.M{"broker_id": id})

	return err
}

func (q *mongoQueueRepository) getCollection() db.Collection {
	return q.SessionDB.DB(enviroment.Conf.Db.Name).C(queue_collection)
}
