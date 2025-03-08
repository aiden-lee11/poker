package table_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"poker/table"
)

func TestMain(t *testing.T) {
	testTable := table.Table{Deck: table.BaseStringDeck}

	// Add players to the table using AddPlayer so they get a reference to the table if needed
	testTable.AddPlayer(1000)
	testTable.AddPlayer(1000)

	testTable.ShuffleDeck()
	testTable.DistributeCards()

	for i, player := range testTable.Players {
		fmt.Printf("player[%d].HoleCards: %v\n", i, player.HoleCards)
	}

	assert.Len(t, testTable.Deck, 48, "Expected 4 cards to be popped, should have a length of 48 after, got %d\n", len(testTable.Deck))

	testTable.ShowFlopCards()
	assert.Len(t, testTable.CommunityCards, 3, "Expected 3 cards in the community, got %d", len(testTable.CommunityCards))
	// fmt.Printf("testTable.CommunityCards: %v\n", testTable.CommunityCards)

	testTable.ShowTurnCard()
	assert.Len(t, testTable.CommunityCards, 4, "Expected 4 cards in the community, got %d", len(testTable.CommunityCards))
	// fmt.Printf("testTable.CommunityCards: %v\n", testTable.CommunityCards)

	testTable.ShowRiverCard()
	assert.Len(t, testTable.CommunityCards, 5, "Expected 5 cards in the community, got %d", len(testTable.CommunityCards))
	fmt.Printf("testTable.CommunityCards: %v\n", testTable.CommunityCards)

	player1 := &testTable.Players[0]

	// Updated EvalHand call
	player1Hand, strength := player1.EvalHand()

	fmt.Printf("player1Hand: %v\n", player1Hand)
	fmt.Printf("strength: %v\n", strength)
}
