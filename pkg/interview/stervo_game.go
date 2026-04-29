package interview

import "sync"

/**
 * Game: Stervo
 *
# Stervo is a card-based board game where players buy cards in exchange for colored gems. In this game, today, we care about two things, gems and cards.

# Players can have any number of gems of five different colors: (B)lue, (W)hite, (G)reen, (R)ed, and (Y)ellow.

# Players can exchange gems for cards. A card appears as such:

# +----------+
# |        G |
# |          |
# |          |
# |  3W      |
# |  2G      |
# |  1R      |
# +----------+

# This indicates that the card costs 3 (W)hite gems, 2 (G)reen gems, and 1 (R)ed. The “G” in the upper right indicates the color of the card (this will be useful later)

# For this entire problem, we want to keep things simple by assuming that there is only one player.

# The data model and structure of the program is up to you.


Task 1:
# We want to write a can_purchase() function such that, given a card and the player's gem collection, it returns true if the player can afford the card, and false otherwise.

Task 2:
# Let's create a function called purchase() that takes as input a card and a player's collection of gems and checks if the player has enough gems to afford the card. If the player can afford the card,
# the function will add it to the player's hand and deduct the card’s cost from the player's gem collection. The function should return "true" if the player can afford the card, and "false" if the player can't afford it.

Task 3:
# Now, we want to introduce the concept of discounts. The color of each card in the player's hand will be used to give a discount in the respective price of that gem’s color when purchasing a new card.
# For example, if the player has 2 white cards and 1 red card in their hand and wants to buy a card that costs 4 white gems and 1 red gem, the discounted price of the card would be 2 white gems.

# Imagine, now, that we are running this code in parallel where there may be multiple requests from the same player to purchase cards at the same time. How would your design change to accommodate this?
*
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
