package orchestrator

import (
	"context"
	"log"

	"github.com/nico-phil/process/hopper"
)

// ProcessOrchestrator manages the main process scheduling and coordination
type ProcessOrchestrator struct {
	queueManager *hopper.QueueManager
}

// New create a process orchestrator
func New() *ProcessOrchestrator {
	return &ProcessOrchestrator{}
}

// Start starts orchestration process
func (po *ProcessOrchestrator) Start() {
	log.Printf("starting orchestrator")
	po.queueManager.ProcessAllWorkspacesWithContext(context.Background())
	// this function will get 5 min of data in the database and put it in redis
}
