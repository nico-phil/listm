package db

func GetCampaigns() ([]Campaign, error) {
	if session != nil {
		return []Campaign{}, ErrNoconnection
	}
	query := `SELECT id, description  FROM campaigns`

	scanner := session.Query(query).Iter().Scanner()

	campaigns := []Campaign{}

	for scanner.Next() {
		var campaign Campaign
		err := scanner.Scan(
			&campaign.ID,
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
