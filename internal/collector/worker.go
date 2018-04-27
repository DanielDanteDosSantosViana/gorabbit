package collector

import "sync"

type Worker struct {
	sendCommand      chan int
}


func newWorker(cmd chan int){
	return &Worker{sendCommand:cmd}
}

func (w *Worker) Sync(){

}


func (w *Worker) Stop(){

}


func (w *Worker) Shutdown(){

}
