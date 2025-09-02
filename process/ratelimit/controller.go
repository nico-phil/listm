package ratelimit

import (
	"fmt"
	"log"
	"time"

	"github.com/nico-phil/process/db"
	"github.com/nico-phil/process/redis"
)

// RateController manages rate limiting for campaigns
type RateController struct {
}

// NewRateController return a new rate contoller
func NewRateController() *RateController {
	return &RateController{}
}

// RateCalculation contains the calculated rate information
type RateCalculation struct {
	CampaignID             string
	WorkspaceID            string
	MaxRatePerMinute       int
	CurrentCallsInProgress int
	CalculatedRate         int
	QueueDepth             int
	TimeWindow             time.Duration
}

// CalculateInjectionRate calculate how many leads to inject for the next 5 minutes
func (rc *RateController) CalculateInjectionRate(campaign db.Campaign) (*RateCalculation, error) {

	// get current calls in progres for this workspace
	currentCalls, err := redis.GetCallCount(campaign.WorkspaceID)
	if err != nil {
		log.Printf("failed to get current call count for workspace %s", campaign.WorkspaceID)
		return nil, fmt.Errorf("failed to get current call count for workspace %s", campaign.WorkspaceID)
	}

	// Get current queue depth
	queueLength, err := redis.GetQueueLength(campaign.WorkspaceID)
	if err != nil {
		log.Printf("failed to get queue length for workspace %s", campaign.WorkspaceID)
		return nil, fmt.Errorf("failed to get queue length for workspace %s", campaign.WorkspaceID)
	}

	// calculate available capacity
	maxRate := campaign.MaxRatePerMin
	if maxRate <= 0 {
		maxRate = 60 // 60 min
	}

	// for 5-minutes window, we want to inject enougth lead for the next 5 minutes
	timeWindow := 5 * time.Minute
	windowMinutes := int(timeWindow.Minutes())

	// total capacity for the time window
	totalCapacity := maxRate * windowMinutes

	//Account for calls already in progress
	availableCapacity := totalCapacity - currentCalls

	//Account for lead already in the queue
	availableCapacity -= queueLength

	// Ensure we don't go negative
	if availableCapacity < 0 {
		availableCapacity = 0
	}

	calculation := &RateCalculation{
		CampaignID:             campaign.ID,
		WorkspaceID:            campaign.WorkspaceID,
		MaxRatePerMinute:       maxRate,
		CurrentCallsInProgress: currentCalls,
		CalculatedRate:         maxRate,
		QueueDepth:             queueLength,
		TimeWindow:             timeWindow,
	}

	log.Printf("Rate calculation for campaign %s: max=%d/min, current=%d, queue=%d, calculated=%d for %v window",
		campaign.ID, maxRate, currentCalls, queueLength, availableCapacity, timeWindow)

	return calculation, nil
}

// CanInjectLeads checks if we can inject more leads based on rate limits
func (rc *RateController) CanInjectLeads(workspaceID string, maxRatePerMinute int) (bool, int, error) {
	// Get current calls in progress
	currentCalls, err := redis.GetCallCount(workspaceID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get current call count: %v", err)
	}

	// Get current queue depth
	queueDepth, err := redis.GetQueueLength(workspaceID)
	if err != nil {
		return false, 0, fmt.Errorf("failed to get queue depth: %v", err)
	}

	// Calculate buffer for next 5 minutes
	bufferCapacity := maxRatePerMinute * 5

	// Total current load
	currentLoad := int(currentCalls + queueDepth)

	// Check if we're under capacity
	canInject := currentLoad < bufferCapacity
	availableCapacity := bufferCapacity - currentLoad

	if availableCapacity < 0 {
		availableCapacity = 0
	}

	log.Printf("Can inject check for workspace %s: current_load=%d, buffer_capacity=%d, can_inject=%v, available=%d",
		workspaceID, currentLoad, bufferCapacity, canInject, availableCapacity)

	return canInject, availableCapacity, nil
}
