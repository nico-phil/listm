package db

import (
	"context"
	"fmt"
)

// GetCampaigns retrive all campaign from the database
func GetCampaigns() ([]Campaign, error) {
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
			return []Campaign{}, err
		}

		campaigns = append(campaigns, campaign)
	}

	if err := scanner.Err(); err != nil {
		return []Campaign{}, err
	}

	return campaigns, nil
}

// GetLists retrive all lists fro the database
func GetLists() ([]List, error) {
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
			return nil, err
		}

		lists = append(lists, list)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

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
