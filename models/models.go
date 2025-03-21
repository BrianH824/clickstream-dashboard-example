package models

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type ClickStreamEvent struct {
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	ArticleID  *string   `json:"article_id,omitempty"` // a pointer so this can be nil,
	EventType  string    `json:"event_type"`           //		for events unrelated to an article
	Category   string    `json:"category"`
	DeviceType string    `json:"device_type"`
}

var articleEventTypes = []string{"page_view", "article_share", "comment", "like"}
var deviceTypes = []string{"mobile", "desktop", "tablet"}

var userUUIDs []string
var loggedInUsers = map[string]bool{}

func CreateUser() string {
	userID := uuid.NewString()
	userUUIDs = append(userUUIDs, userID)
	loggedInUsers[userID] = false
	return userID
}

func selectUser() string {

	// slow down user creation after 500 users
	random := rand.Intn(10)
	if len(userUUIDs) < 500 && random > 6 {
		return CreateUser()
	} else if random > 8 {
		return CreateUser()
	}

	randomIndex := rand.Intn(len(userUUIDs))
	return userUUIDs[randomIndex]
}

func GenerateRandomEvent() ClickStreamEvent {
	userID := selectUser()

	// Check if this user has logged in
	isLoggedIn := loggedInUsers[userID]
	var eventType string
	var articleID *string

	if !isLoggedIn {
		loggedInUsers[userID] = true

		return ClickStreamEvent{
			Timestamp:  time.Now(),
			UserID:     userID,
			EventType:  "login",
			Category:   "admin",
			DeviceType: deviceTypes[rand.Intn(len(deviceTypes))],
		}
	}

	eventType = articleEventTypes[rand.Intn(len(articleEventTypes))]
	articleUUID := uuid.NewString()
	articleID = &articleUUID

	return ClickStreamEvent{
		Timestamp:  time.Now(),
		UserID:     userID,
		ArticleID:  articleID,
		EventType:  eventType,
		Category:   "news",
		DeviceType: deviceTypes[rand.Intn(len(deviceTypes))],
	}
}
