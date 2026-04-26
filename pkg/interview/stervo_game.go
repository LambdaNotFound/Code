package interview

/**
 * Game: Stervo
 *
 * Players collect colored gems (B/W/G/R/Y) and buy cards using gems
 * Each card has a cost (gems required) and a color (the gem color it represents)
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

// "I'm modeling the card cost as a map from Color to int because not every card uses every color
// — a sparse map is cleaner than a fixed-size array.
// The player's gem inventory is the same shape. canPurchase is just a check that for every required color,
// the player has at least that many gems."

type Card struct {
	color Color         // the card's own color (for discounts later)
	cost  map[Color]int // gems required to buy
}

type Player struct {
	gems map[Color]int // current gems
	hand []*Card       // owned cards
}

// "I'm checking canPurchase first, then mutating.
// In a single-threaded context this is fine. Once we add concurrency in part 4,
// this becomes a check-then-act race condition we'll need to fix."

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
