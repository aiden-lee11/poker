package table_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"poker/table"
	"testing"
)

func TestMain(t *testing.T) {
	player1 := table.Player{StackSize: 1000}
	player2 := table.Player{StackSize: 1000}

	testTable := table.Table{Deck: table.BaseStringDeck, Players: []table.Player{player1, player2}}

	testTable.ShuffleDeck()
	testTable.DistributeCards()

	for i, player := range testTable.Players {
		fmt.Printf("player[%d].HoleCards: %v\n", i, player.HoleCards)
	}

	assert.Len(t, testTable.Deck, 48, "Expected 4 cards to be popped, should have a length of 48 after, got %d\n", len(testTable.Deck))

	testTable.ShowFlopCards()

	assert.Len(t, testTable.CommunityCards, 3, "Expected 3 cards in the community got %d", len(testTable.CommunityCards))

	fmt.Printf("testTable.CommunityCards: %v\n", testTable.CommunityCards)

	testTable.ShowTurnCard()

	assert.Len(t, testTable.CommunityCards, 4, "Expected 4 cards in the community got %d", len(testTable.CommunityCards))

	fmt.Printf("testTable.CommunityCards: %v\n", testTable.CommunityCards)

	testTable.ShowRiverCard()

	assert.Len(t, testTable.CommunityCards, 5, "Expected 5 cards in the community got %d", len(testTable.CommunityCards))

	fmt.Printf("testTable.CommunityCards: %v\n", testTable.CommunityCards)

}
