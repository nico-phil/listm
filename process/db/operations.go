package db

import (
	"context"
	"fmt"
	"log"
)

// GetCampaigns retrive all campaign from the database
func GetAllCampaigns() ([]Campaign, error) {
	if session == nil {
		return []Campaign{}, ErrNoConnection
	}
	query := `SELECT id, workspace_id, name, description, active, max_rate_per_min, dial_start_hour, dial_end_hour, dial_days, createdat, modifiedat FROM campaigns`

	scanner := session.Query(query).Iter().Scanner()

	campaigns := []Campaign{}

	for scanner.Next() {
		var campaign Campaign
		err := scanner.Scan(
			&campaign.ID,
			&campaign.WorkspaceID,
			&campaign.Name,
			&campaign.Description,
			&campaign.Active,
			&campaign.MaxRatePerMin,
			&campaign.DialStartHour,
			&campaign.DialEndHour,
			&campaign.DialDays,
			&campaign.CreatedAt,
			&campaign.ModifiedAt,
		)

		if err != nil {
			log.Printf("error reading campaigns: %v", err)
			return []Campaign{}, fmt.Errorf("db: error reading campaigns %v", err)
		}

		campaigns = append(campaigns, campaign)
	}

	if err := scanner.Err(); err != nil {
		return []Campaign{}, fmt.Errorf("db: error reading campaigns %v", err)
	}

	log.Printf("retrieved %d campaigns", len(campaigns))
	return campaigns, nil
}

// GetLists retrives all lists fro the database
func GetAllLists() ([]List, error) {
	if session == nil {
		return []List{}, ErrNoConnection
	}

	query := "SELECT listnumber, campaignid, workspace_id, listname, active, createdat, updatedat FROM lists"
	scanner := session.Query(query).WithContext(context.Background()).Iter().Scanner()

	lists := []List{}
	for scanner.Next() {
		var list List
		err := scanner.Scan(
			&list.ListNumber,
			&list.CampaignID,
			&list.WorkspaceID,
			&list.ListName,
			&list.Active,
			&list.CreatedAt,
			&list.UpdatedAt,
		)

		if err != nil {
			log.Printf("error reading list: %v", err)
			return nil, err
		}

		lists = append(lists, list)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading list: %v", err)
		return nil, err
	}

	log.Printf("retrieved %d lists", len(lists))
	return lists, nil
}

// GetActiveListByCampaign retrives all active lists for a spcecific campaign from the database
func GetActiveListByCampaign(ctx context.Context, campaignID string) ([]List, error) {
	if session == nil {
		return nil, ErrNoConnection
	}
	query := `SELECT listnumber, listname, workspace_id, campaignid, active, createdat, updatedat from lists where campaignid=? AND active=true ALLOW FILTERING`

	var lists []List
	scanner := session.Query(query, campaignID).WithContext(ctx).Iter().Scanner()
	for scanner.Next() {
		var list List
		err := scanner.Scan(
			&list.ListNumber,
			&list.ListName,
			&list.WorkspaceID,
			&list.CampaignID,
			&list.Active,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("db: failed to get lists for campaign: %s : %w", campaignID, err)
		}

		lists = append(lists, list)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("db: failed to close iterator: %s : %w", campaignID, err)
	}

	return lists, nil
}

func GetCampaignsByWorkspace(workspaceID string) ([]Campaign, error) {
	if session == nil {
		return []Campaign{}, ErrNoConnection
	}
	query := `SELECT id, workspace_id, name, description, active, max_rate_per_min, dial_start_hour, dial_end_hour, dial_days, createdat, modifiedat FROM campaigns WHERE workspace_id = ?`

	scanner := session.Query(query, workspaceID).Iter().Scanner()

	campaigns := []Campaign{}

	for scanner.Next() {
		var campaign Campaign
		err := scanner.Scan(
			&campaign.ID,
			&campaign.WorkspaceID,
			&campaign.Name,
			&campaign.Description,
			&campaign.Active,
			&campaign.MaxRatePerMin,
			&campaign.DialStartHour,
			&campaign.DialEndHour,
			&campaign.DialDays,
			&campaign.CreatedAt,
			&campaign.ModifiedAt,
		)

		if err != nil {
			log.Printf("error reading campaigns for workspace %s: %v", workspaceID, err)
			return []Campaign{}, fmt.Errorf("db: error reading campaigns for workspace %s: %v", workspaceID, err)
		}

		campaigns = append(campaigns, campaign)
	}

	if err := scanner.Err(); err != nil {
		return []Campaign{}, err
	}

	log.Printf("db: retrieved %d for workspace %s", len(campaigns), workspaceID)

	return campaigns, nil

}

// GetLeadsCount counts leads for list (listnumber -> count)
func GetLeadsCount(worksapceID string) (map[string]int, error) {
	if session == nil {
		return nil, ErrNoConnection
	}

	query := `SELECT listnumber FROM list_data where workspace_id=? `
	scanner := session.Query(query, worksapceID).Iter().Scanner()

	var listnumbers []string
	for scanner.Next() {
		var listnumber string
		err := scanner.Scan(&listnumber)
		if err != nil {
			log.Printf("error reading listnumber for workspace: %s", worksapceID)
			return nil, fmt.Errorf("error reading listnumber for workspace: %s", worksapceID)
		}

		listnumbers = append(listnumbers, listnumber)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading listnumber for workspace: %s", worksapceID)
		return nil, fmt.Errorf("error reading listnumber for workspace: %s", worksapceID)

	}

	countMap := map[string]int{}
	for _, listnumber := range listnumbers {
		countMap[listnumber]++
	}

	return countMap, nil
}

// GetListsByWorkspace retreive lists for a single workspace
func GetListsByWorkspace(workspaceID string) ([]List, error) {
	if session == nil {
		return nil, ErrNoConnection
	}
	query := `SELECT listnumber, listname, workspace_id, campaignid, active, createdat, updatedat from lists where workspace_id=?`

	var lists []List
	scanner := session.Query(query, workspaceID).Iter().Scanner()
	for scanner.Next() {
		var list List
		err := scanner.Scan(
			&list.ListNumber,
			&list.ListName,
			&list.WorkspaceID,
			&list.CampaignID,
			&list.Active,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			log.Printf("db: error reading lists for workspace: %v", err)
			return nil, fmt.Errorf("db: failed to get lists for workspace: %s : %w", workspaceID, err)
		}

		lists = append(lists, list)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("db: failed to close iterator: %s : %w", workspaceID, err)
	}

	log.Printf("db: retrieved %d list for workspace %s", len(lists), workspaceID)
	return lists, nil
}
