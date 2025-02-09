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

	assert.Equal(t, 48, len(testTable.Deck), "Expected 4 cards to be popped, should have a length of 48 after, got %d\n", len(testTable.Deck))
}
