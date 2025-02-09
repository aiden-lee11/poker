package table

import (
	"math/rand"
)

type Card string

type Player struct {
	StackSize   int
	HoleCards   [2]Card
	PlayingHand bool
	// Position    int
}

type Table struct {
	Deck           []Card
	Players        []Player
	CommunityCards []Card
	PotSize        int
}

func (table *Table) ShuffleDeck() {
	shuffledDeck := make([]Card, len(BaseStringDeck))
	copy(shuffledDeck, BaseStringDeck)

	rand.Shuffle(len(shuffledDeck), func(i, j int) {
		shuffledDeck[i], shuffledDeck[j] = shuffledDeck[j], shuffledDeck[i]
	})

	table.Deck = shuffledDeck
}

func (table *Table) DistributeCards() {
	numPlayers := len(table.Players)

	for i := 0; i < (numPlayers * 2); i++ {
		table.Players[i%numPlayers].HoleCards[i/numPlayers] = table.Deck[i]
	}

	table.popCardFromDeck(numPlayers * 2)
}

func (table *Table) AddPlayer(stackSize int) {
	table.Players = append(table.Players, Player{StackSize: stackSize})
}

func (table *Table) popCardFromDeck(numToPop int) {
	table.Deck = table.Deck[numToPop:]
}

// // table package?
// func bet()

// func fold()

// func check()
