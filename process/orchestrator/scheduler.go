package orchestrator

import (
	"context"
	"log"

	"github.com/nico-phil/process/hopper"
)

type Orchestrator struct {
	queueManager *hopper.QueueManager
}

func New() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Start() {
	log.Printf("starting orchestrator")
	o.queueManager.ProcessAllWorkspacesWithContext(context.Background())
	// this function will get 5 min of data in the database and put it in redis
}
