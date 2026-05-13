package oodesign

import (
	"reflect"
	"testing"
)

func TestFilterAds(t *testing.T) {
	ads := []Ad{
		{ID: "A", TargetLocations: map[string]bool{"US": true, "CA": true}, Age: 21},
		{ID: "B", TargetLocations: map[string]bool{"US": true}, Age: 13},
		{ID: "C", TargetLocations: map[string]bool{"CA": true}, Age: 25},
		{ID: "D", TargetLocations: map[string]bool{"TX": true, "WA": true}, Age: 30},
		{ID: "E", TargetLocations: map[string]bool{"US": true}, Age: 21},
	}

	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name string
		rule Rule
		want []string
	}{
		{
			name: "AND: minAge and location — only A matches both",
			rule: Rule{
				TargetLocations: map[string]bool{"US": true},
				MinAge:          intPtr(21),
				Operator:        "AND",
			},
			want: []string{"A", "E"},
		},
		{
			name: "AND: minAge, maxAge, and location",
			rule: Rule{
				TargetLocations: map[string]bool{"US": true},
				MinAge:          intPtr(21),
				MaxAge:          intPtr(22),
				Operator:        "AND",
			},
			want: []string{"A", "E"},
		},
		{
			name: "AND: age range excludes all",
			rule: Rule{
				MinAge:   intPtr(50),
				MaxAge:   intPtr(60),
				Operator: "AND",
			},
			want: nil,
		},
		{
			name: "OR: location or minAge — broad match",
			rule: Rule{
				TargetLocations: map[string]bool{"TX": true},
				MinAge:          intPtr(25),
				Operator:        "OR",
			},
			want: []string{"C", "D"}, // E: age=21 < 25, no TX → no match
		},
		{
			name: "OR: either location matches",
			rule: Rule{
				TargetLocations: map[string]bool{"WA": true},
				MinAge:          intPtr(21),
				Operator:        "OR",
			},
			want: []string{"A", "C", "D", "E"},
		},
		{
			name: "location only (AND with single condition)",
			rule: Rule{
				TargetLocations: map[string]bool{"CA": true},
				Operator:        "AND",
			},
			want: []string{"A", "C"},
		},
		{
			name: "minAge only",
			rule: Rule{
				MinAge:   intPtr(21),
				Operator: "AND",
			},
			want: []string{"A", "C", "D", "E"},
		},
		{
			name: "maxAge only",
			rule: Rule{
				MaxAge:   intPtr(13),
				Operator: "AND",
			},
			want: []string{"B"},
		},
		{
			name: "rule location lowercase normalized to match ad uppercase",
			rule: Rule{
				TargetLocations: map[string]bool{"us": true},
				Operator:        "AND",
			},
			want: []string{"A", "B", "E"},
		},
		{
			name: "no conditions — all ads match",
			rule: Rule{
				Operator: "AND",
			},
			want: []string{"A", "B", "C", "D", "E"},
		},
		{
			name: "location with no overlap — no match",
			rule: Rule{
				TargetLocations: map[string]bool{"NY": true},
				Operator:        "AND",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterAds(ads, tt.rule)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nrule: %+v\ngot:  %v\nwant: %v", tt.rule, got, tt.want)
			}
		})
	}
}
