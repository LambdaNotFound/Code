package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolveTopLevelIssuers(t *testing.T) {
	t.Run("simple_hierarchy", func(t *testing.T) {
		hierarchy := map[string][]string{
			"Root": {"ChildA", "ChildB"},
			"ChildA": {"GrandChildA"},
		}
		issuers, err := resolveTopLevelIssuers(hierarchy, []string{"GrandChildA", "ChildB"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Root", "Root"}, issuers)
	})

	t.Run("top_level_company_resolves_to_itself", func(t *testing.T) {
		hierarchy := map[string][]string{
			"Parent": {"Child"},
		}
		issuers, err := resolveTopLevelIssuers(hierarchy, []string{"Parent"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"Parent"}, issuers)
	})

	t.Run("multiple_levels_deep", func(t *testing.T) {
		hierarchy := map[string][]string{
			"A": {"B"},
			"B": {"C"},
			"C": {"D"},
		}
		issuers, err := resolveTopLevelIssuers(hierarchy, []string{"D"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"A"}, issuers)
	})

	t.Run("cycle_returns_error", func(t *testing.T) {
		hierarchy := map[string][]string{
			"A": {"B"},
			"B": {"A"},
		}
		_, err := resolveTopLevelIssuers(hierarchy, []string{"A"})
		assert.Error(t, err)
	})

	t.Run("empty_loan_list", func(t *testing.T) {
		issuers, err := resolveTopLevelIssuers(map[string][]string{}, []string{})
		assert.NoError(t, err)
		assert.Nil(t, issuers)
	})
}
