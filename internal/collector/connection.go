package collector

import (
	"github.com/streadway/amqp"
	"fmt"
	"log"
)

type Connection struct {
	workers  map[string]*Worker
	Cmd string
	urlConnection string
	conn *amqp.Connection
}

func NewConnection(urlConnection string)*Connection{
  return &Connection{workers:make(map[string]*Worker),urlConnection:urlConnection}
}

func (q * Connection) Execute() error{
	switch q.Cmd {

	case "connect":
		return q.connect();

	case "stop":
		return q.close();

	default:
		return ErrUnknownCommand
	}
}

func (q * Connection) AddCmd(cmd string) Command{
	log.Println(q)
	q.Cmd = cmd
	return q
}

func (c * Connection) connect() error{
	conn, err := amqp.Dial(c.urlConnection)
	failOnError(err, "Failed to connect to RabbitMQ")
	c.conn = conn
	return err
}


func (c * Connection) close() error{
	return c.conn.Close()
}



func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}