package hopper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nico-phil/process/db"
)

// QueueManager manages the hopper  queue system
type QueueManager struct {
}

// NewQueueManager created a new queue manager
func NewQueueManager() *QueueManager {
	return &QueueManager{}
}

// ProcessAllWorkspacesWithContext process all worspaces
func (qm *QueueManager) ProcessAllWorkspacesWithContext(ctx context.Context) error {
	campaigns, err := db.GetAllCampaigns()
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

// ProcessWorkspaceWithContext processes  a single workspace with context
func (qm *QueueManager) ProcessWorkspaceWithContext(ctx context.Context, worksapceID string, campaigns []db.Campaign) error {

	activeCampgaignWithSchedule := qm.GetActiveCampignsWithSchedule(worksapceID, campaigns)
	log.Printf("activeCampgaignWithSchedule for %s %v\n", worksapceID, activeCampgaignWithSchedule)

	if len(activeCampgaignWithSchedule) == 0 {
		log.Printf("no active campaign found for workspace %s", worksapceID)
		return nil
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

// GetActiveCampignsWithSchedule retreives active campaign that are ready to process
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

// ProcessCampaignWithContext processes a single campaign with context
func (qm *QueueManager) ProcessCampaignWithContext(ctx context.Context, campaign db.Campaign) (int, error) {
	log.Printf("processing campaign %s", campaign.ID)

	// get all list for this spcecific campaign
	lists, err := db.GetActiveListByCampaign(ctx, campaign.ID)
	if err != nil {
		log.Printf("failed get lists for campaign: %s with error: %v", campaign.ID, err)
		return 0, fmt.Errorf("failed to get lists for campaign: %s with error: %v", campaign.ID, err)
	}

	if len(lists) == 0 {
		log.Printf("No active lists found for campaign %s", campaign.ID)
		return 0, nil
	}

	leadsCount, err := db.GetLeadsCount(campaign.WorkspaceID)
	if err != nil {
		log.Printf("error")
	}
	fmt.Println("leadsCount:", leadsCount)

	totalLeadsAvailable := 0

	for _, list := range lists {
		count, ok := leadsCount[list.ListNumber]
		if ok {
			totalLeadsAvailable += count
		}
	}

	if totalLeadsAvailable == 0 {
		log.Printf("No dialable lead available for campaign %s", campaign.ID)
		return 0, nil
	}

	//inject lead proportionally across lists
	totalInjected := 0
	remainingToInject := 300
	for _, list := range lists {
		if remainingToInject <= 0 {
			break
		}

		listLeadCount, ok := leadsCount[list.ListNumber]
		if !ok || listLeadCount == 0 {
			continue
		}

		// Calculate proportional share for this list
		proportion := float64(listLeadCount) / float64(totalLeadsAvailable)
		listInjectCount := int(300 * proportion)

		// Ensure we don't exceed remaining capacity
		if listInjectCount > remainingToInject {
			listInjectCount = remainingToInject
		}

		// Ensure minimum of 1 if there are leads and capacity
		if listInjectCount < 1 && listLeadCount > 0 && remainingToInject > 0 {
			listInjectCount = 1
		}

		if listInjectCount > 0 {
			// injectleadfrom list
			injected, err := qm.InjectLeadsFromList(campaign, list, listInjectCount)
			if err != nil {
				log.Printf("failed to inject leads from list %s", list.ListNumber)
				continue
			}

			totalInjected += injected
			remainingToInject += injected
		}
	}

	return totalInjected, nil
}

// InjectLeadsFromList injects leads from list to queue system
func (qm *QueueManager) InjectLeadsFromList(campaign db.Campaign, list db.List, listInjectCount int) (int, error) {
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
