package table

import (
	"fmt"
	"math/rand"
	"poker/eval"
	"reflect"
	"slices"
	"strings"
)

type Player struct {
	StackSize   int
	PlayerID    string
	HoleCards   [2]eval.Card
	PlayingHand bool
	WinOdds     float32
	SplitOdds   float32
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

func (table *Table) ShuffleDeck(deck []eval.Card) []eval.Card {
	shuffledDeck := make([]eval.Card, len(deck))
	copy(shuffledDeck, deck)

	rand.Shuffle(len(shuffledDeck), func(i, j int) {
		shuffledDeck[i], shuffledDeck[j] = shuffledDeck[j], shuffledDeck[i]
	})

	return shuffledDeck
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

func (t *Table) CommunityCardNames() []string {
	names := eval.CardNames(t.CommunityCards)

	if len(names) == 0 {
		return []string{}
	}

	return names
}

func (t *Table) ValidBet(betSize int) bool {
	return (t.MostRecentRaise.BetAmount == 0) || (t.MostRecentRaise.BetAmount*2 <= betSize)
}

// thinking for sims we do a map of players to wins
// loop over like 1000 iterations for simulations with shuffled decks each time
// eval the winner for each and then add that to the counter at the end loop over the map
// and update the players odds field with (num wins / sims run)
func (t *Table) SimulateOdds() {
	numSimulations := 1000
	wins := make(map[string]int)
	splits := make(map[string]int)

	// change to round based lol
	n := len(t.CommunityCards)

	fmt.Printf("cards that are static t.CommunityCards: %v\n", t.CommunityCards)
	for i := 0; i < numSimulations; i++ {
		newDeck := t.ShuffleDeck(t.Deck)
		before := t.CommunityCards

		if n == 3 {
			t.CommunityCards = append(t.CommunityCards, []eval.Card{newDeck[0], newDeck[2]}...)
		} else if n == 4 {
			t.CommunityCards = append(t.CommunityCards, newDeck[0])
		}

		_, winners := t.HandleEvaluateHands()
		if len(winners) > 1 {
			for _, winner := range winners {
				splits[winner] += 1
				wins[winner] += 1
			}
		} else {
			wins[winners[0]] += 1
		}

		t.CommunityCards = before
	}

	for _, player := range t.Players {
		player.WinOdds = float32(wins[player.PlayerID]) / float32(numSimulations)
		player.SplitOdds = float32(splits[player.PlayerID]) / float32(numSimulations)
		fmt.Printf("Player %s has %f odds of winning and %f of splitting\n", player.PlayerID, player.WinOdds, player.SplitOdds)
	}
}

// similar to above this should be handled by server eventually
// would be triggered when a table has done the full round after river
// or when all active players cannot bet further IN TERMS OF DETERMINING WINNER
// returns stringified hand and list of winners
func (t *Table) HandleEvaluateHands() ([]string, []string) {
	winners := []string{t.Players[0].PlayerID}
	// need to handle if player 1 has the best hand but wasnt playing it lol
	resHand, highest := t.Players[0].EvalHand(t)

	for i := 1; i < len(t.Players); i++ {
		player := t.Players[i]

		if !player.PlayingHand {
			continue
		}

		hand, val := player.EvalHand(t)
		// val here means like two pair one pair etc

		if val > highest {
			highest = val
			resHand = hand
			winners = []string{player.PlayerID}
		} else if val == highest {
			best := eval.HandleTie(hand, resHand)
			if !reflect.DeepEqual(resHand, best) {
				winners = []string{player.PlayerID}
				resHand = best
			} else {
				winners = append(winners, player.PlayerID)
			}
		}
	}

	return resHand.Stringify(), winners
}
