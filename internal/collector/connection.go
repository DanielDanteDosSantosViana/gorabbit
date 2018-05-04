package collector

type Connection struct {
	Workers []*Worker
}

func (c * Connection) Execute() error{
	return nil
}

func (c * Connection) BuildCommand(cmdRequest *CommandRequest) Command{
	return nil
}