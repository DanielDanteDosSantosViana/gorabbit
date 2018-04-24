package collector

import "sync"

type Worker struct {
	destructor sync.Once
	sendCommand      sync.Mutex
	m          sync.Mutex
}
