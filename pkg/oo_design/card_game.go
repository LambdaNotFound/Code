package oodesign

/** Card Game

Game
 ├── players []*Player      ← was hardcoded to 2, now N
 ├── deal()                 ← already N-player, no change
 ├── playRound()            ← extract one round into its own method
 ├── play()                 ← loop until done
 └── winner() string        ← generalized across N players

Player
 ├── Name string
 ├── Hand []int
 └── WinPile []int

*/

type Player struct {
	Name    string
	Hand    []int
	WinPile []int
}

type Game struct {
	players  []*Player
	numCards int
}

func (g *Game) deal()              // unchanged
func (g *Game) playRound()         // NEW — extract the inner loop body
func (g *Game) allHaveCards() bool // NEW — termination condition
func (g *Game) winner() string     // generalize from 2 → N
func (g *Game) tiebreak() string   // generalize from 2 → N

func NewGame(names []string, numCards int) *Game {
	players := make([]*Player, len(names))
	for i, name := range names {
		players[i] = &Player{Name: name}
	}
	return &Game{players: players, numCards: numCards}
}

func (g *Game) Play() string {
	g.deal()
	for g.allHaveCards() {
		g.playRound()
	}
	return g.winner()
}

// Single responsibility: Each method does one thing: deal, play a round, find winner
// Open/closed Extending to N players doesn't modify Play() or deal() — only playRound
// Encapsulation Game owns all state; caller just does NewGame(...).Play()

// 2-player playRound (question 1)
func (g *Game) playRound() {
	card0 := g.players[0].Hand[0]
	card1 := g.players[1].Hand[0]
	g.players[0].Hand = g.players[0].Hand[1:]
	g.players[1].Hand = g.players[1].Hand[1:]

	if card0 > card1 {
		g.players[0].WinPile = append(g.players[0].WinPile, card0, card1)
	} else {
		g.players[1].WinPile = append(g.players[1].WinPile, card0, card1)
	}
}

// N-player playRound (question 2) — same signature, just generalize the body
func (g *Game) playRound() {
	highCard, winner := 0, g.players[0]
	var roundCards []int

	for _, p := range g.players { // ← only change: loop instead of hardcode
		card := p.Hand[0]
		p.Hand = p.Hand[1:]
		roundCards = append(roundCards, card)
		if card > highCard {
			highCard = card
			winner = p
		}
	}
	winner.WinPile = append(winner.WinPile, roundCards...)
}
