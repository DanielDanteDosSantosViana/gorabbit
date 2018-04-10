package main

import (
	"fmt"
	"strings"
	"io"
	"errors"
	"os"
	"flag"
	"github.com/boltdb/bolt"
	"os/user"
)


var (

	DB  *bolt.DB

	ErrUsage = errors.New("usage")

	ErrNotFoundCurrentUser = errors.New("not found current user")

	ErrUnknownCommand = errors.New("unknown command")

)

func main() {

	m := NewMain()

	if DB !=nil{
		defer DB.Close()
	}
	if err := m.Run(os.Args[1:]...); err == ErrUsage {
		os.Exit(2)
	} else if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type Main struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewMain() *Main {

	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}


func (m *Main) Run(args ...string) error {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Fprintln(m.Stderr, m.Usage())
		return ErrUsage
	}

	switch args[0] {
	case "help":
		fmt.Fprintln(m.Stderr, m.Usage())
		return ErrUsage

	case "init":
		return newInitCommand(m).Run(args[1:]...)

	case "broker":
		return ErrUsage
	case "database":
		return ErrUsage

	default:
		return ErrUnknownCommand
	}
}


func (m *Main) Usage() string {
	return strings.TrimLeft(`

gorabbit is a tool for inspecting data to broker (RabbitMQ).

Usage:
	gorabbit command [arguments]
The commands are:
    init       initialize settings for gorabbit.
    broker     run command to manipulate broker.
    database   manage connections to databases for manipulate data to broker.
    info       info about project
    help       print this screen
Use "gorabbit [command] -h" for more information about a command.
`, "\n")
}


type InitCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

}

func newInitCommand(m *Main) *InitCommand {
	return &InitCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}



func (cmd *InitCommand) Run(args ...string) error {

	db,err := createConnectionDB()
	if err!=nil{
		return err
	}
	DB = db
	return nil
}


type BrokerCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

}


func (cmd *BrokerCommand) Run(args ...string) error {


	return nil
}

func (cmd *BrokerCommand) ParseFlags(args []string) (*BrokerOptions, error) {
	var options BrokerOptions

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.Add, "add", "", "")
	fs.StringVar(&options.Remove, "remove", "", "")
	fs.StringVar(&options.Remove, "list", "", "")
	fs.StringVar(&options.Export, "export", "", "")
	fs.SetOutput(cmd.Stderr)
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return &options, nil
}


type BrokerOptions struct {
	Add   string
	Remove     string
	Export      string
}



func createConfigDir() (string , error){
	userCurrent, err := user.Current();
	if err != nil {
		fmt.Println(err)
		return "",ErrNotFoundCurrentUser
	}
	path := userCurrent.HomeDir+"/.gorabbit"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700);
	}
	return path,nil
}


func createConnectionDB() (*bolt.DB,error){
	path,err := createConfigDir()
	if err!=nil{
		return nil, err
	}
	path = path+"/gorabbit.db"

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return db, err;
	}
	return db,nil
}