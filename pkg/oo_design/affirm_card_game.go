package oodesign

import "math/rand/v2"

/** Card Game

"""
this is a two player card game
the game starts with a deck of 52 cards represented as unique integers [1...52]
the cards are randomly shuffled and then dealt out to both players evenly
on each turn:
both players turn over their top-most card
the player with the higher valued card takes the cards and puts them in their scoring pile (scoring 1 point per card)
this continues until all the players have no cards left
the player with the highest score wins
if they have the same number of cards in their win pile, tiebreaker goes to the player with the highest card in their win pile

Be able to play the game with N players.
An input to the game will now be a list of strings (of length N) indicating the player names.
The deck contains M cards of distinct integers.
It is not guaranteed M % N == 0. If there are leftover cards they should randomly be handed out to remaining players.
i.e. with 17 cards and 5 people: 2 people get 4 cards and 3 get 3 cards
For example the input: game(["Joe", "Jill", "Bob"], 5) would be a game between 3 players and 5 cards.
you should print the name of the player that won the game.
"""

Game
 ├── players []*Player      ← was hardcoded to 2, now N
 ├── deal()                 ← already N-player, no change
 ├── playRound()            ← extract one round into its own method
 ├── play()                 ← loop until done
 └── winner() string        ← generalized across N players

func (g *CardGame) deal()              // unchanged
func (g *CardGame) playRound()         // NEW — extract the inner loop body
func (g *CardGame) allHaveCards() bool // NEW — termination condition
func (g *CardGame) winner() string     // generalize from 2 → N

Player
 ├── Name string
 ├── Hand []int
 └── WinPile []int

*/

type CardGamePlayer struct {
	Name    string
	Hand    []int
	WinPile []int // winner takes all
}

type CardGame struct {
	players  []*CardGamePlayer
	numCards int
}

func NewGame(names []string, numCards int) *CardGame { // 52 cards, [1 ... 52]
	players := make([]*CardGamePlayer, len(names))
	for i, name := range names {
		players[i] = &CardGamePlayer{Name: name}
	}
	return &CardGame{players: players, numCards: numCards}
}

func (g *CardGame) deal() {
	// Create and shuffle deck
	deck := make([]int, g.numCards)
	for i := range deck {
		deck[i] = i + 1
	}
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	// Deal cards round-robin
	for i, card := range deck {
		g.players[i%len(g.players)].Hand = append(g.players[i%len(g.players)].Hand, card)
	}
}

// Works for N=2 and N=any — just loops over g.players
func (g *CardGame) allHaveCards() bool {
	for _, p := range g.players {
		if len(p.Hand) == 0 {
			return false
		}
	}
	return true
}

// 2-player playRound (question 1)
func (g *CardGame) playRound2PlayerMode() {
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
func (g *CardGame) playRoundNPlayerMode() {
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

// N-player playRound (question 2) — same signature, just generalize the body
func (g *CardGame) winner() string {
	// Find winner by pile size
	best := g.players[0]
	for _, p := range g.players[1:] {
		if len(p.WinPile) > len(best.WinPile) {
			best = p
		}
	}

	// Collect all tied players
	var tied []*CardGamePlayer
	for _, p := range g.players {
		if len(p.WinPile) == len(best.WinPile) {
			tied = append(tied, p)
		}
	}

	if len(tied) == 1 {
		return tied[0].Name
	}

	maxCard := func(pile []int) int {
		max := 0
		for _, c := range pile {
			if c > max {
				max = c
			}
		}
		return max
	}

	tiebreak := func(players []*CardGamePlayer) string {
		best := players[0]
		for _, p := range players[1:] {
			if maxCard(p.WinPile) > maxCard(best.WinPile) {
				best = p
			}
		}
		return best.Name
	}

	// Tiebreaker among tied players
	return tiebreak(tied)
}

// Single responsibility: Each method does one thing: deal, play a round, find winner
// Open/closed Extending to N players doesn't modify Play() or deal() — only playRound
// Encapsulation Game owns all state; caller just does NewGame(...).Play()
func (g *CardGame) Play() string {
	g.deal()
	for g.allHaveCards() {
		g.playRoundNPlayerMode()
	}
	return g.winner()
}
