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

// Need to incorporate this type of player hand struct with our eval
// func (player *Player) evalHand() int {
// 	cards :=
// }

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

func (table *Table) ShowFlopCards() {
	table.CommunityCards = append(table.CommunityCards, table.Deck[1:4]...)

	table.popCardFromDeck(4)
}

func (table *Table) ShowTurnCard() {
	table.CommunityCards = append(table.CommunityCards, table.Deck[1])
	table.popCardFromDeck(2)
}

func (table *Table) ShowRiverCard() {
	table.CommunityCards = append(table.CommunityCards, table.Deck[1])
	table.popCardFromDeck(2)
}

func (table *Table) AddPlayer(stackSize int) {
	table.Players = append(table.Players, Player{StackSize: stackSize})
}

func (table *Table) popCardFromDeck(numToPop int) {
	table.Deck = table.Deck[numToPop:]
}
