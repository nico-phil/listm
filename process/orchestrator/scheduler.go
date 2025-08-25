package orchestrator

import (
	"context"
	"log"

	"github.com/nico-phil/process/db"
)

type Orchestrator struct{}

func New() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Start() {
	log.Printf("starting orchestrator")
	processAllWorkspacesWithContext(context.Background())
	// this function will get 5 min of data in the database and put it in redis
}

func processAllWorkspacesWithContext(ctx context.Context) error {
	// get capaigns or each worksapce
	// get active list for each compaign
	// get leads for each list
	campaigns, err := db.GetCampaigns()
	if err != nil {
		log.Printf("failed to get campaign from db %v", err)
		return err
	}

	workspaces := map[string][]db.Campaign{}

	for _, c := range campaigns {
		if c.Active {
			workspaces[c.WorkspaceID] = append(workspaces[c.WorkspaceID], c)
		}
	}

	totalWorkspaces := len(workspaces)
	processedWorkspaces := 0

	for workspaceID, campaign := range workspaces {
		err := processWorkspaceWithContext(ctx, workspaceID, campaign)
		if err != nil {
			log.Printf("failed to process workspace %s: %v", workspaceID, err)
			continue
		}

		log.Printf("successfully process workspace %s", workspaceID)

		processedWorkspaces++
	}
	log.Printf("processed %d/%d workspaces succesfully", processedWorkspaces, totalWorkspaces)
	return nil
}

func processWorkspaceWithContext(ctx context.Context, worksapceID string, campaign []db.Campaign) error {
	return nil
}
