package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"os"
	"os/user"
	"strings"
)

var (
	ErrUsage               = errors.New("usage")
	ErrUnknownCommand      = errors.New("unknown command")
	ErrNameRequired        = errors.New(" -n (name of target) required")
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
		return newInitCommand(m).Run(args[1:]...)

	case "broker":
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

func newInitCommand(m *Main) *TargetCommand {
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
		return newAddCommand(cmd).Run(args[1:]...)

	case "remove":
		return newRemoveCommand(cmd).Run(args[1:]...)

	case "list":
		return newListCommand(cmd).Run(args[1:]...)

	default:
		return ErrUnknownCommand

	}
	return nil
}

type TargetArgs struct {
}

type TargetResult struct {
}

type AddCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newAddCommand(m *TargetCommand) *AddCommand {
	return &AddCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *AddCommand) Run(args ...string) error {
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

func (m *AddCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target add [options]

Add' is required to save local or remote gorrabit deamon
information to send commands when needed.

The options are:
    -n    Name of target.
    -p    Port of daemon process.
    -h  Host of daemon local or remote, by default localhost.
    help  print this screen

`, "\n")
}

func (cmd *AddCommand) ParseOptions(args []string) (*Options, error) {
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

type ListCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newListCommand(t *TargetCommand) *ListCommand {
	return &ListCommand{
		Stdin:  t.Stdin,
		Stdout: t.Stdout,
		Stderr: t.Stderr,
	}
}

func (cmd *ListCommand) Run(args ...string) error {
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
		c := tx.Bucket([]byte("target")).Cursor()
		fmt.Fprintln(cmd.Stdout, " NAME           HOST       PORT")
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

func (m *ListCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target list

lists all the targets added.

The options are:
    help  print this screen
`, "\n")
}

type RemoveCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newRemoveCommand(m *TargetCommand) *RemoveCommand {
	return &RemoveCommand{
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *RemoveCommand) Run(args ...string) error {
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

		err := bucket.Delete([]byte(target.Name))
		if err != nil {
			fmt.Fprintln(cmd.Stdout, "could not delete target with name %s : %s", target.Name, err)
			return err
		}
		return err
	})

	return nil
}

func (m *RemoveCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit target remove [options]

remove target by name.

The options are:
    -n    Name of target.
    help  print this screen

`, "\n")
}

func (cmd *RemoveCommand) ParseOptions(args []string) (*Options, error) {
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
	Name string
	Host string
	Port string
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
    help       print this screen

Use "gorabbit target [commands] -h" for more information about a commands.
`, "\n")
}

type BrokerCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (cmd *BrokerCommand) Run(args ...string) error {

	return nil
}

func (cmd *BrokerCommand) Usage() string {
	return strings.TrimLeft(`

Usage: gorabbit broker [commands]

Target add, remove or list connections local or remote existes in deamon
gorabbit.

The commands are:
    add     add target daemon local or remote for gorabbit.
    remove  remove targets.
    list    list all targets.
    help       print this screen

Use "gorabbit target [commands] -h" for more information about a commands.
`, "\n")
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
	Add    string
	Remove string
	Export string
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
