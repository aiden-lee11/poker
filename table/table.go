package table

import (
	"fmt"
	"math/rand"
	"poker/eval"
)

type Card string

type Player struct {
	StackSize   int
	PlayerID    string
	HoleCards   [2]Card
	PlayingHand bool
}

type Table struct {
	ID               string
	Deck             []Card
	Players          []*Player
	CommunityCards   []Card
	PotSize          int
	CurrentTurnIndex int
	// i think this is how im gonna handle when orbits should be over
	MostRecentRaise string
}

// Need to incorporate this type of player hand struct with our eval
// only call this after flop should make other func for preflop vals
func (player *Player) EvalHand(table *Table) (eval.Hand, int) {
	if table == nil {
		fmt.Println("invalid table reference")
		return eval.Hand{}, 0
	}

	handString := append([]Card{}, player.HoleCards[:]...)
	handString = append(handString, table.CommunityCards...)

	handInts := make([]eval.Card, len(handString))
	for i, card := range handString {
		handInts[i] = eval.Card(CardToBits[card])
	}

	return eval.EvalHand(handInts)
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
	table.Players = append(table.Players, &Player{StackSize: stackSize})
}

func (table *Table) popCardFromDeck(numToPop int) {
	table.Deck = table.Deck[numToPop:]
}

// for these we always index at 1 because we burn first
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
