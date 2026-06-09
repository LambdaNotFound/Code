package interview

import (
	"time"
)

type Offer struct {
	ID        string
	StartTime time.Time
	EndTime   time.Time
	MaxRedeem int
}

type RedeemEvent struct {
	UserID  string
	Time    time.Time
	OfferID string
}

var cutoff = time.Date(2015, 12, 31, 0, 0, 0, 0, time.UTC)

// Phase 1: build per-user per-offer redeem counts from events
func buildRedeemCounts(events []RedeemEvent) map[string]map[string]int {
	// userCounts[userID][offerID] = net redeem count
	userCounts := make(map[string]map[string]int)

	for _, event := range events {
		if _, exists := userCounts[event.UserID]; !exists {
			userCounts[event.UserID] = make(map[string]int)
		}
		userCounts[event.UserID][event.OfferID]++
	}

	return userCounts
}

// Phase 2: for each user, collect offers they can still redeem by cutoff
func getRedeemableOffers(
	userCounts map[string]map[string]int,
	offers map[string]Offer,
) map[string][]Offer {
	result := make(map[string][]Offer)

	for userID, offerCounts := range userCounts {
		for offerID, count := range offerCounts {
			offer, exists := offers[offerID]
			if !exists {
				continue
			}

			// Offer must have started, not yet expired, and user must be under limit
			offerStarted := offer.StartTime.Before(cutoff)
			offerActive := offer.EndTime.After(cutoff)
			underLimit := count < offer.MaxRedeem

			if offerStarted && offerActive && underLimit {
				result[userID] = append(result[userID], offer)
			}
		}
	}

	return result
}
