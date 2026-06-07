package interview

import (
	"fmt"
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
			offerStarted := !offer.StartTime.After(cutoff)
			offerActive := offer.EndTime.After(cutoff)
			underLimit := count < offer.MaxRedeem

			if offerStarted && offerActive && underLimit {
				result[userID] = append(result[userID], offer)
			}
		}
	}

	return result
}

func testOffer() {
	offers := map[string]Offer{
		"o1": {
			ID:        "o1",
			StartTime: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC),
			MaxRedeem: 3,
		},
		"o2": {
			ID:        "o2",
			StartTime: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC),
			MaxRedeem: 2,
		},
		"o3": {
			ID:        "o3",
			StartTime: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC), // expired before cutoff
			MaxRedeem: 5,
		},
	}

	events := []RedeemEvent{
		// user1: redeemed o1 twice → 1 left, redeemed o2 twice → at limit
		{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
		{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
		{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o2"},
		{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o2"},

		// user2: redeemed o1 three times then unredeemed once → count=2, 1 left
		{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
		{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
		{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
		{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},

		// user1: redeemed o3 but it's expired → should not appear
		{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o3"},
	}

	userCounts := buildRedeemCounts(events)
	result := getRedeemableOffers(userCounts, offers)

	for userID, redeemableOffers := range result {
		fmt.Printf("%s can still redeem:\n", userID)
		for _, offer := range redeemableOffers {
			fmt.Printf("  - offer %s (max: %d)\n", offer.ID, offer.MaxRedeem)
		}
	}
}
