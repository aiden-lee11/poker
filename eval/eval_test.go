package eval_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"poker/eval"
	"testing"
)

func TestMain(t *testing.T) {
	cards := []eval.Card{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}

	combos := eval.GenerateCombinations(cards)

	for i, combo := range combos {
		fmt.Printf("The %d combo is %v\n", i, combo)
	}
}

func TestEval(t *testing.T) {
	aceHigh := eval.Hand{Cards: []eval.Card{
		//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
		0b00001000_00000000_01001011_00100101, // King of Diamonds
		0b00000000_00001000_00010011_00000111, // Five of Spades
		0b00000010_00000000_10001001_00011101, // Jack of Clubs
		0b00000100_00000000_00011010_00011111, // Queen of Hearts
		0b00010000_00000000_00011100_00101001, // Ace of Spades
	}}

	straight := eval.Hand{Cards: []eval.Card{
		//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
		0b00000000_00000001_10000010_00000010, // Two of Clubs
		0b00000000_00000010_01000011_00000011, // Three of Diamonds
		0b00000000_00000100_00100100_00000101, // Four of Hearts
		0b00000000_00001000_00010101_00000111, // Five of Spades
		0b00000000_00010000_10000110_00001011, // Six of Clubs
	}}

	aceHighEval := aceHigh.EvaluateHand()
	straightEval := straight.EvaluateHand()

	res := straightEval > aceHighEval
	assert.True(t, res, "Expected straight to be higher, but that shit was not higher, straight val %v, aceHigh val %v", straightEval, aceHighEval)
}

func TestFindBestHand(t *testing.T) {
	allCards := []eval.Card{
		//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
		0b00000000_00000001_10000010_00000010, // Two of Clubs
		0b00000000_00000010_01000011_00000011, // Three of Diamonds
		0b00000000_00000100_00100100_00000101, // Four of Hearts
		0b00000000_00001000_00010101_00000111, // Five of Spades
		0b00000000_00010000_10000110_00001011, // Six of Clubs
		0b00001000_00000000_01001011_00100101, // King of Diamonds
		0b00010000_00000000_00011100_00101001, // Ace of Spades
	}

	combinations := eval.GenerateCombinations(allCards)

	bestHand := eval.FindBestHand(combinations)
	bestHand.SortHand()

	expectedBest := eval.Hand{Cards: []eval.Card{
		//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
		0b00000000_00000001_10000010_00000010, // Two of Clubs
		0b00000000_00000010_01000011_00000011, // Three of Diamonds
		0b00000000_00000100_00100100_00000101, // Four of Hearts
		0b00000000_00001000_00010101_00000111, // Five of Spades
		0b00000000_00010000_10000110_00001011, // Six of Clubs
	}}

	fmt.Printf("bestHand: %v\n", bestHand)
	fmt.Printf("expectedBest: %v\n", expectedBest)

	for i, card := range bestHand.Cards {
		assert.Equal(t, card, expectedBest.Cards[i], "Expected the straight to be the best hand, got %v", bestHand)
	}

}
