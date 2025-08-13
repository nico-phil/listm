package db

func GetCampaigns() ([]Campaign, error) {
	if session == nil {
		return []Campaign{}, ErrNoconnection
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

// func GetLists() {}
