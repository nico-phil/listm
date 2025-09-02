package main

import (
	"fmt"

	"github.com/nico-phil/process/db"
)

func main() {

	err := db.NewClient()
	if err != nil {
		return
	}

	// orchestrator := orchestrator.New()
	// orchestrator.Start()

	c, _ := db.GetListsByWorkspace("workspace-1")
	fmt.Println(c)
}
