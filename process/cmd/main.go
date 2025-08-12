package main

import (
	"log"

	"github.com/nico-phil/process/db"
	"github.com/nico-phil/process/orchestrator"
)

func main() {

	err := db.NewClient()
	if err != nil {
		return
	}

	log.Printf("successfully connected to db")

	orchestrator := orchestrator.New()
	orchestrator.Start()
}
