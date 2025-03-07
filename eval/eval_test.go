package eval_test

import (
	// "fmt"
	// "github.com/stretchr/testify/assert"
	"reflect"
	"poker/eval"
	"testing"
)

//func TestMain(t *testing.T) {
//	cards := []eval.Card{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}

//	combos := eval.GenerateCombinations(cards)

//	for i, combo := range combos {
//		fmt.Printf("The %d combo is %v\n", i, combo)
//	}
//}

func TestEval(t *testing.T) {
	aceHigh := eval.Hand{Cards: []eval.Card{
		eval.KingOfDiamonds,
		eval.FiveOfSpades,
		eval.JackOfClubs,
		eval.QueenOfHearts,
		eval.AceOfSpades,
	}}

	straight := eval.Hand{Cards: []eval.Card{
		eval.TwoOfClubs,
		eval.ThreeOfDiamonds,
		eval.FourOfHearts,
		eval.FiveOfSpades,
		eval.SixOfClubs,
	}}

	aceHighEval := aceHigh.EvaluateHand()
	if aceHighEval != eval.HighCard {
		t.Errorf(`Expected %v to be high card but got %v`, aceHigh, aceHighEval)
	}

	straightEval := straight.EvaluateHand()
	if straightEval != eval.Straight {
		t.Errorf(`Expected %v to be straight but got %v`, straight, straightEval)
	}

	if straightEval <= aceHighEval {
		t.Errorf(`straight hand %v ranks lower than ace high hand %v, which is incorrect`, straight.Stringify(), aceHigh.Stringify())
	}
}

func TestFindBestHand(t *testing.T) {
	allCards := []eval.Card{
		eval.TwoOfClubs,
		eval.ThreeOfDiamonds,
		eval.FourOfHearts,
		eval.FiveOfSpades,
		eval.SixOfClubs,
		eval.KingOfDiamonds,
		eval.QueenOfSpades,
	}

	bestHand, strength := eval.EvalHand(allCards)
	if strength != eval.Straight {
		t.Errorf(`Expected straight but got %v`, strength)
	}

	bestHand.SortHand()
	expectedBest := eval.Hand{Cards: []eval.Card{
		eval.SixOfClubs,
		eval.FiveOfSpades,
		eval.FourOfHearts,
		eval.ThreeOfDiamonds,
		eval.TwoOfClubs,
	}}

	if !reflect.DeepEqual(bestHand.Cards, expectedBest.Cards) {
		t.Errorf(`Expected best hand to be %v but got %v`, expectedBest.Stringify(), bestHand.Stringify())
	}
}

// handString: [Four of Clubs Two of Hearts Four of Hearts Five of Clubs Nine of Diamonds Jack of Clubs Three of Diamonds]

// func TestRandomHandEval(t *testing.T) {
	// allCards := []eval.CardBits{
	// 	//xxxAKQJT 98765432 CDHSrrrr xxPPPPPP
	// 	0b00000000_00000100_10000010_00000101,
	// 	0b00000000_00000001_00100000_00000010,
	// 	0b00000000_00000100_00100010_00000101,
	// 	0b00000000_00001000_10000011_00000111,
	// 	0b00000000_10000000_01000111_00010011,
	// 	0b00000010_00000000_10001001_00011101,
	// 	0b00000000_00000010_01000001_00000011,
	// }

	// bestHand, strength := eval.EvalHand(allCards)
	// bestHand.SortHand()

	// fmt.Printf("strength: %v\n", strength)
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
// }
