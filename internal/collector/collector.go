package collector


import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	repo "github.com/DanielDanteDosSantosViana/gorabbit/internal/event/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
)

type CommandRequest struct {
	Cmd string
}

type Command interface {
	Execute() error
	BuildCommand(cmd CommandRequest) Command
}

type Collector struct {
	workers  map[string]*Worker
	eventRepo repo.EventRepository
}

func NewCollector(sessiondb db.Session) (*Collector,error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/all")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	repository := repo.NewMongoEventRepository(sessiondb)
	return open(repository)
}

func open(eventRepo repo.EventRepository)(*Collector,error){
	return &Collector{eventRepo},nil
}

func (c * Collector)Run(cmd Command) error{
	return cmd.Execute()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}