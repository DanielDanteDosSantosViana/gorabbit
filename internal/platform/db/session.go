package db

import (
	"errors"
	"fmt"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/enviroment"
	"gopkg.in/mgo.v2"
)

func NewSession() (Session, error) {
	db, err := mgo.Dial(enviroment.Conf.Db.Mongo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("An error occurred while attempting to open db connection. %v", err))
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("An error occurred while trying to verify db connection. %v", err))
	}

	return MongoSession{db}, nil
}
