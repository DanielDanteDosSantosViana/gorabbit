package repository

import (
	models "github.com/DanielDanteDosSantosViana/gorabbit/internal/broker"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
)

type mongoBrokerRepository struct {
	SessionDB db.Session
}

func NewMongoBrokerRepository(session db.Session) BrokerRepository {

	return &mongoBrokerRepository{session}
}

func (m *mongoBrokerRepository) Store(a *models.Broker) (int64, error) {

	return 10, nil
}
