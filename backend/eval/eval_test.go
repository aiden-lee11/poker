package eval_test

import (
	"reflect"
	"testing"

	"poker/eval"

	"github.com/stretchr/testify/assert"
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

func TestGenerateCombinations_LengthAndSize(t *testing.T) {
	allCards := []eval.Card{
		eval.TwoOfClubs,
		eval.ThreeOfDiamonds,
		eval.FourOfHearts,
		eval.FiveOfSpades,
		eval.SixOfClubs,
		eval.KingOfDiamonds,
		eval.QueenOfSpades,
	}

	combinations := eval.GenerateCombinations(allCards)
	assert.Len(t, combinations, 21)
	for _, c := range combinations {
		assert.Len(t, c.Cards, 5)
	}
}

func TestGenerateCombinations_BadInput(t *testing.T) {
	bad := []eval.Card{eval.TwoOfClubs}
	assert.Nil(t, eval.GenerateCombinations(bad))
}

func TestSortHand_SortsDescending(t *testing.T) {
	unsorted := eval.Hand{Cards: []eval.Card{
		eval.TwoOfClubs,
		eval.AceOfHearts,
		eval.NineOfSpades,
		eval.FiveOfDiamonds,
		eval.KingOfClubs,
	}}
	unsorted.SortHand()
	expected := []eval.Card{
		eval.AceOfHearts,
		eval.KingOfClubs,
		eval.NineOfSpades,
		eval.FiveOfDiamonds,
		eval.TwoOfClubs,
	}
	assert.Equal(t, expected, unsorted.Cards)
}

func TestHandleTie_ChoosesHigher(t *testing.T) {
	high := eval.Hand{Cards: []eval.Card{
		eval.AceOfSpades,
		eval.KingOfHearts,
		eval.NineOfSpades,
		eval.FiveOfDiamonds,
		eval.ThreeOfClubs,
	}}
	low := eval.Hand{Cards: []eval.Card{
		eval.KingOfSpades,
		eval.QueenOfHearts,
		eval.NineOfClubs,
		eval.FiveOfHearts,
		eval.ThreeOfDiamonds,
	}}
	best := eval.HandleTie(high, low)
	assert.Equal(t, high, best)
}

func TestEvaluateHand_Categories(t *testing.T) {
	// One Pair
	onePair := eval.Hand{Cards: []eval.Card{
		eval.TwoOfClubs, eval.TwoOfDiamonds, eval.SevenOfSpades, eval.NineOfHearts, eval.AceOfSpades,
	}}
	assert.Equal(t, eval.OnePair, onePair.EvaluateHand())

	// Two Pair
	twoPair := eval.Hand{Cards: []eval.Card{
		eval.TwoOfClubs, eval.TwoOfDiamonds, eval.FiveOfClubs, eval.FiveOfDiamonds, eval.AceOfHearts,
	}}
	assert.Equal(t, eval.TwoPair, twoPair.EvaluateHand())

	// Trips
	trips := eval.Hand{Cards: []eval.Card{
		eval.ThreeOfClubs, eval.ThreeOfDiamonds, eval.ThreeOfHearts, eval.NineOfSpades, eval.KingOfDiamonds,
	}}
	assert.Equal(t, eval.Trips, trips.EvaluateHand())

	// Straight (already covered in TestEval, but keep another)
	straight := eval.Hand{Cards: []eval.Card{
		eval.SixOfClubs, eval.FiveOfSpades, eval.FourOfHearts, eval.ThreeOfDiamonds, eval.TwoOfClubs,
	}}
	assert.Equal(t, eval.Straight, straight.EvaluateHand())

	// Flush
	flush := eval.Hand{Cards: []eval.Card{
		eval.TwoOfSpades, eval.FourOfSpades, eval.SixOfSpades, eval.NineOfSpades, eval.KingOfSpades,
	}}
	assert.Equal(t, eval.Flush, flush.EvaluateHand())

	// Full House
	fullHouse := eval.Hand{Cards: []eval.Card{
		eval.ThreeOfClubs, eval.ThreeOfDiamonds, eval.ThreeOfHearts, eval.FiveOfClubs, eval.FiveOfDiamonds,
	}}
	assert.Equal(t, eval.FullHouse, fullHouse.EvaluateHand())

	// Quads
	quads := eval.Hand{Cards: []eval.Card{
		eval.TwoOfSpades, eval.TwoOfHearts, eval.TwoOfDiamonds, eval.TwoOfClubs, eval.AceOfHearts,
	}}
	assert.Equal(t, eval.Quads, quads.EvaluateHand())

	// Straight Flush
	straightFlush := eval.Hand{Cards: []eval.Card{
		eval.TwoOfClubs, eval.ThreeOfClubs, eval.FourOfClubs, eval.FiveOfClubs, eval.SixOfClubs,
	}}
	assert.Equal(t, eval.StraightFlush, straightFlush.EvaluateHand())
}

func TestCardNamesAndStringify(t *testing.T) {
	cards := []eval.Card{eval.AceOfSpades, eval.TwoOfHearts, eval.KingOfClubs}
	names := eval.CardNames(cards)
	assert.Equal(t, []string{"Ace of Spades", "Two of Hearts", "King of Clubs"}, names)

	h := eval.Hand{Cards: cards}
	assert.Equal(t, names, h.Stringify())
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
