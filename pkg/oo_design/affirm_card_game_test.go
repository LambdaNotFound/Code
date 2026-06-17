package oodesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGame(t *testing.T) {
	g := NewGame([]string{"Alice", "Bob"}, 52)
	assert.Equal(t, 2, len(g.players))
	assert.Equal(t, "Alice", g.players[0].Name)
	assert.Equal(t, "Bob", g.players[1].Name)
	assert.Equal(t, 52, g.numCards)
}

func Test_deal(t *testing.T) {
	t.Run("even split", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 6)
		g.deal()
		assert.Equal(t, 3, len(g.players[0].Hand))
		assert.Equal(t, 3, len(g.players[1].Hand))
	})

	t.Run("uneven split 5 cards 3 players", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob", "Charlie"}, 5)
		g.deal()
		// round-robin: players 0 and 1 get an extra card
		assert.Equal(t, 2, len(g.players[0].Hand))
		assert.Equal(t, 2, len(g.players[1].Hand))
		assert.Equal(t, 1, len(g.players[2].Hand))
	})

	t.Run("cards are unique integers in valid range", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 10)
		g.deal()
		seen := make(map[int]bool)
		for _, p := range g.players {
			for _, card := range p.Hand {
				assert.False(t, seen[card], "duplicate card %d", card)
				assert.True(t, card >= 1 && card <= 10)
				seen[card] = true
			}
		}
		assert.Equal(t, 10, len(seen))
	})
}

func Test_allHaveCards(t *testing.T) {
	t.Run("all players have cards", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 4)
		g.players[0].Hand = []int{1, 2}
		g.players[1].Hand = []int{3, 4}
		assert.True(t, g.allHaveCards())
	})

	t.Run("one player has no cards", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 2)
		g.players[0].Hand = []int{1}
		g.players[1].Hand = []int{}
		assert.False(t, g.allHaveCards())
	})

	t.Run("all players have no cards", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 0)
		assert.False(t, g.allHaveCards())
	})
}

func Test_playRound2PlayerMode(t *testing.T) {
	t.Run("player 0 wins with higher card", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 4)
		g.players[0].Hand = []int{10, 3}
		g.players[1].Hand = []int{5, 7}
		g.playRound2PlayerMode()
		assert.Equal(t, []int{3}, g.players[0].Hand)
		assert.Equal(t, []int{7}, g.players[1].Hand)
		assert.Equal(t, []int{10, 5}, g.players[0].WinPile)
		assert.Nil(t, g.players[1].WinPile)
	})

	t.Run("player 1 wins with higher card", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 4)
		g.players[0].Hand = []int{3, 10}
		g.players[1].Hand = []int{9, 1}
		g.playRound2PlayerMode()
		assert.Equal(t, []int{10}, g.players[0].Hand)
		assert.Equal(t, []int{1}, g.players[1].Hand)
		assert.Nil(t, g.players[0].WinPile)
		assert.Equal(t, []int{3, 9}, g.players[1].WinPile)
	})
}

func Test_playRoundNPlayerMode(t *testing.T) {
	t.Run("first player wins with highest card", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob", "Charlie"}, 6)
		g.players[0].Hand = []int{10, 1}
		g.players[1].Hand = []int{5, 2}
		g.players[2].Hand = []int{3, 6}
		g.playRoundNPlayerMode()
		assert.Equal(t, []int{1}, g.players[0].Hand)
		assert.Equal(t, []int{2}, g.players[1].Hand)
		assert.Equal(t, []int{6}, g.players[2].Hand)
		assert.Equal(t, []int{10, 5, 3}, g.players[0].WinPile)
		assert.Nil(t, g.players[1].WinPile)
		assert.Nil(t, g.players[2].WinPile)
	})

	t.Run("last player wins with highest card", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob", "Charlie"}, 3)
		g.players[0].Hand = []int{2}
		g.players[1].Hand = []int{5}
		g.players[2].Hand = []int{9}
		g.playRoundNPlayerMode()
		assert.Empty(t, g.players[0].Hand)
		assert.Empty(t, g.players[1].Hand)
		assert.Empty(t, g.players[2].Hand)
		assert.Nil(t, g.players[0].WinPile)
		assert.Nil(t, g.players[1].WinPile)
		assert.Equal(t, []int{2, 5, 9}, g.players[2].WinPile)
	})
}

func Test_winner(t *testing.T) {
	t.Run("clear winner by pile size", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 4)
		g.players[0].WinPile = []int{1, 2, 3}
		g.players[1].WinPile = []int{4}
		assert.Equal(t, "Alice", g.winner())
	})

	t.Run("tiebreak: bob has highest card", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 4)
		g.players[0].WinPile = []int{1, 3}
		g.players[1].WinPile = []int{2, 4}
		assert.Equal(t, "Bob", g.winner())
	})

	t.Run("tiebreak: alice has highest card", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob"}, 4)
		g.players[0].WinPile = []int{5, 1}
		g.players[1].WinPile = []int{2, 4}
		assert.Equal(t, "Alice", g.winner())
	})

	t.Run("three players, middle wins by pile size", func(t *testing.T) {
		g := NewGame([]string{"Alice", "Bob", "Charlie"}, 9)
		g.players[0].WinPile = []int{1, 2}
		g.players[1].WinPile = []int{3, 4, 5}
		g.players[2].WinPile = []int{6}
		assert.Equal(t, "Bob", g.winner())
	})
}

func Test_Play(t *testing.T) {
	t.Run("2-player game returns a player name", func(t *testing.T) {
		names := []string{"Alice", "Bob"}
		g := NewGame(names, 52)
		assert.Contains(t, names, g.Play())
	})

	t.Run("3-player game returns a player name", func(t *testing.T) {
		names := []string{"Joe", "Jill", "Bob"}
		g := NewGame(names, 5)
		assert.Contains(t, names, g.Play())
	})
}
