package orchestrator

import "fmt"

type Orchestrator struct{}

func New() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Start() {
	fmt.Println("orchestrator start")
}
