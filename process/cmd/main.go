package main

import (
	"fmt"

	"github.com/nico-phil/process/orchestrator"
)

func main() {
	fmt.Println("Hello from process")
	orchestrator := orchestrator.New()
	orchestrator.Start()
}
