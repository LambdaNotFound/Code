package utils

import "sync"

type Color int

const (
	Red Color = iota // iota represents the zero-indexed integer ordinal number of the current
	Blue
	White
	Green
	Yellow
)

func (c Color) String() string {
	return []string{"R", "B", "W", "G", "Y"}[c]
}

type Card struct {
	color Color
	cost  map[Color]int
}

func NewCard(color Color, cost map[Color]int) Card {
	return Card{color: color, cost: cost}
}

type Player struct {
	gems  map[Color]int
	cards []Card
	// Task 3
	discount map[Color]int
	// Task 4
	mutex sync.Mutex
}

func NewPlayer(gems map[Color]int) *Player {
	return &Player{gems: gems, cards: []Card{}, discount: make(map[Color]int)}
}

// Task 1
func (p *Player) canPurchase(card Card) bool {
	for color, cost := range card.cost {
		actualCost := max(0, cost-p.discount[color])
		if p.gems[color]-actualCost < 0 {
			return false
		}
	}
	return true
}

// Task 2
func (p *Player) Purchase(card Card) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.canPurchase(card) {
		return false
	}

	for color, cost := range card.cost {
		actualCost := max(0, cost-p.discount[color])
		p.gems[color] -= actualCost
	}
	p.cards = append(p.cards, card)
	p.discount[card.color]++ // Task 3, O(1) incremental update
	return true
}
