package interview

import (
	"sort"
	"strconv"
	"strings"
)

/**
 * 1. recipe name is case insensitive
 * 2. recipe[0-9] is the recipe id
 */

type DigitalRecipeManager struct {
	recipes   map[string]*Recipe // recipeId -> Recipe
	nameIndex map[string]string  // lowercase name -> recipeId
	counter   int
	// level 3
	users map[string]bool // userId -> exists
}

type Recipe struct {
	id          string
	name        string
	ingredients []string
	steps       []string
	// level 4
	history []RecipeVersion // NEW
}
type RecipeVersion struct {
	version     int
	name        string
	ingredients []string
	steps       []string
}

// Helper — call this at the end of any mutation
// r.snapshot() // only addition needed
func (r *Recipe) snapshot() {
	r.history = append(r.history, RecipeVersion{
		version:     len(r.history) + 1,
		name:        r.name,
		ingredients: append([]string{}, r.ingredients...),
		steps:       append([]string{}, r.steps...),
	})
}

func NewDigitalRecipeManager() *DigitalRecipeManager {
	return &DigitalRecipeManager{
		recipes:   make(map[string]*Recipe),
		nameIndex: make(map[string]string),
		users:     make(map[string]bool),
	}
}

// level 1: CRUD
func (d *DigitalRecipeManager) AddRecipe(name string, ingredients []string, steps []string) *string {
	lower := strings.ToLower(name)
	if _, exists := d.nameIndex[lower]; exists {
		return nil
	}

	d.counter++
	id := "recipe" + strconv.Itoa(d.counter)

	r := &Recipe{id: id, name: name, ingredients: ingredients, steps: steps}
	d.recipes[id] = r
	d.nameIndex[lower] = id
	return &id
}

func (d *DigitalRecipeManager) GetRecipe(recipeId string) []string {
	r, ok := d.recipes[recipeId]
	if !ok {
		return []string{}
	}
	return []string{
		r.name,
		strings.Join(r.ingredients, ", "),
		strings.Join(r.steps, ", "),
	}
}

// recipeX, coffee => recipeX, tea
func (d *DigitalRecipeManager) UpdateRecipe(recipeId string, name string, ingredients []string, steps []string) bool {
	r, ok := d.recipes[recipeId]
	if !ok {
		return false
	}
	lower := strings.ToLower(name)
	if existing, exists := d.nameIndex[lower]; exists && existing != recipeId {
		return false
	}

	// Remove old name index
	delete(d.nameIndex, strings.ToLower(r.name))
	d.nameIndex[lower] = recipeId

	r.name = name
	r.ingredients = ingredients
	r.steps = steps
	return true
}

func (d *DigitalRecipeManager) DeleteRecipe(recipeId string) bool {
	r, ok := d.recipes[recipeId]
	if !ok {
		return false
	}
	delete(d.nameIndex, strings.ToLower(r.name))
	delete(d.recipes, recipeId)
	return true
}

// level 2
func (d *DigitalRecipeManager) SearchRecipesByIngredient(ingredient string) []string {
	lowerIngredient := strings.ToLower(ingredient)

	var matched []*Recipe
	for _, r := range d.recipes {
		for _, ing := range r.ingredients {
			if strings.ToLower(ing) == lowerIngredient {
				matched = append(matched, r)
				break
			}
		}
	}

	// Sort by ingredient count asc, then recipe ID asc (sort string)
	sort.Slice(matched, func(i, j int) bool {
		if len(matched[i].ingredients) != len(matched[j].ingredients) {
			return len(matched[i].ingredients) < len(matched[j].ingredients)
		}
		return matched[i].id < matched[j].id
	})

	result := make([]string, len(matched))
	for i, r := range matched {
		result[i] = r.id
	}
	return result
}

func (d *DigitalRecipeManager) ListRecipes(sortBy string) []string {
	all := make([]*Recipe, 0, len(d.recipes))
	for _, r := range d.recipes {
		all = append(all, r)
	}

	if sortBy == "ingredient_count" {
		sort.Slice(all, func(i, j int) bool {
			if len(all[i].ingredients) != len(all[j].ingredients) {
				return len(all[i].ingredients) < len(all[j].ingredients)
			}
			return all[i].id < all[j].id
		})
	} else {
		// Default: sort by name lexicographically, tie-break by recipe ID
		sort.Slice(all, func(i, j int) bool {
			ni := strings.ToLower(all[i].name)
			nj := strings.ToLower(all[j].name)
			if ni != nj {
				return ni < nj
			}
			return all[i].id < all[j].id
		})
	}

	result := make([]string, len(all))
	for i, r := range all {
		result[i] = r.id
	}
	return result
}

// level 3
func (d *DigitalRecipeManager) AddUser(userId string) bool {
	if d.users[userId] {
		return false
	}
	d.users[userId] = true
	return true
}

func (d *DigitalRecipeManager) EditRecipe(userId string, recipeId string, newName string, newIngredients []string, newSteps []string) bool {
	// User must exist
	if !d.users[userId] {
		return false
	}
	// Recipe must exist
	r, ok := d.recipes[recipeId]
	if !ok {
		return false
	}
	// Name conflict check (case-insensitive, allow same recipe)
	lower := strings.ToLower(newName)
	if existing, exists := d.nameIndex[lower]; exists && existing != recipeId {
		return false
	}
	delete(d.nameIndex, strings.ToLower(r.name))

	r.name = newName
	r.ingredients = newIngredients
	r.steps = newSteps
	d.nameIndex[lower] = recipeId
	return true
}

// level 4
func (d *DigitalRecipeManager) GetRecipeHistory(recipeId string) []map[string]interface{} {
	r, ok := d.recipes[recipeId]
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, len(r.history))
	for i, v := range r.history {
		result[i] = map[string]interface{}{
			"version":     v.version,
			"name":        v.name,
			"ingredients": v.ingredients,
			"steps":       v.steps,
		}
	}
	return result
}

func (d *DigitalRecipeManager) RollbackRecipe(userId string, recipeId string, version int) bool {
	if !d.users[userId] {
		return false
	}
	r, ok := d.recipes[recipeId]
	if !ok {
		return false
	}
	if version < 1 || version > len(r.history) {
		return false
	}
	target := r.history[version-1]
	lower := strings.ToLower(target.name)
	if existing, exists := d.nameIndex[lower]; exists && existing != recipeId {
		return false
	}
	// Apply old version onto current fields
	delete(d.nameIndex, strings.ToLower(r.name))
	r.name = target.name
	r.ingredients = append([]string{}, target.ingredients...)
	r.steps = append([]string{}, target.steps...)
	d.nameIndex[lower] = recipeId
	r.snapshot() // records it as a new version
	return true
}
