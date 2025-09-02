package db

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"
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

// GetActiveCampaignsWithSchedule retrives active campaign between a time window
func GetActiveCampaignsWithSchedule(workspaceID string, currentTime time.Time) ([]Campaign, error) {
	if session == nil {
		return []Campaign{}, ErrNoConnection
	}

	query := "SELECT id, workspace_id, name, description, active, max_rate_per_min, dial_start_hour, dial_end_hour, dial_days, createdat, modifiedat FROM campaigns WHERE workspace_id = ? AND active = true ALLOW FILTERING"

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

		currentHour := 3
		currentDay := currentTime.Day()

		if currentHour >= campaign.DialStartHour && currentHour <= campaign.DialEndHour && slices.Contains(campaign.DialDays, currentDay) {
			campaigns = append(campaigns, campaign)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading campaigns for workspace %s: %v", workspaceID, err)
		return []Campaign{}, fmt.Errorf("db: error reading campaigns for workspace %s: %v", workspaceID, err)
	}

	log.Printf("db: retrieved %d for workspace %s", len(campaigns), workspaceID)

	return campaigns, nil

}

func GetDialableLeads(workspaceID, listNumber string, limit int) ([]ListData, error) {

	if session == nil {
		return []ListData{}, ErrNoConnection
	}

	// add limit
	query := "SELECT leadid, listnumber, workspace_id, phonenumber, firstname, lastname, zipcode, extradata, callcount, dialable, inserteddate, lastcalldate, callstatus FROM list_data WHERE workspace_id = ? AND listnumber = ? AND dialable = true LIMIT ? ALLOW FILTERING"

	scanner := session.Query(query, workspaceID, listNumber, limit).Iter().Scanner()

	var leads []ListData
	for scanner.Next() {
		var lead ListData
		err := scanner.Scan(
			&lead.LeadID,
			&lead.ListNumber,
			&lead.WorkspaceID,
			&lead.PhoneNumber,
			&lead.FirstName,
			&lead.LastName,
			&lead.ZipCode,
			&lead.ExtraData,
			&lead.CallCount,
			&lead.Dialable,
			&lead.InsertedDate,
			&lead.LastCallDate,
			&lead.CallStatus,
		)

		if err != nil {
			log.Printf("db: error reading dialable leads fro workspace %s, list %s: %v", workspaceID, listNumber, err)
			return []ListData{}, fmt.Errorf("db: error reading dialable leads fro workspace %s, list %s: %v", workspaceID, listNumber, err)
		}

		leads = append(leads, lead)
	}

	log.Printf("retrieved %d leads for workspace %s, list %s", len(leads), workspaceID, listNumber)
	return leads, nil
}

// UpdateLeadDialStatus updates lead dialable status and call count
func UpdateLeadDialStatus(workspaceID, listNumber, leadID string, dialable bool) error {
	if session == nil {
		return ErrNoConnection
	}

	now := time.Now()
	var query string

	var count int
	query = "select callcount from list_data where workspace_id= ? AND listnumber= ? AND leadid= ?"
	if err := session.Query(query, workspaceID, listNumber, leadID).Scan(&count); err != nil {
		log.Printf("Error getting callcount: %v", err)
		return err
	}

	if dialable {
		// Setting back to dialable, decrement call count
		query = "UPDATE list_data SET dialable = ?, callcount = ? WHERE workspace_id = ? AND listnumber = ? AND leadid = ?"
		if err := session.Query(query, dialable, count-1, workspaceID, listNumber, leadID).Exec(); err != nil {
			log.Printf("Error updating lead status: %v", err)
			return err
		}

	} else {
		// Setting to non-dialable, increment call count and update last call date
		query = "UPDATE list_data SET dialable = ?, callcount = ?, lastcalldate = ? WHERE workspace_id = ? AND listnumber = ? AND leadid = ?"
		if err := session.Query(query, dialable, count+1, now, workspaceID, listNumber, leadID).Exec(); err != nil {
			log.Printf("Error updating lead status: %v", err)
			return err
		}
		return nil
	}

	log.Printf("Updated lead %s dialable status to %v", leadID, dialable)
	return nil
}

// GetLeadCounts returns count of dialable leads per list for a workspace
func GetLeadCounts(workspaceID string) (map[string]int, error) {
	if session == nil {
		return nil, ErrNoConnection
	}

	query := "SELECT listnumber FROM list_data WHERE workspace_id = ? AND dialable = true  ALLOW FILTERING"
	iter := session.Query(query, workspaceID).Iter()

	counts := make(map[string]int)
	var listNumbers []string
	var listNumber string
	// var count int

	for iter.Scan(&listNumber) {
		listNumbers = append(listNumbers, listNumber)
		listNumber = ""
	}

	for _, v := range listNumbers {
		counts[v]++
	}

	if err := iter.Close(); err != nil {
		log.Printf("Error reading lead counts for workspace %s: %v", workspaceID, err)
		return nil, err
	}

	log.Printf("Retrieved lead counts for %d lists in workspace %s", len(counts), workspaceID)
	return counts, nil

}

// TODO: Add GetNonDialableLeadIDs function for complete cleanup of deactivated lists
// This function should return all lead IDs that are marked as non-dialable for a specific list
// func GetNonDialableLeadIDs(workspaceID, listNumber string) ([]string, error) {
//     query := "SELECT leadid FROM list_data WHERE workspace_id = ? AND listnumber = ? AND dialable = false ALLOW FILTERING"
//     // Implementation needed for cleanup/service.go to reset ALL non-dialable leads in deactivated lists
// }

// GetActiveListsByCampaign retrieves active lists for a specific campaign
func GetActiveListsByCampaign(campaignID string) ([]List, error) {
	if session == nil {
		return nil, ErrNoConnection
	}

	query := "SELECT listNumber, listName, campaignId, workspace_id, active, createdAt, updatedAt FROM lists WHERE campaignId = ? AND active = true ALLOW FILTERING"
	iter := session.Query(query, campaignID).Iter()

	var lists []List
	var list List

	for iter.Scan(&list.ListNumber, &list.ListName, &list.CampaignID,
		&list.WorkspaceID, &list.Active, &list.CreatedAt, &list.UpdatedAt) {
		lists = append(lists, list)
		list = List{} // Reset for next iteration
	}

	if err := iter.Close(); err != nil {
		log.Printf("Error reading active lists for campaign %s: %v", campaignID, err)
		return nil, err
	}

	log.Printf("Retrieved %d active lists for campaign %s", len(lists), campaignID)
	return lists, nil
}

// BatchUpdateLeadsDialable updates multiple leads' dialable status
func BatchUpdateLeadsDialable(workspaceID string, listNumber string, leadIDs []string, dialable bool) error {
	if session == nil {
		return ErrNoConnection
	}

	now := time.Now()
	var query string

	if dialable {
		query = "UPDATE list_data SET dialable = ?, callcount = callcount - 1 WHERE workspace_id = ? AND listnumber = ? AND leadid = ?"
	} else {
		query = "UPDATE list_data SET dialable = ?, callcount = callcount + 1, lastcalldate = ? WHERE workspace_id = ? AND listnumber = ? AND leadid = ?"
	}

	for _, leadID := range leadIDs {
		var err error
		if dialable {
			err = session.Query(query, dialable, workspaceID, listNumber, leadID).Exec()
		} else {
			err = session.Query(query, dialable, now, workspaceID, listNumber, leadID).Exec()
		}

		if err != nil {
			log.Printf("Error updating lead %s status: %v", leadID, err)
			return err
		}
	}

	log.Printf("Batch updated %d leads dialable status to %v", len(leadIDs), dialable)
	return nil
}

// UpdateLeadStatus updates lead status and related fields
func UpdateLeadStatus(workspaceID, listNumber, leadID, status string) error {
	if session == nil {
		return ErrNoConnection
	}

	now := time.Now()
	query := "UPDATE list_data SET callstatus = ?, lastcalldate = ? WHERE workspace_id = ? AND listnumber = ? AND leadid = ?"

	if err := session.Query(query, status, now, workspaceID, listNumber, leadID).Exec(); err != nil {
		log.Printf("[%s]: Error updating lead status: %v", workspaceID, err)
		return err
	}

	log.Printf("[%s]: Updated lead %s status to %s", workspaceID, leadID, status)
	return nil
}

// GetLeadByID retrieves a specific lead by ID
func GetLeadByID(workspaceID, listNumber, leadID string) (*ListData, error) {
	if session == nil {
		return nil, ErrNoConnection
	}

	query := "SELECT leadid, listnumber, workspace_id, phonenumber, firstname, lastname, zipcode, extradata, callcount, dialable, inserteddate, lastcalldate, callstatus FROM list_data WHERE workspace_id = ? AND listnumber = ? AND leadid = ?"

	var lead ListData
	if err := session.Query(query, workspaceID, listNumber, leadID).Scan(
		&lead.LeadID, &lead.ListNumber, &lead.WorkspaceID, &lead.PhoneNumber,
		&lead.FirstName, &lead.LastName, &lead.ZipCode, &lead.ExtraData,
		&lead.CallCount, &lead.Dialable, &lead.InsertedDate,
		&lead.LastCallDate, &lead.CallStatus); err != nil {
		log.Printf("Error reading lead %s: %v", leadID, err)
		return nil, err
	}

	log.Printf("Retrieved lead %s", leadID)
	return &lead, nil
}
