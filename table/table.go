package table

import (
	"fmt"
	"math/rand"
	"poker/eval"
)

type Card string

type Player struct {
	StackSize   int
	HoleCards   [2]Card
	PlayingHand bool
	Table       *Table
}

type Table struct {
	Deck           []Card
	Players        []Player
	CommunityCards []Card
	PotSize        int
}

// Need to incorporate this type of player hand struct with our eval
// only call this after flop should make other func for preflop vals
func (player *Player) EvalHand() (eval.Hand, int) {
	if player.Table == nil {
		fmt.Println("Player is not associated with a table.")
		return eval.Hand{}, 0
	}

	handString := append([]Card{}, player.HoleCards[:]...)
	handString = append(handString, player.Table.CommunityCards...)

	handInts := make([]eval.Card, len(handString))
	for i, card := range handString {
		handInts[i] = eval.Card(CardToBits[card])
	}

	fmt.Printf("handString: %v\n", handString)

	// fmt.Print("handInts: [\n")
	// for i, card := range handInts {
	// 	if i > 0 {
	// 		fmt.Print(", \n")
	// 	}
	// 	fmt.Printf("0b%032b", card)
	// }
	// fmt.Println("]")
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
	table.Players = append(table.Players, Player{StackSize: stackSize, Table: table})
}

func (table *Table) popCardFromDeck(numToPop int) {
	table.Deck = table.Deck[numToPop:]
}
