package main

import (
	"fmt"

	"github.com/nico-phil/process/db"
	"github.com/nico-phil/process/redis"
	"github.com/nico-phil/process/tz"
)

func main() {

	err := db.NewClient()
	if err != nil {
		return
	}

	err = redis.InitRedis()
	if err != nil {
		return
	}

	// orchestrator := orchestrator.New()
	// orchestrator.Start()

	// , _ := db.GetAllCampaigns()

	// result, _ := redis.IncrementCallCount("workspace-1")
	// fmt.Println(result, result)

	fmt.Println(tz.DownLoadZipData())
}
