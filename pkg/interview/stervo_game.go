package interview

import "sync"

/**
 * Game: Stervo
 *
 * Players collect colored gems (B/W/G/R/Y) and buy cards using gems
 * Each card has a cost (gems required) and a color (the gem color it represents)
 */

type Color int

const (
	Blue Color = iota
	White
	Green
	Red
	Yellow
)

func (c Color) String() string {
	return []string{"B", "W", "G", "R", "Y"}[c]
}

type Card struct {
	color Color         // the card's own color (for discounts later)
	cost  map[Color]int // gems required to buy
}

type Player struct {
	gems map[Color]int // current gems
	hand []*Card       // owned cards
	// Part 4: Concurrency
	mu sync.Mutex
}

// Part 1: Data model + can_purchase
/*
// canPurchase: does player have enough gems for every required color?
func canPurchase(p *Player, c *Card) bool {
	for color, required := range c.cost {
		if p.gems[color] < required {
			return false
		}
	}
	return true
}
*/

// Part 2: purchase()
/*
func purchase(p *Player, c *Card) bool {
	if !canPurchase(p, c) {
		return false
	}
	// deduct gems
	for color, required := range c.cost {
		p.gems[color] -= required
	}
	// add to hand
	p.hand = append(p.hand, c)
	return true
}
*/

// Part 3: Discounts
// discounts: how many gems player saves per color, based on owned cards
func discountsOf(p *Player) map[Color]int {
	d := map[Color]int{}
	for _, card := range p.hand {
		d[card.color]++
	}
	return d
}

// effectiveCost: cost after discounts (clamped at 0)
func effectiveCost(c *Card, p *Player) map[Color]int {
	discount := discountsOf(p)
	eff := map[Color]int{}
	for color, required := range c.cost {
		net := required - discount[color]
		if net < 0 {
			net = 0
		}
		eff[color] = net
	}
	return eff
}

// updated canPurchase
func canPurchase(p *Player, c *Card) bool {
	eff := effectiveCost(c, p)
	for color, required := range eff {
		if p.gems[color] < required {
			return false
		}
	}
	return true
}

// updated purchase
func purchase(p *Player, c *Card) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !canPurchase(p, c) {
		return false
	}
	eff := effectiveCost(c, p)
	for color, required := range eff {
		p.gems[color] -= required
	}
	p.hand = append(p.hand, c)
	return true
}

// Part 4: Concurrency
// This is where the question gets interesting.
// Multiple goroutines may call purchase() for the same player concurrently — classic check-then-act race.

// Option 2: Per-player request queue

/*
type Player struct {
    gems    map[Color]int
    hand    []*Card
    reqCh   chan purchaseReq
}

type purchaseReq struct {
    card    *Card
    resultCh chan bool
}

// single goroutine processes all requests for this player
func (p *Player) run() {
    for req := range p.reqCh {
        result := doPurchase(p, req.card)
        req.resultCh <- result
    }
}

func purchase(p *Player, c *Card) bool {
    resultCh := make(chan bool, 1)
    p.reqCh <- purchaseReq{c, resultCh}
    return <-resultCh
}
*/
