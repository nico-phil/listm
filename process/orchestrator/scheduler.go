package orchestrator

import (
	"context"
	"fmt"
	"log"
	"time"

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

func processWorkspaceWithContext(ctx context.Context, worksapceID string, campaigns []db.Campaign) error {

	activeCampgaignWithSchedule := getActiveCampignsWithSchedule(worksapceID, campaigns)
	fmt.Printf("activeCampgaignWithSchedule for %s %v\n", worksapceID, activeCampgaignWithSchedule)

	if len(activeCampgaignWithSchedule) == 0 {
		log.Printf("no active campaign found for workspace %s", worksapceID)
	}

	log.Printf("found %d active campaigns for %s", len(activeCampgaignWithSchedule), worksapceID)

	for _, campaign := range activeCampgaignWithSchedule {
		// process the campaign for the worskspace
		fmt.Println(campaign)

	}
	return nil
}

func getActiveCampignsWithSchedule(worksapceID string, campaigns []db.Campaign) []db.Campaign {

	campaignsWithSchedule := []db.Campaign{}

	for _, campaign := range campaigns {
		currentTime := time.Now()
		currentHour := currentTime.Hour()
		currentWeekDay := currentTime.Weekday()

		if currentHour >= campaign.DialStartHour && currentHour <= campaign.DialEndHour && contains(currentWeekDay, campaign.DialDays) {
			campaignsWithSchedule = append(campaignsWithSchedule, campaign)
		}
	}

	return campaignsWithSchedule
}

func contains(currentWeekDay time.Weekday, days []int) bool {
	for _, day := range days {
		if int(currentWeekDay) == day {
			return true
		}
	}

	return false
}
