package orchestrator

import (
	"fmt"
	"log"

	"github.com/nico-phil/process/db"
)

type Orchestrator struct{}

func New() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Start() {
	log.Printf("starting orchestrator")
	processAllWorkspace()
}

func processAllWorkspace() error {
	// get capaigns or each worksapce
	// get active list for each compaign
	// get leads for each list
	campaigns, err := db.GetCampaigns()
	if err != nil {
		log.Printf("error getting db campaigns %v", err)
		return err
	}

	activeCampaigns := []db.Campaign{}

	for _, c := range campaigns {
		if c.Active {
			activeCampaigns = append(activeCampaigns, c)
		}
	}

	fmt.Printf("activeCampaigns: %+v", activeCampaigns)

	return nil
}
