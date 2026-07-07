package interview

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testOffers = map[string]Offer{
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
		EndTime:   time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC), // expires before cutoff
		MaxRedeem: 5,
	},
}

var testEvents = []RedeemEvent{
	// user1: redeemed o1 twice (1 slot left), o2 twice (at limit), o3 once (expired)
	{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
	{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
	{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o2"},
	{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o2"},
	{UserID: "user1", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o3"},
	// user2: redeemed o1 four times (over the limit of 3)
	{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
	{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
	{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
	{UserID: "user2", Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), OfferID: "o1"},
}

func Test_buildRedeemCounts(t *testing.T) {
	counts := buildRedeemCounts(testEvents)

	assert.Equal(t, 2, counts["user1"]["o1"], "user1 redeemed o1 twice")
	assert.Equal(t, 2, counts["user1"]["o2"], "user1 redeemed o2 twice")
	assert.Equal(t, 1, counts["user1"]["o3"], "user1 redeemed o3 once")
	assert.Equal(t, 4, counts["user2"]["o1"], "user2 redeemed o1 four times")
}

func Test_buildRedeemCounts_empty(t *testing.T) {
	counts := buildRedeemCounts(nil)
	assert.Empty(t, counts)
}

func Test_getRedeemableOffers_user1(t *testing.T) {
	counts := buildRedeemCounts(testEvents)
	result := getRedeemableOffers(counts, testOffers)

	// user1: o1 has 1 slot left and is active → redeemable
	// user1: o2 is at limit (2/2) → not redeemable
	// user1: o3 expired before cutoff → not redeemable
	assert.ElementsMatch(t, []Offer{testOffers["o1"]}, result["user1"])
}

func Test_getRedeemableOffers_user2(t *testing.T) {
	counts := buildRedeemCounts(testEvents)
	result := getRedeemableOffers(counts, testOffers)

	// user2: o1 redeemed 4 times, limit is 3 → not redeemable
	assert.Empty(t, result["user2"])
}

func Test_getRedeemableOffers_expiredOffer(t *testing.T) {
	counts := map[string]map[string]int{
		"user1": {"o3": 1},
	}
	result := getRedeemableOffers(counts, testOffers)

	// o3 ended 2015-06-01, cutoff is 2015-12-31 → expired, not redeemable
	assert.Empty(t, result["user1"])
}

func Test_getRedeemableOffers_atLimit(t *testing.T) {
	counts := map[string]map[string]int{
		"user1": {"o2": 2}, // exactly at MaxRedeem=2
	}
	result := getRedeemableOffers(counts, testOffers)

	assert.Empty(t, result["user1"])
}

func Test_getRedeemableOffers_unknownOffer(t *testing.T) {
	counts := map[string]map[string]int{
		"user1": {"unknown": 1},
	}
	result := getRedeemableOffers(counts, testOffers)

	assert.Empty(t, result["user1"])
}
