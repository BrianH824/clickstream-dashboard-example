package models

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type ClickStreamEvent struct {
	Timestamp       time.Time `json:"timestamp"`
	UserID          string    `json:"user_id"`
	ArticleID       *string   `json:"article_id,omitempty"` // a pointer so this can be nil,
	ArticleCategory string    `json:"category"`             //		for events unrelated to an article
	EventType       string    `json:"event_type"`
	DeviceType      string    `json:"device_type"`
}

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

	// return an existing user
	randomIndex := rand.Intn(len(userUUIDs))
	return userUUIDs[randomIndex]
}

func GenerateRandomEvent() ClickStreamEvent {
	userID := selectUser()

	var (
		eventType string
		articleID *string
	)

	if !loggedInUsers[userID] {
		loggedInUsers[userID] = true

		return ClickStreamEvent{
			Timestamp:  time.Now(),
			UserID:     userID,
			EventType:  "login",
			DeviceType: deviceTypes[rand.Intn(len(deviceTypes))],
		}
	} else {
		eventType = randomWeightedEventType()
		articleUUID := uuid.NewString()
		articleID = &articleUUID
	}

	return ClickStreamEvent{
		Timestamp:       time.Now(),
		UserID:          userID,
		ArticleID:       articleID,
		ArticleCategory: randomCategory(),
		EventType:       eventType,
		DeviceType:      deviceTypes[rand.Intn(len(deviceTypes))],
	}
}

func randomWeightedEventType() string {
	pick := rand.Intn(100)
	switch {
	case pick > 90:
		return "article_share"
	case pick > 80:
		return "comment"
	case pick > 50:
		return "like"
	default:
		return "page view"
	}
}

func randomCategory() string {
	categories := []string{"sports", "politics", "tech", "health", "entertainment"}
	return categories[rand.Intn(len(categories))]
}
