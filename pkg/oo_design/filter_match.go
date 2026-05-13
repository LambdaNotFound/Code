package oodesign

import (
	"strconv"
	"strings"
)

/*
Two patterns working together:

Strategy — Predicate is the interface; NumPredicate, SetPredicate, StrPredicate are interchangeable implementations. FilterAds only calls pred.Match(r) and never knows which concrete type it's talking to.

Composite — AndPredicate and OrPredicate both hold []Predicate and themselves implement Predicate. This means they can be nested arbitrarily:

AndPredicate
├── OrPredicate          ← locations (TX or WA)
│   ├── SetPredicate{TX}
│   └── SetPredicate{WA}
└── NumPredicate{age >= 21}
The Composite pattern is what makes buildPredicate powerful — the location inner OrPredicate slots into the outer AndPredicate seamlessly, because everything satisfies the same Predicate interface. FilterAds doesn't need to know the tree exists; it just calls Match.

Factory — buildPredicate(rule Rule) is also present, translating the Rule data into a predicate tree. It's the seam between the domain model and the filter engine.
*/

// interface: only behaviors
type Predicate interface {
	Match(r Record) bool
}

type Record struct {
	Scalars map[string]string
	Sets    map[string]map[string]bool
}

// struct: only properties
type NumPredicate struct {
	col string
	val int
	op  string
}

func (np NumPredicate) Match(r Record) bool {
	number, err := strconv.Atoi(r.Scalars[np.col])
	if err != nil {
		return false
	}

	switch np.op {
	case ">":
		return number > np.val
	case ">=":
		return number >= np.val
	case "=":
		return number == np.val
	case "<=":
		return number <= np.val
	case "<":
		return number < np.val
	}

	return false
}

type SetPredicate struct {
	col string
	val string
}

func (p SetPredicate) Match(r Record) bool {
	return r.Sets[p.col][strings.ToUpper(p.val)]
}

type AndPredicate struct {
	predicates []Predicate
}

func (ap *AndPredicate) Match(r Record) bool {
	for _, p := range ap.predicates {
		if !p.Match(r) {
			return false
		}
	}
	return true
}

type OrPredicate struct {
	predicates []Predicate
}

func (op *OrPredicate) Match(r Record) bool {
	for _, p := range op.predicates {
		if p.Match(r) {
			return true
		}
	}
	return false
}

/**
 * Single-Rule Ad Filter
 */
type Ad struct {
	ID              string
	TargetLocations map[string]bool
	Age             int
}

type Rule struct {
	TargetLocations map[string]bool
	MinAge          *int // pointer so nil = "not set"
	MaxAge          *int
	Operator        string // "AND" | "OR"
}

func buildPredicate(rule Rule) Predicate {
	var preds []Predicate
	if rule.TargetLocations != nil {
		var locPreds []Predicate
		for loc := range rule.TargetLocations {
			locPreds = append(locPreds, SetPredicate{"targetLocations", loc})
		}
		preds = append(preds, &OrPredicate{locPreds})
	}
	if rule.MinAge != nil {
		preds = append(preds, NumPredicate{"age", *rule.MinAge, ">="})
	}
	if rule.MaxAge != nil {
		preds = append(preds, NumPredicate{"age", *rule.MaxAge, "<="})
	}
	if rule.Operator == "OR" {
		return &OrPredicate{preds}
	}
	return &AndPredicate{preds}
}

func adToRecord(ad Ad) Record {
	return Record{
		Scalars: map[string]string{
			"age": strconv.Itoa(ad.Age),
		},
		Sets: map[string]map[string]bool{
			"targetLocations": ad.TargetLocations,
		},
	}
}

func FilterAds(ads []Ad, rule Rule) []string {
	pred := buildPredicate(rule)
	var result []string
	for _, ad := range ads {
		if pred.Match(adToRecord(ad)) {
			result = append(result, ad.ID)
		}
	}
	return result
}
