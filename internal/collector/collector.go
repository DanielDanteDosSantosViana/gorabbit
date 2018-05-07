package collector


import (
	repo "github.com/DanielDanteDosSantosViana/gorabbit/internal/event/repository"
	"github.com/DanielDanteDosSantosViana/gorabbit/internal/platform/db"
	"errors"
)


var ErrUnknownCommand = errors.New("unknown command")

type CommandRequest struct {
	Cmd string `json:"cmd" validate:"required"`
	Type string `json:"type" validate:"required"`
	connId string `json:"conn_id"`
	workId string `json:"work_id"`
}

type Command interface {
	Execute() error
}

type Collector struct {
	connections map[string]*Connection
	eventRepo repo.EventRepository
}

func NewCollector(sessiondb db.Session) (*Collector,error) {
	repository := repo.NewMongoEventRepository(sessiondb)
	return open(repository)
}

func open(eventRepo repo.EventRepository)(*Collector,error){
	connections := make(map[string]*Connection)
	return &Collector{eventRepo:eventRepo,connections:connections},nil
}

func (c * Collector)AddCollector(id string, urlConnection string){
	c.connections[id] = NewConnection(urlConnection)
}

func (c * Collector)addWorker(){

}

func(c *Collector)buildCommand(cmd * CommandRequest) (Command, error){
	switch cmd.Type {

	case "queue":
		return c.connections[cmd.connId].workers[cmd.workId].AddCmd(cmd.Cmd),nil

	case "broker":
		return c.connections[cmd.connId].AddCmd(cmd.Cmd),nil;

	default:
		return nil, ErrUnknownCommand
	}
}

func (c * Collector)Execute(cmdRequest *CommandRequest) error{
	if cmd,err:= c.buildCommand(cmdRequest);err!=nil{
		return err
	}else{
		return cmd.Execute()
	}
}
