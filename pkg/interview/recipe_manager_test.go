package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── AddRecipe / GetRecipe ──────────────────────────────────────────────────────

func Test_RecipeManager_AddRecipe(t *testing.T) {
	testCases := []struct {
		name        string
		recipeName  string
		ingredients []string
		steps       []string
		wantNil     bool
		wantID      string
	}{
		{
			name:        "first_recipe_gets_id_recipe1",
			recipeName:  "Coffee",
			ingredients: []string{"water", "beans"},
			steps:       []string{"brew"},
			wantID:      "recipe1",
		},
		{
			name:        "second_recipe_gets_incremented_id",
			recipeName:  "Tea",
			ingredients: []string{"water", "leaves"},
			steps:       []string{"steep"},
			wantID:      "recipe2",
		},
	}

	d := NewDigitalRecipeManager()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := d.AddRecipe(tc.recipeName, tc.ingredients, tc.steps)
			if tc.wantNil {
				assert.Nil(t, id)
			} else {
				assert.NotNil(t, id)
				assert.Equal(t, tc.wantID, *id)
			}
		})
	}

	t.Run("duplicate_name_case_insensitive_returns_nil", func(t *testing.T) {
		d2 := NewDigitalRecipeManager()
		d2.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
		id := d2.AddRecipe("coffee", []string{"water"}, []string{"brew"})
		assert.Nil(t, id)
	})
}

func Test_RecipeManager_GetRecipe(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		setup    func(*DigitalRecipeManager)
		wantLen  int
		wantName string
	}{
		{
			name:    "get_nonexistent_recipe_returns_empty",
			id:      "recipe99",
			setup:   func(_ *DigitalRecipeManager) {},
			wantLen: 0,
		},
		{
			name: "get_existing_recipe_returns_name_ingredients_steps",
			id:   "recipe1",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water", "beans"}, []string{"boil", "brew"})
			},
			wantLen:  3,
			wantName: "Coffee",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			result := d.GetRecipe(tc.id)
			assert.Len(t, result, tc.wantLen)
			if tc.wantLen > 0 {
				assert.Equal(t, tc.wantName, result[0])
			}
		})
	}
}

// ── UpdateRecipe ──────────────────────────────────────────────────────────────

func Test_RecipeManager_UpdateRecipe(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(*DigitalRecipeManager)
		recipeId    string
		newName     string
		ingredients []string
		steps       []string
		wantOk      bool
	}{
		{
			name: "update_existing_recipe_succeeds",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
			},
			recipeId:    "recipe1",
			newName:     "Espresso",
			ingredients: []string{"water", "beans"},
			steps:       []string{"press"},
			wantOk:      true,
		},
		{
			name:        "update_nonexistent_recipe_fails",
			setup:       func(_ *DigitalRecipeManager) {},
			recipeId:    "recipe99",
			newName:     "X",
			ingredients: nil,
			steps:       nil,
			wantOk:      false,
		},
		{
			name: "update_name_conflict_with_other_recipe_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.AddRecipe("Tea", []string{"water"}, []string{"steep"})
			},
			recipeId:    "recipe1",
			newName:     "tea", // case-insensitive conflict with recipe2
			ingredients: []string{"water"},
			steps:       []string{"x"},
			wantOk:      false,
		},
		{
			name: "update_to_same_name_succeeds",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
			},
			recipeId:    "recipe1",
			newName:     "Coffee", // same name, same recipe ID
			ingredients: []string{"water", "beans"},
			steps:       []string{"brew"},
			wantOk:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			ok := d.UpdateRecipe(tc.recipeId, tc.newName, tc.ingredients, tc.steps)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}

// ── DeleteRecipe ──────────────────────────────────────────────────────────────

func Test_RecipeManager_DeleteRecipe(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(*DigitalRecipeManager)
		recipeId string
		wantOk   bool
	}{
		{
			name: "delete_existing_recipe_succeeds",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
			},
			recipeId: "recipe1", wantOk: true,
		},
		{
			name:     "delete_nonexistent_recipe_fails",
			setup:    func(_ *DigitalRecipeManager) {},
			recipeId: "recipe99", wantOk: false,
		},
		{
			name: "deleted_name_can_be_reused",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.DeleteRecipe("recipe1")
			},
			recipeId: "recipe1", wantOk: false, // already deleted
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			ok := d.DeleteRecipe(tc.recipeId)
			assert.Equal(t, tc.wantOk, ok)
		})
	}

	t.Run("deleted_name_allows_re_add", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
		d.DeleteRecipe("recipe1")
		id := d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
		assert.NotNil(t, id)
	})
}

// ── SearchRecipesByIngredient ─────────────────────────────────────────────────

func Test_RecipeManager_SearchRecipesByIngredient(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*DigitalRecipeManager)
		ingredient string
		want       []string
	}{
		{
			name:       "no_matches_returns_empty",
			setup:      func(_ *DigitalRecipeManager) {},
			ingredient: "sugar", want: []string{},
		},
		{
			name: "case_insensitive_ingredient_match",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"Water", "beans"}, []string{"brew"})
			},
			ingredient: "water",
			want:       []string{"recipe1"},
		},
		{
			name: "results_sorted_by_ingredient_count_then_id",
			setup: func(d *DigitalRecipeManager) {
				// recipe1: 3 ingredients
				d.AddRecipe("A", []string{"egg", "milk", "sugar"}, []string{"mix"})
				// recipe2: 1 ingredient
				d.AddRecipe("B", []string{"egg"}, []string{"fry"})
				// recipe3: 2 ingredients
				d.AddRecipe("C", []string{"egg", "butter"}, []string{"bake"})
			},
			ingredient: "egg",
			want:       []string{"recipe2", "recipe3", "recipe1"},
		},
		{
			name: "tiebreak_by_recipe_id_when_same_ingredient_count",
			setup: func(d *DigitalRecipeManager) {
				// recipe1 and recipe2 both have 2 ingredients containing "egg"
				d.AddRecipe("A", []string{"egg", "milk"}, []string{"mix"})
				d.AddRecipe("B", []string{"egg", "butter"}, []string{"fry"})
			},
			ingredient: "egg",
			want:       []string{"recipe1", "recipe2"}, // recipe1 < recipe2 lexicographically
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			got := d.SearchRecipesByIngredient(tc.ingredient)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ── ListRecipes ───────────────────────────────────────────────────────────────

func Test_RecipeManager_ListRecipes(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func(*DigitalRecipeManager)
		sortBy string
		want   []string
	}{
		{
			name: "sort_by_name_alphabetical",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Zebra cake", []string{"flour"}, []string{"bake"})
				d.AddRecipe("Apple pie", []string{"apple", "flour"}, []string{"bake"})
				d.AddRecipe("Muffin", []string{"flour", "egg"}, []string{"mix"})
			},
			sortBy: "name",
			want:   []string{"recipe2", "recipe3", "recipe1"},
		},
		{
			name: "sort_by_ingredient_count_asc",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("A", []string{"x", "y", "z"}, []string{"step"})
				d.AddRecipe("B", []string{"x"}, []string{"step"})
				d.AddRecipe("C", []string{"x", "y"}, []string{"step"})
			},
			sortBy: "ingredient_count",
			want:   []string{"recipe2", "recipe3", "recipe1"},
		},
		{
			name: "empty_manager_returns_empty",
			setup: func(_ *DigitalRecipeManager) {},
			sortBy: "name",
			want:   []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			got := d.ListRecipes(tc.sortBy)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ── AddUser / EditRecipe ──────────────────────────────────────────────────────

func Test_RecipeManager_AddUser(t *testing.T) {
	t.Run("new_user_succeeds", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		assert.True(t, d.AddUser("u1"))
	})

	t.Run("duplicate_user_fails", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		d.AddUser("u1")
		assert.False(t, d.AddUser("u1"))
	})
}

func Test_RecipeManager_EditRecipe(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(*DigitalRecipeManager)
		userId      string
		recipeId    string
		newName     string
		ingredients []string
		steps       []string
		wantOk      bool
	}{
		{
			name: "valid_user_edits_recipe",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
			},
			userId: "u1", recipeId: "recipe1",
			newName: "Latte", ingredients: []string{"water", "milk"}, steps: []string{"steam"},
			wantOk: true,
		},
		{
			name: "unknown_user_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
			},
			userId: "nobody", recipeId: "recipe1",
			newName: "X", ingredients: nil, steps: nil,
			wantOk: false,
		},
		{
			name: "nonexistent_recipe_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
			},
			userId: "u1", recipeId: "recipe99",
			newName: "X", ingredients: nil, steps: nil,
			wantOk: false,
		},
		{
			name: "name_conflict_with_other_recipe_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.AddRecipe("Tea", []string{"water"}, []string{"steep"})
			},
			userId: "u1", recipeId: "recipe1",
			newName: "Tea", ingredients: []string{"water"}, steps: []string{"x"},
			wantOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			ok := d.EditRecipe(tc.userId, tc.recipeId, tc.newName, tc.ingredients, tc.steps)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}

// ── GetRecipeHistory / RollbackRecipe ─────────────────────────────────────────

func Test_RecipeManager_GetRecipeHistory(t *testing.T) {
	t.Run("nonexistent_recipe_returns_nil", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		assert.Nil(t, d.GetRecipeHistory("recipe99"))
	})

	t.Run("history_empty_until_snapshot_called", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
		// No snapshot recorded by AddRecipe — history is empty
		h := d.GetRecipeHistory("recipe1")
		assert.Empty(t, h)
	})

	t.Run("snapshot_records_version", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
		d.recipes["recipe1"].snapshot()

		h := d.GetRecipeHistory("recipe1")
		assert.Len(t, h, 1)
		assert.Equal(t, 1, h[0]["version"])
		assert.Equal(t, "Coffee", h[0]["name"])
	})
}

func Test_RecipeManager_RollbackRecipe(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(*DigitalRecipeManager)
		userId   string
		recipeId string
		version  int
		wantOk   bool
	}{
		{
			name: "rollback_to_version_1_succeeds",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.recipes["recipe1"].snapshot() // version 1
				d.EditRecipe("u1", "recipe1", "Espresso", []string{"water", "beans"}, []string{"press"})
			},
			userId: "u1", recipeId: "recipe1", version: 1,
			wantOk: true,
		},
		{
			name: "rollback_invalid_version_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.recipes["recipe1"].snapshot() // only version 1 exists
			},
			userId: "u1", recipeId: "recipe1", version: 5,
			wantOk: false,
		},
		{
			name: "rollback_version_zero_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.recipes["recipe1"].snapshot()
			},
			userId: "u1", recipeId: "recipe1", version: 0,
			wantOk: false,
		},
		{
			name: "rollback_unknown_user_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
				d.recipes["recipe1"].snapshot()
			},
			userId: "nobody", recipeId: "recipe1", version: 1,
			wantOk: false,
		},
		{
			name: "rollback_nonexistent_recipe_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
			},
			userId: "u1", recipeId: "recipe99", version: 1,
			wantOk: false,
		},
		{
			name: "rollback_name_conflict_with_another_recipe_fails",
			setup: func(d *DigitalRecipeManager) {
				d.AddUser("u1")
				// recipe1 has v1="Alpha", then renamed to "Beta"
				d.AddRecipe("Alpha", []string{"x"}, []string{"step"})
				d.recipes["recipe1"].snapshot() // v1: Alpha
				d.EditRecipe("u1", "recipe1", "Beta", []string{"x"}, []string{"step"})
				// recipe2 now owns "Alpha"
				d.AddRecipe("Alpha", []string{"y"}, []string{"step"})
			},
			// Rolling back recipe1 to v1 would restore name "Alpha", but recipe2 owns it
			userId: "u1", recipeId: "recipe1", version: 1,
			wantOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDigitalRecipeManager()
			tc.setup(d)
			ok := d.RollbackRecipe(tc.userId, tc.recipeId, tc.version)
			assert.Equal(t, tc.wantOk, ok)
		})
	}

	t.Run("rollback_restores_old_name_and_creates_new_snapshot", func(t *testing.T) {
		d := NewDigitalRecipeManager()
		d.AddUser("u1")
		d.AddRecipe("Coffee", []string{"water"}, []string{"brew"})
		d.recipes["recipe1"].snapshot() // version 1: Coffee

		d.EditRecipe("u1", "recipe1", "Espresso", []string{"water", "beans"}, []string{"press"})
		d.recipes["recipe1"].snapshot() // version 2: Espresso

		ok := d.RollbackRecipe("u1", "recipe1", 1) // roll back to Coffee
		assert.True(t, ok)

		info := d.GetRecipe("recipe1")
		assert.Equal(t, "Coffee", info[0])

		// A new snapshot should have been added
		h := d.GetRecipeHistory("recipe1")
		assert.Equal(t, 3, len(h))
	})
}
