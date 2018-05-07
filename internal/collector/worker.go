package collector

type Worker struct {
	sendCommand      chan int
	Cmd string
}


func (q * Worker) AddCmd(cmd string) Command{
	q.Cmd = cmd
	return q
}

func (q * Worker) Execute() error{
	return nil;
}

func newWorker(cmd chan int){
}

func (w *Worker) Sync(){

}


func (w *Worker) Stop(){

}


func (w *Worker) Shutdown(){

}
