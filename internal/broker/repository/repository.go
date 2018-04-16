package repository

import models "github.com/DanielDanteDosSantosViana/gorabbit/internal/broker"

type BrokerRepository interface {
	Store(a *models.Broker) (int64, error)
}
