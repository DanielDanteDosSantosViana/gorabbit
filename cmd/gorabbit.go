package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/parnurzeal/gorequest"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
)

var (
	ErrUsage          = errors.New("usage")
	ErrUnknownCommand = errors.New("unknown command")
	ErrNameRequired   = errors.New(" -n (name of target) required")
	ErrTargetRequired = errors.New(" -t (target) required")

	ErrNotFoundCurrentUser = errors.New("not found current user")
)

func main() {
	m := NewMain()

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

	case "target":
		return newTargetCommand(m).Run(args[1:]...)

	case "broker":
		return newBrokerCommand(m).Run(args[1:]...)

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
    target     add target daemon local or remote for gorabbit.
    broker     run command to manipulate broker.
    info       info about project
    help       print this screen
Use "gorabbit [command] -h" for more information about a command.
`, "\n")
}

type TargetCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newTargetCommand(m *Main) *TargetCommand {
	return &TargetCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *TargetCommand) Run(args ...string) error {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}
	switch args[0] {
	case "help":
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage

	case "add":
		return newAddTargetCommand(cmd).Run(args[1:]...)

	case "remove":
		return newRemoveTargetCommand(cmd).Run(args[1:]...)

	case "list":
		return newListTargetCommand(cmd).Run(args[1:]...)

	default:
		return ErrUnknownCommand

	}

	return nil
}

type AddTargetCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newAddTargetCommand(m *TargetCommand) *AddTargetCommand {
	return &AddTargetCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *AddTargetCommand) Run(args ...string) error {
	if len(args) == 0 {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	if args[0] == "help" {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	options, err := cmd.ParseOptions(args)

	if err != nil {
		return err
	}

	db, err := createConnectionDB()
	if err != nil {
		return err
	}

	target := &Target{Name: options.Name, Host: options.Host, Port: options.Port}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("target"))
		if err != nil {
			err := fmt.Errorf("could not encode target %s: %s", target.Name, err)
			fmt.Print(err.Error())
			return err
		}
		enc, err := target.encode()
		if err != nil {

			err := fmt.Errorf("could not encode target %s: %s", target.Name, err)
			fmt.Print(err.Error())
			return err
		}
		err = bucket.Put([]byte(target.Name), enc)
		return err
	})

	return err
}

func (m *AddTargetCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target add [options]

Add' is required to save local or remote gorabbit
information to send commands when needed.

The options are:
    -n    Name of target.
    -p    Port of target process.
    -h    Host of target local or remote, by default localhost.
    help  print this screen

`, "\n")
}

func (cmd *AddTargetCommand) ParseOptions(args []string) (*Options, error) {
	var options Options

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.Host, "h", "localhost", " host of target")
	fs.StringVar(&options.Name, "n", "", "name of target")
	fs.StringVar(&options.Port, "p", "2222", " port of target")
	fs.SetOutput(cmd.Stderr)
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if options.Name == "" {
		return &options, ErrNameRequired
	}

	return &options, nil
}

type ListTargetCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newListTargetCommand(t *TargetCommand) *ListTargetCommand {
	return &ListTargetCommand{
		Stdin:  t.Stdin,
		Stdout: t.Stdout,
		Stderr: t.Stderr,
	}
}

func (cmd *ListTargetCommand) Run(args ...string) error {
	if len(args) > 0 {
		if args[0] == "help" {
			fmt.Fprintln(cmd.Stderr, cmd.Usage())
			return ErrUsage
		}
	}

	db, err := createConnectionDB()
	if err != nil {
		return err
	}

	target := &Target{}

	db.View(func(tx *bolt.Tx) error {
		fmt.Fprintln(cmd.Stdout, " NAME           HOST       PORT")
		bucket := tx.Bucket([]byte("target"))
		if bucket==nil{
			return nil
		}

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			target, err = target.decode(v)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.Stdout, " %s 	    %s         %s \n", target.Name, target.Host, target.Port)
		}
		return nil
	})

	return nil
}

func (m *ListTargetCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target list

lists all the targets added.

The options are:
    help  print this screen
`, "\n")
}

type RemoveTargetCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newRemoveTargetCommand(m *TargetCommand) *RemoveTargetCommand {
	return &RemoveTargetCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *RemoveTargetCommand) Run(args ...string) error {
	if len(args) == 0 {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	if args[0] == "help" {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	options, err := cmd.ParseOptions(args)

	if err != nil {
		return err
	}

	db, err := createConnectionDB()
	if err != nil {
		return err
	}

	target := &Target{Name: options.Name}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("target"))
		if bucket==nil{
			fmt.Fprintln(cmd.Stdout, "not found target ")
			return nil
		}

		err = bucket.Delete([]byte(target.Name))
		if err != nil {
			fmt.Fprintln(cmd.Stdout, "could not delete target with name %s : %s", target.Name, err)
			return err
		}
		return err
	})

	return nil
}

func (m *RemoveTargetCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target remove [options]

remove target by name.

The options are:
    -n    Name of target.
    help  print this screen

`, "\n")
}

func (cmd *RemoveTargetCommand) ParseOptions(args []string) (*Options, error) {
	var options Options

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.Name, "n", "", "name of target")
	fs.SetOutput(cmd.Stderr)
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if options.Name == "" {
		return &options, ErrNameRequired
	}
	return &options, nil
}

type Target struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
}

func (t *Target) encode() ([]byte, error) {
	enc, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (t *Target) decode(data []byte) (*Target, error) {
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type Options struct {
	Name   string
	Host   string
	Port   string
	Target string
}

func (m *TargetCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target [commands]

Target add, remove or list connections local or remote existes in deamon
gorabbit.

The commands are:
    add     add target daemon local or remote for gorabbit.
    remove  remove targets.
    list    list all targets.
    help    print this screen

Use "gorabbit target [commands] -h" for more information about a commands.
`, "\n")
}

type BrokerCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newBrokerCommand(m *Main) *BrokerCommand {
	return &BrokerCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *BrokerCommand) Run(args ...string) error {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}
	switch args[0] {
	case "help":
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage

	case "add":
		return newAddBrokerCommand(cmd).Run(args[1:]...)

	case "remove":
		return newRemoveBrokerCommand(cmd).Run(args[1:]...)

	default:
		return ErrUnknownCommand

	}

	return nil
}

func (cmd *BrokerCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit broker [commands]

broker command is usage to manage conections of brokers( RabbitMQ ).

The commands are:
    add     add brokers on local or remote gorabbit.
    remove  remove brokers.
    list    list all brokers.
    help    print this screen

Use "gorabbit broker [commands] -h" for more information about a commands.
`, "\n")
}

type AddBrokerCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newAddBrokerCommand(b *BrokerCommand) *AddBrokerCommand {
	return &AddBrokerCommand{
		Stdin:  b.Stdin,
		Stdout: b.Stdout,
		Stderr: b.Stderr,
	}
}

func (cmd *AddBrokerCommand) Run(args ...string) error {
	if len(args) == 0 {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	if args[0] == "help" {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	options, err := cmd.ParseOptions(args)
	if err != nil {
		return err
	}

	broker := &Broker{Name: options.Name, Port: options.Port, Host: options.Host}
	db, err := createConnectionDB()
	if err != nil {
		return err
	}

	err = db.View(func(tx *bolt.Tx) error {
		target := &Target{}
		targetB := tx.Bucket([]byte("target")).Get([]byte(options.Target))
		target, err := target.decode(targetB)
		if err != nil {
			return err
		}
		request := gorequest.New()
		uri := "http://" + target.Host + ":" + target.Port + "/v1/brokers"
		resp, _, errs := request.Post(uri).
			Set("Notes", "target is coming!").
			Send(broker).
			End()
		if len(errs) > 0 {
			return errs[0]
		}
		if resp.StatusCode != http.StatusCreated {
			fmt.Fprintln(cmd.Stdout, "error to create broker, check the target are up.\n")
			return nil
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type Broker struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
}

func (m *AddBrokerCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit broker add [options]

add broker in defined target.

The options are:

    -t  target of gorabbit (required).
    -n  name of broker (required).
    -h  host of broker (default=localhost) (required).
    -p  port (default=5672) (required).
     help  print this screen

`, "\n")
}

func (cmd *AddBrokerCommand) ParseOptions(args []string) (*Options, error) {
	var options Options

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.Target, "t", "", "name of target")
	fs.StringVar(&options.Host, "h", "localhost", "host of broker")
	fs.StringVar(&options.Port, "p", "5672", "port of broker")
	fs.StringVar(&options.Name, "n", "", "name of broker")

	fs.SetOutput(cmd.Stderr)
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if options.Name == "" {
		return &options, ErrNameRequired
	}

	if options.Target == "" {
		return &options, ErrTargetRequired
	}
	return &options, nil
}

type RemoveBrokerCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newRemoveBrokerCommand(b *BrokerCommand) *RemoveBrokerCommand {
	return &RemoveBrokerCommand{
		Stdin:  b.Stdin,
		Stdout: b.Stdout,
		Stderr: b.Stderr,
	}
}

func (cmd *RemoveBrokerCommand) Run(args ...string) error {
	if len(args) == 0 {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	if args[0] == "help" {
		fmt.Fprintln(cmd.Stderr, cmd.Usage())
		return ErrUsage
	}

	return nil
}

func (m *RemoveBrokerCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit broker remove id [options]

remove broker by id.

The options are:
    help  print this screen

`, "\n")
}

func createConfigDir() (string, error) {
	userCurrent, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return "", ErrNotFoundCurrentUser
	}
	path := userCurrent.HomeDir + "/.gorabbit"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
	return path, nil
}

func createConnectionDB() (*bolt.DB, error) {
	path, err := createConfigDir()
	if err != nil {
		return nil, err
	}
	path = path + "/gorabbit.db"

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return db, err
	}
	return db, nil
}
