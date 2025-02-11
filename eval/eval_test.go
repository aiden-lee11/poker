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

	bestHand, strength := eval.EvalHand(allCards)
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

	fmt.Printf("strength: %v\n", strength)

}

// handString: [Four of Clubs Two of Hearts Four of Hearts Five of Clubs Nine of Diamonds Jack of Clubs Three of Diamonds]

func TestRandomHandEval(t *testing.T) {
	allCards := []eval.Card{
		//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
		0b00000000_00000100_10000010_00000101,
		0b00000000_00000001_00100000_00000010,
		0b00000000_00000100_00100010_00000101,
		0b00000000_00001000_10000011_00000111,
		0b00000000_10000000_01000111_00010011,
		0b00000010_00000000_10001001_00011101,
		0b00000000_00000010_01000001_00000011,
	}

	bestHand, strength := eval.EvalHand(allCards)
	bestHand.SortHand()

	fmt.Printf("strength: %v\n", strength)
	//expectedBest := eval.Hand{Cards: []eval.Card{
	//	//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
	//	0b00000000_00000001_10000010_00000010, // Two of Clubs
	//	0b00000000_00000010_01000011_00000011, // Three of Diamonds
	//	0b00000000_00000100_00100100_00000101, // Four of Hearts
	//	0b00000000_00001000_00010101_00000111, // Five of Spades
	//	0b00000000_00010000_10000110_00001011, // Six of Clubs
	//}}

	// 	fmt.Printf("bestHand: %v\n", bestHand)
	// 	fmt.Printf("expectedBest: %v\n", expectedBest)

	// for i, card := range bestHand.Cards {
	// 	assert.Equal(t, card, expectedBest.Cards[i], "Expected the straight to be the best hand, got %v", bestHand)
	// }

	// fmt.Printf("strength: %v\n", strength)

}
