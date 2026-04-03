package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HitCounter(t *testing.T) {
	t.Run("no_hits_returns_zero", func(t *testing.T) {
		hc := ConstructorHitCounter()
		assert.Equal(t, 0, hc.GetHits(100))
	})

	t.Run("hit_within_window_counted", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(1)
		assert.Equal(t, 1, hc.GetHits(1))
	})

	t.Run("hit_outside_300s_window_excluded", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(1)
		// timestamp=302: window is (2, 302], so timestamp=1 falls out
		assert.Equal(t, 0, hc.GetHits(302))
	})

	t.Run("hit_at_boundary_excluded", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(1)
		// window is (301-300, 301] = (1, 301], so exactly timestamp=1 is excluded
		assert.Equal(t, 0, hc.GetHits(301))
	})

	t.Run("hit_just_inside_boundary_included", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(2)
		// window at t=301 is (1, 301], so timestamp=2 is included
		assert.Equal(t, 1, hc.GetHits(301))
	})

	t.Run("multiple_hits_counted", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(1)
		hc.Hit(2)
		hc.Hit(3)
		assert.Equal(t, 3, hc.GetHits(4))
	})

	t.Run("only_in_window_hits_counted", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(1)   // outside window at t=302
		hc.Hit(100) // inside window at t=302
		hc.Hit(200) // inside window at t=302
		assert.Equal(t, 2, hc.GetHits(302))
	})

	t.Run("leetcode_example", func(t *testing.T) {
		hc := ConstructorHitCounter()
		hc.Hit(1)
		hc.Hit(2)
		hc.Hit(3)
		assert.Equal(t, 3, hc.GetHits(4))
		hc.Hit(300)
		assert.Equal(t, 4, hc.GetHits(300))
		assert.Equal(t, 3, hc.GetHits(301)) // hit(1) falls out
	})
}
