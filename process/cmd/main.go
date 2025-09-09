package main

import (
	"github.com/nico-phil/process/db"
	"github.com/nico-phil/process/redis"
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

	// const geoNamesZipURL = "http://download.geonames.org/export/zip/US.zip"

	// const dataDir = "/data"

	// const zipFilePath = dataDir + "/US.zip"

	// fmt.Println(tz.DownLoadZipData(context.TODO(), http.Client{}, geoNamesZipURL, zipFilePath))
}
