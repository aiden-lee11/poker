package table_test

import (
	"testing"

	"poker/eval"
	"poker/table"

	"github.com/stretchr/testify/assert"
)

func TestTableFlow(t *testing.T) {
	testTable := table.Table{Deck: eval.UnshuffledDeck}

	// Add players
	testTable.AddPlayer(1000)
	testTable.AddPlayer(1000)

	// Keep deck unshuffled for determinism
	testTable.DistributeCards()

	// Verify hole cards dealt round-robin
	assert.Equal(t, eval.TwoOfSpades, testTable.Players[0].HoleCards[0])
	assert.Equal(t, eval.ThreeOfSpades, testTable.Players[1].HoleCards[0])
	assert.Equal(t, eval.FourOfSpades, testTable.Players[0].HoleCards[1])
	assert.Equal(t, eval.FiveOfSpades, testTable.Players[1].HoleCards[1])

	assert.Len(t, testTable.Deck, 48, "Expected 4 cards to be popped, should have a length of 48 after, got %d\n", len(testTable.Deck))

	testTable.ShowFlopCards()
	assert.Len(t, testTable.CommunityCards, 3, "Expected 3 cards in the community, got %d", len(testTable.CommunityCards))

	testTable.ShowTurnCard()
	assert.Len(t, testTable.CommunityCards, 4, "Expected 4 cards in the community, got %d", len(testTable.CommunityCards))

	testTable.ShowRiverCard()
	assert.Len(t, testTable.CommunityCards, 5, "Expected 5 cards in the community, got %d", len(testTable.CommunityCards))

	// Sanity: evaluate a player's hand returns a rank
	player1 := testTable.Players[0]
	_, strength := player1.EvalHand(&testTable)
	assert.GreaterOrEqual(t, strength, eval.HighCard)
}

func TestCommunityCardNames(t *testing.T) {
	tbl := table.Table{}
	assert.Empty(t, tbl.CommunityCardNames())

	tbl.CommunityCards = []eval.Card{eval.TwoOfClubs, eval.ThreeOfDiamonds}
	names := tbl.CommunityCardNames()
	assert.Equal(t, []string{"Two of Clubs", "Three of Diamonds"}, names)
}

func TestValidBet(t *testing.T) {
	tbl := table.Table{}

	// No previous raise: any bet is valid
	tbl.MostRecentRaise = table.Bet{BetAmount: 0}
	assert.True(t, tbl.ValidBet(1))
	assert.True(t, tbl.ValidBet(100))

	// Previous raise of 50: next valid raise must be >= 100
	tbl.MostRecentRaise = table.Bet{BetAmount: 50}
	assert.False(t, tbl.ValidBet(75))
	assert.True(t, tbl.ValidBet(100))
	assert.True(t, tbl.ValidBet(150))
}

func TestSetPositionsPreFlop(t *testing.T) {
	tbl := table.Table{
		Players: []*table.Player{
			{PlayerID: "A"},
			{PlayerID: "B"},
			{PlayerID: "C"},
		},
		BigBlindIndex: 1,
		BigBlindCost:  20,
		Round:         table.PreFlop,
	}

	tbl.SetPositions()

	assert.Equal(t, "B", tbl.MostRecentRaise.PlayerID)
	assert.Equal(t, 20, tbl.MostRecentRaise.BetAmount)
	assert.Equal(t, table.PreFlop, tbl.MostRecentRaise.Round)
	assert.True(t, tbl.MostRecentRaise.Start)
	assert.Equal(t, 2, tbl.CurrentTurnIndex)
}

func TestSetPositionsPostFlop(t *testing.T) {
	tbl := table.Table{
		Players: []*table.Player{
			{PlayerID: "A"},
			{PlayerID: "B"},
			{PlayerID: "C"},
		},
		BigBlindIndex: 1,
		Round:         table.Flop,
	}

	tbl.SetPositions()

	// Dealer is index (BB-2) mod n = 2
	assert.Equal(t, "C", tbl.MostRecentRaise.PlayerID)
	assert.Equal(t, 0, tbl.MostRecentRaise.BetAmount)
	assert.Equal(t, table.Flop, tbl.MostRecentRaise.Round)
	assert.True(t, tbl.MostRecentRaise.Start)
	// Action starts at small blind: (BB-1) mod n = 0
	assert.Equal(t, 0, tbl.CurrentTurnIndex)
}

func TestSetBigBlindIndexWrap(t *testing.T) {
	tbl := table.Table{
		Players:       []*table.Player{{}, {}, {}},
		BigBlindIndex: 2,
	}
	tbl.SetBigBlindIndex()
	assert.Equal(t, 0, tbl.BigBlindIndex)
}

func TestHoleCardNames(t *testing.T) {
	p := table.Player{HoleCards: [2]eval.Card{eval.AceOfSpades, eval.KingOfHearts}}
	assert.Equal(t, []string{"Ace of Spades", "King of Hearts"}, p.HoleCardNames())
}

func TestHandleEvaluateHands_TieOnBoard(t *testing.T) {
	tbl := table.Table{
		Players: []*table.Player{
			{PlayerID: "P1", PlayingHand: true, HoleCards: [2]eval.Card{eval.AceOfSpades, eval.AceOfHearts}},
			{PlayerID: "P2", PlayingHand: true, HoleCards: [2]eval.Card{eval.KingOfSpades, eval.KingOfHearts}},
		},
		CommunityCards: []eval.Card{eval.TwoOfClubs, eval.ThreeOfClubs, eval.FourOfClubs, eval.FiveOfClubs, eval.SixOfClubs},
	}

	_, winners := tbl.HandleEvaluateHands()
	assert.ElementsMatch(t, []string{"P1", "P2"}, winners)
}

func TestPrintTableDetails_DoesNotPanic(t *testing.T) {
	tbl := &table.Table{}
	// Should not panic
	tbl.PrintTableDetails()
}

func TestShuffleDeckReturnsPermutation(t *testing.T) {
	tbl := &table.Table{}
	original := eval.UnshuffledDeck
	shuffled := tbl.ShuffleDeck(original)

	// Same length
	assert.Equal(t, len(original), len(shuffled))

	// Same multiset of names
	count := map[string]int{}
	for _, c := range original {
		count[c.Name]++
	}
	for _, c := range shuffled {
		count[c.Name]--
	}
	for _, v := range count {
		assert.Equal(t, 0, v)
	}
}

func TestSimulateOdds_DoesNotPanic(t *testing.T) {
	tbl := table.Table{
		Players: []*table.Player{
			{PlayerID: "P1", PlayingHand: true, HoleCards: [2]eval.Card{eval.AceOfSpades, eval.KingOfSpades}},
			{PlayerID: "P2", PlayingHand: true, HoleCards: [2]eval.Card{eval.QueenOfHearts, eval.JackOfHearts}},
		},
		Deck:           eval.UnshuffledDeck,
		CommunityCards: []eval.Card{eval.TwoOfClubs, eval.ThreeOfDiamonds, eval.FourOfHearts},
	}

	// Should run without panic
	tbl.SimulateOdds()
}
