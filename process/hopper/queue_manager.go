package hopper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nico-phil/process/db"
)

// QueueManager manages the I/O to the queue system
type QueueManager struct {
}

func NewQueueManager() *QueueManager {
	return &QueueManager{}
}
func (qm *QueueManager) ProcessAllWorkspacesWithContext(ctx context.Context) error {
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
		err := qm.ProcessWorkspaceWithContext(ctx, workspaceID, campaign)
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

func (qm *QueueManager) ProcessWorkspaceWithContext(ctx context.Context, worksapceID string, campaigns []db.Campaign) error {

	activeCampgaignWithSchedule := qm.GetActiveCampignsWithSchedule(worksapceID, campaigns)
	fmt.Printf("activeCampgaignWithSchedule for %s %v\n", worksapceID, activeCampgaignWithSchedule)

	if len(activeCampgaignWithSchedule) == 0 {
		log.Printf("no active campaign found for workspace %s", worksapceID)
	}

	log.Printf("found %d active campaigns for %s", len(activeCampgaignWithSchedule), worksapceID)

	for _, campaign := range activeCampgaignWithSchedule {
		// process the campaign for the worskspace
		injected, err := qm.ProcessCampaignWithContext(ctx, campaign)
		if err != nil {

		}
		fmt.Println(injected)

	}
	return nil
}

func (qm *QueueManager) GetActiveCampignsWithSchedule(worksapceID string, campaigns []db.Campaign) []db.Campaign {

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

func (qm *QueueManager) ProcessCampaignWithContext(ctx context.Context, campaign db.Campaign) (int, error) {
	log.Printf("processing campaign %s", campaign.ID)

	// get all list for this spcecific campaign
	lists, err := db.GetActiveListByCampaing(campaign.ID)
	if err != nil {
		log.Printf("failed get lists for: %s with error: %v", err)
		return 0, fmt.Errorf("failed get lists for: %s with error: %v", err)
	}

	activeLists := make(map[string]db.List)
	for _, list := range lists {
		if list.Active {
			activeLists[list.Listnumber] = list
		}
	}
	return 0, nil
}

func contains(currentWeekDay time.Weekday, days []int) bool {
	for _, day := range days {
		if int(currentWeekDay) == day {
			return true
		}
	}

	return false
}
