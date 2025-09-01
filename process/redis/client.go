package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nico-phil/process/config"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client = nil
	ctx               = context.Background()
)

// QueuedLead represents a lead in the queue - moved here to avoid circular imports
type QueuedLead struct {
	LeadID       string            `json:"lead_id"`
	ListNumber   string            `json:"list_number"`
	WorkspaceID  string            `json:"workspace_id"`
	CampaignID   string            `json:"campaign_id"`
	PhoneNumber  string            `json:"phone_number"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	ZipCode      string            `json:"zip_code"`
	ExtraData    map[string]string `json:"extra_data"`
	QueuedAt     time.Time         `json:"queued_at"`
	CallAttempts int               `json:"call_attempts"`
	CallStatus   string            `json:"call_status"`
}

// InitRedis initiate the redis client
func InitRedis() error {
	addr := config.GetRedisArr()
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.GetRedisPasswrod(),
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("failed to connect t redis: %v", err)
		return fmt.Errorf("failed to connect to redis %q", err)
	}

	log.Printf("connected to redis: %s", pong)

	return nil

}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return rdb
}

// IncrementCallCount increments call in progress
func IncrementCallCount(workspaceID string) (int, error) {
	key := fmt.Sprintf("ws_%s_calls_in_progress", workspaceID)
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment call count for workspace %s, error: %v", workspaceID, err)
	}

	return int(count), nil
}

// DecrementCallCount decrements call in progress
func DecrementCallCount(workspaceID string) (int, error) {
	key := fmt.Sprintf("ws_%s_calls_in_progress", workspaceID)
	count, err := rdb.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement call count for workspace %s, error: %v", workspaceID, err)
	}

	if count <= 0 {
		rdb.Set(ctx, key, 0, time.Hour)
		count = 0
	}

	return int(count), nil
}

// GetCallCount retreive the amount of call in progress
func GetCallCount(workspaceID string) (int, error) {
	key := fmt.Sprintf("ws_%s_calls_in_progress", workspaceID)
	countStr, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve call count: %v", err)
	}

	count, _ := strconv.Atoi(countStr)

	return count, nil
}

// GetQueueLength retrieve length of the queue for a single workspace
func GetQueueLength(workspaceID string) (int, error) {
	key := fmt.Sprintf("ws_%s", workspaceID)
	length, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get length of the queue")
	}

	return int(length), nil
}

// GetWorkspaceQueues returns all workspace queue keys
func GetWorkspaceQueues() ([]string, error) {
	// keys, err := rdb.Keys(ctx, "ws_*").Result()
	keys, _, err := rdb.Scan(ctx, 0, "ws_*", 100).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace queues: %v", err)
	}

	// Filter out call count keys
	var queueKeys []string
	for _, key := range keys {
		if !isCallCountKey(key) {
			queueKeys = append(queueKeys, key)
		}
	}

	return queueKeys, nil
}

// isCallCountKey checks if a key is a call count key
func isCallCountKey(key string) bool {
	return len(key) > 18 && key[len(key)-18:] == "_calls_in_progress"
}

// CacheCampaignRate caches campaign max rate per minute
func CacheCampaignRate(campaignID string, maxRate int) error {
	key := fmt.Sprintf("campaign_%s_max_rate", campaignID)
	err := rdb.Set(ctx, key, maxRate, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to cache campaign rate: %v", err)
	}

	log.Printf("Cached campaign %s max rate: %d", campaignID, maxRate)
	return nil
}

// GetCachedCampaignRate retrieves cached campaign max rate
func GetCachedCampaignRate(campaignID string) (int, error) {
	key := fmt.Sprintf("campaign_%s_max_rate", campaignID)
	rate, err := rdb.Get(ctx, key).Int()
	if err != nil {
		return 0, fmt.Errorf("failed to get cached campaign rate: %v", err)
	}

	return rate, nil
}

// QueueLead inserts lead for a workspace
func QueueLead(workspaceID string, lead QueuedLead) error {
	key := fmt.Sprintf("ws_%s", workspaceID)

	jsonLead, err := json.Marshal(lead)
	if err != nil {
		fmt.Errorf("failed to marshal lead %v", err)
	}
	_, err = rdb.LPush(ctx, key, jsonLead).Result()
	if err != nil {
		return fmt.Errorf("failed to queue lead for workspace: %s : %v", workspaceID, err)
	}

	return nil
}

func DequeueLead(workspaceID string) (*QueuedLead, error) {
	key := fmt.Sprintf("ws_%s", workspaceID)
	leadJson, err := rdb.RPop(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue lead for workspace %s, %v", workspaceID, err)
	}

	var lead QueuedLead
	err = json.Unmarshal([]byte(leadJson), &lead)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarsha lead: %v", err)
	}

	return &lead, nil

}
