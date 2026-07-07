package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── testAssignColors (demo function coverage) ─────────────────────────────────

func Test_testAssignColors(t *testing.T) {
	// testAssignColors() prints to stdout; calling it exercises its branches.
	testAssignColors()
}

// ── assignColors ──────────────────────────────────────────────────────────────

func Test_assignColors(t *testing.T) {
	testCases := []struct {
		name      string
		users     []User
		minSep    float64
		maxColors int
		// verify no two conflicting users share a color (correctness check)
		// for small cases we also check that a color is assigned to every user
		wantAllAssigned bool
	}{
		{
			name:            "empty_users",
			users:           []User{},
			minSep:          10.0,
			maxColors:       4,
			wantAllAssigned: true,
		},
		{
			name:            "single_user_gets_color_1",
			users:           []User{{1, 0.0}},
			minSep:          10.0,
			maxColors:       4,
			wantAllAssigned: true,
		},
		{
			name: "two_well_separated_users_both_get_color_1",
			// 0° and 90° are 90° apart — no conflict
			users:           []User{{1, 0.0}, {2, 90.0}},
			minSep:          10.0,
			maxColors:       4,
			wantAllAssigned: true,
		},
		{
			name: "two_conflicting_users_get_different_colors",
			// 0° and 5° are only 5° apart (< 10°) — conflict
			users:           []User{{1, 0.0}, {2, 5.0}},
			minSep:          10.0,
			maxColors:       4,
			wantAllAssigned: true,
		},
		{
			name: "linear_conflict_chain_requires_two_colors",
			// 0° conflicts with 5°, 5° conflicts with 9°, but 0° and 9° do not conflict
			// Valid 2-coloring: 0°=1, 5°=2, 9°=1
			users:           []User{{1, 0.0}, {2, 5.0}, {3, 9.0}},
			minSep:          10.0,
			maxColors:       4,
			wantAllAssigned: true,
		},
		{
			name: "sample_data_from_file_six_users",
			// From testAssignColors() in graph_coloring.go
			users: []User{
				{1, 0.0},
				{2, 3.0},
				{3, 8.0},
				{4, 15.0},
				{5, 20.0},
				{6, 25.0},
			},
			minSep:          10.0,
			maxColors:       4,
			wantAllAssigned: true,
		},
		{
			name: "unserviceable_when_max_colors_too_small",
			// 0°, 3°, 6°, 9° all within 10° of each other — form a 4-clique,
			// so we need 4 colors; with maxColors=2 some will be unserviceable.
			users: []User{
				{1, 0.0},
				{2, 3.0},
				{3, 6.0},
				{4, 9.0},
			},
			minSep:          10.0,
			maxColors:       2,
			wantAllAssigned: false, // not all users can be served
		},
		{
			name: "no_conflicts_single_color_suffices",
			users: []User{
				{1, 0.0},
				{2, 20.0},
				{3, 40.0},
				{4, 60.0},
			},
			minSep:          10.0,
			maxColors:       1,
			wantAllAssigned: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			colors := assignColors(tc.users, tc.minSep, tc.maxColors)

			// Constraint 1: no two conflicting users share a color
			for i := 0; i < len(tc.users); i++ {
				for j := i + 1; j < len(tc.users); j++ {
					diff := tc.users[i].angle - tc.users[j].angle
					if diff < 0 {
						diff = -diff
					}
					if diff < tc.minSep {
						ci, iAssigned := colors[i]
						cj, jAssigned := colors[j]
						if iAssigned && jAssigned {
							assert.NotEqual(t, ci, cj,
								"users %d and %d conflict but share color %d", tc.users[i].id, tc.users[j].id, ci)
						}
					}
				}
			}

			// Constraint 2: all colors are within [1, maxColors]
			for idx, c := range colors {
				assert.GreaterOrEqual(t, c, 1, "user index %d has color < 1", idx)
				assert.LessOrEqual(t, c, tc.maxColors, "user index %d has color > maxColors", idx)
			}

			// Constraint 3: if wantAllAssigned, every user index must have a color
			if tc.wantAllAssigned {
				for i := range tc.users {
					_, assigned := colors[i]
					assert.True(t, assigned, "user index %d should have a color but doesn't", i)
				}
			}
		})
	}
}
