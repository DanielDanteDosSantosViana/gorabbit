package collector


import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Collector struct {
	workers  map[uint16]*Worker
	closed int32

}

func NewConn() (*Collector,error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/all")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	return open()
}

func open()(*Collector,error){
	return &Collector{},nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}