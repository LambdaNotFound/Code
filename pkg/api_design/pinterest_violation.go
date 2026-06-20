package apidesign

import "sort"

/*
 * You are given a list of violation records: (int pinId, string policy, string date)
 * Given a policy, return the number of unique pins that violated that policy.
 * Given a date range [startDate, endDate], return the number of unique pins that violated any policy in that range.
 * Given a date range [startDate, endDate], return the number of unique pins per policy in that range.
 */

type ViolationRecord struct {
	PinID  int
	Policy string
	Date   string // "YYYY-MM-DD" — lexicographic order == chronological order
}

type ViolationTracker struct {
	policyPins     map[string]map[int]struct{}            // query 1: policy → unique pins
	datePins       map[string]map[int]struct{}            // query 2: date → unique pins
	policyDatePins map[string]map[string]map[int]struct{} // query 3: policy → date → unique pins
	// alternative for query 3: date → policy → unique pins
	// datePolicyPins map[string]map[string]map[int]struct{}
	// flipping the nesting avoids iterating all P policies for every range query;
	// only dates actually in the range are touched → O(K + V) instead of O(P·K + V)
	sortedDates []string // sorted unique dates for binary search
}

func NewViolationTracker(records []ViolationRecord) *ViolationTracker {
	vt := &ViolationTracker{
		policyPins:     make(map[string]map[int]struct{}),
		datePins:       make(map[string]map[int]struct{}),
		policyDatePins: make(map[string]map[string]map[int]struct{}),
	}
	dateSet := map[string]struct{}{}
	for _, r := range records {
		vt.add(r.PinID, r.Policy, r.Date)
		dateSet[r.Date] = struct{}{}
	}
	for d := range dateSet {
		vt.sortedDates = append(vt.sortedDates, d)
	}
	sort.Strings(vt.sortedDates)
	return vt
}

func (vt *ViolationTracker) add(pinID int, policy, date string) {
	if vt.policyPins[policy] == nil {
		vt.policyPins[policy] = map[int]struct{}{}
	}
	vt.policyPins[policy][pinID] = struct{}{}

	if vt.datePins[date] == nil {
		vt.datePins[date] = map[int]struct{}{}
	}
	vt.datePins[date][pinID] = struct{}{}

	if vt.policyDatePins[policy] == nil {
		vt.policyDatePins[policy] = map[string]map[int]struct{}{}
	}
	if vt.policyDatePins[policy][date] == nil {
		vt.policyDatePins[policy][date] = map[int]struct{}{}
	}
	vt.policyDatePins[policy][date][pinID] = struct{}{}
}

// CountByPolicy returns the number of unique pins that violated the given policy.
// O(1)
func (vt *ViolationTracker) CountByPolicy(policy string) int {
	return len(vt.policyPins[policy])
}

// CountUniqueInRange returns the number of unique pins that violated any policy
// within [startDate, endDate] inclusive.
// O(log D + K·N) where D = unique dates, K = dates in range, N = pins per date
func (vt *ViolationTracker) CountUniqueInRange(startDate, endDate string) int {
	lo, hi := vt.dateRange(startDate, endDate)
	unique := map[int]struct{}{}
	for _, date := range vt.sortedDates[lo:hi] {
		for pin := range vt.datePins[date] {
			unique[pin] = struct{}{}
		}
	}
	return len(unique)
}

// CountPerPolicyInRange returns a map of policy → unique pin count for violations
// within [startDate, endDate] inclusive.
// O(P·K + V) where P = unique policies, K = dates in range, V = violations in range.
// Outer loop visits all P policies even if they have no violations in the range.
//
// Alternative with datePolicyPins (date → policy → pins): O(K + V) — only touches
// dates in the range, skipping policies with no activity entirely:
//
//	perPolicy := map[string]map[int]struct{}{}
//	for _, date := range vt.sortedDates[lo:hi] {
//	    for policy, pins := range vt.datePolicyPins[date] {
//	        if perPolicy[policy] == nil {
//	            perPolicy[policy] = map[int]struct{}{}
//	        }
//	        for pin := range pins {
//	            perPolicy[policy][pin] = struct{}{}
//	        }
//	    }
//	}
//	for policy, pins := range perPolicy {
//	    result[policy] = len(pins)
//	}
func (vt *ViolationTracker) CountPerPolicyInRange(startDate, endDate string) map[string]int {
	lo, hi := vt.dateRange(startDate, endDate)
	dates := vt.sortedDates[lo:hi]
	result := map[string]int{}
	for policy, datePins := range vt.policyDatePins {
		unique := map[int]struct{}{}
		for _, date := range dates {
			for pin := range datePins[date] {
				unique[pin] = struct{}{}
			}
		}
		if len(unique) > 0 {
			result[policy] = len(unique)
		}
	}
	return result
}

// dateRange returns [lo, hi) indices into sortedDates for [startDate, endDate] inclusive.
func (vt *ViolationTracker) dateRange(startDate, endDate string) (lo, hi int) {
	lo = sort.SearchStrings(vt.sortedDates, startDate)
	hi = sort.SearchStrings(vt.sortedDates, endDate)
	if hi < len(vt.sortedDates) && vt.sortedDates[hi] == endDate {
		hi++
	}
	return lo, hi
}
