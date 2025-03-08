package table

import (
	"fmt"
	"math/rand"
	"poker/eval"
	"slices"
	"strings"
)

type Player struct {
	StackSize   int
	PlayerID    string
	HoleCards   [2]eval.Card
	PlayingHand bool
}

type BettingRound int

const (
	PreFlop BettingRound = iota
	Flop
	Turn
	River
)

type Bet struct {
	PlayerID  string
	BetAmount int
	Round     BettingRound
	Start     bool
}

type Table struct {
	ID             string
	Deck           []eval.Card
	Players        []*Player
	CommunityCards []eval.Card
	PotSize        int
	// should be big blind + 1 for preflop
	// then small blind for all other rounds
	CurrentTurnIndex int
	BigBlindIndex    int
	MostRecentRaise  Bet          // could be empty ("") when not set
	Round            BettingRound // new field to track the current betting round
	SmallBlindCost   int
	BigBlindCost     int
}

// Need to incorporate this type of player hand struct with our eval
// only call this after flop should make other func for preflop vals
func (player *Player) EvalHand(table *Table) (eval.Hand, int) {
	if table == nil {
		fmt.Println("invalid table reference")
		return eval.Hand{}, 0
	}
	holeCardsAndCommunityCards := slices.Concat(player.HoleCards[:], table.CommunityCards)
	return eval.EvalHand(holeCardsAndCommunityCards)
}

func (p Player) HoleCardNames() []string {
	return eval.CardNames(p.HoleCards[:])
}

func (table *Table) ShuffleDeck() {
	shuffledDeck := make([]eval.Card, len(eval.UnshuffledDeck))
	copy(shuffledDeck, eval.UnshuffledDeck)

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

// SetDefaultMostRecentRaise sets the default most recent raise based on the current betting round.
// For the first round (PreFlop), we default to the big blind
// For other rounds, we default to the small blind
// set current index to start with utg for preflop and sb for postflop
func (t *Table) SetPositions() {
	numPlayers := len(t.Players)
	if numPlayers == 0 {
		return
	}
	if t.Round == PreFlop {
		t.MostRecentRaise = Bet{PlayerID: t.Players[t.BigBlindIndex].PlayerID, BetAmount: t.BigBlindCost, Round: PreFlop, Start: true}
		t.CurrentTurnIndex = (t.BigBlindIndex + 1) % numPlayers
	} else {
		smallBlindIndex := mod(t.BigBlindIndex-1, numPlayers)
		dealerIndex := mod(t.BigBlindIndex-2, numPlayers)
		t.MostRecentRaise = Bet{PlayerID: t.Players[dealerIndex].PlayerID, BetAmount: 0, Round: t.Round, Start: true}
		t.CurrentTurnIndex = smallBlindIndex
	}
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func (t *Table) SetBigBlindIndex() {
	numPlayers := len(t.Players)
	if numPlayers == 0 {
		return
	}

	t.BigBlindIndex = (t.BigBlindIndex + 1) % numPlayers
}

func (t *Table) PrintTableDetails() {
	var sb strings.Builder

	sb.WriteString("Table Details:\n")
	sb.WriteString("==============\n")

	sb.WriteString(fmt.Sprintf("Table ID: %s\n", t.ID))

	// Print Deck
	sb.WriteString("\nDeck:\n")
	if len(t.Deck) == 0 {
		sb.WriteString("  [Empty]\n")
	} else {
		for i, card := range t.Deck {
			sb.WriteString(fmt.Sprintf("  %d: %v\n", i, card))
		}
	}

	// Print Players
	sb.WriteString("\nPlayers:\n")
	if len(t.Players) == 0 {
		sb.WriteString("  [No Players]\n")
	} else {
		for i, player := range t.Players {
			sb.WriteString(fmt.Sprintf("  %d: %v\n", i, player))
		}
	}

	// Print Community Cards
	sb.WriteString("\nCommunity Cards:\n")
	if len(t.CommunityCards) == 0 {
		sb.WriteString("  [No Community Cards]\n")
	} else {
		for i, card := range t.CommunityCards {
			sb.WriteString(fmt.Sprintf("  %d: %v\n", i, card))
		}
	}

	// Print Other Details
	sb.WriteString("\nGame Status:\n")
	sb.WriteString(fmt.Sprintf("  Pot Size: %d\n", t.PotSize))
	sb.WriteString(fmt.Sprintf("  Current Turn Index: %d\n", t.CurrentTurnIndex))
	sb.WriteString(fmt.Sprintf("  Big Blind Index: %d\n", t.BigBlindIndex))
	sb.WriteString(fmt.Sprintf("  Most Recent Raise: %v\n", t.MostRecentRaise))
	sb.WriteString(fmt.Sprintf("  Current Round: %v\n", t.Round))
	sb.WriteString(fmt.Sprintf("  Small Blind Cost: %d\n", t.SmallBlindCost))
	sb.WriteString(fmt.Sprintf("  Big Blind Cost: %d\n", t.BigBlindCost))

	fmt.Println(sb.String())
}

func (t Table) CommunityCardNames() []string {
	return eval.CardNames(t.CommunityCards)
}
