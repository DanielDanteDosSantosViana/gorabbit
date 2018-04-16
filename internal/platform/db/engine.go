package db

import "gopkg.in/mgo.v2"

type Session interface {
	DB(name string) DataLayer
	Close()
	Ping() error
	Clone() Session
}

type MongoSession struct {
	*mgo.Session
}

func (s MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

type DataLayer interface {
	C(name string) Collection
}

type MongoCollection struct {
	*mgo.Collection
}

type Collection interface {
	Find(query interface{}) *mgo.Query
	Count() (n int, err error)
	Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	Update(selector interface{}, update interface{}) error
}

type MongoDatabase struct {
	*mgo.Database
}

func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

func (s MongoSession) Ping() error {
	return s.Session.Ping()
}

func (s MongoSession) Clone() Session {
	session := s.Session.Clone()
	return MongoSession{session}
}
