package main

import (
	"github.com/nico-phil/process/db"
	"github.com/nico-phil/process/orchestrator"
)

func main() {

	err := db.NewClient()
	if err != nil {
		return
	}

	orchestrator := orchestrator.New()
	orchestrator.Start()
}
