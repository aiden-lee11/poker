package eval

import (
	"fmt"
	"log"
	"sort"
)

// Ported version of Paul Senzee's hand evaluator found at https://github.com/christophschmalhofer/poker/blob/master/XPokerEval/XPokerEval.CactusKev.PerfectHash/fast_eval.cpp

// Card representation
// +--------+--------+--------+--------+
// |xxxbbbbb|bbbbbbbb|cdhsrrrr|xxpppppp|
// +--------+--------+--------+--------+
// xxxAKQJT 98765432 CDHSrrrr xxPPPPPP

type Card int32

type Hand struct {
	Cards []Card
}

const (
	HighCard int = iota
	OnePair
	TwoPair
	Trips
	Straight
	Flush
	FullHouse
	Quads
	StraightFlush
	RoyalFlush
)

const (
	rankMask  = 0x00000F00
	suitMask  = 0x0000F000
	primeMask = 0x000000FF
)

// func for generating combos, evaluating a single 5 card hand strength, evaluating best 5 card hand out 7
func GenerateCombinations(cards []Card) []Hand {
	if len(cards) != 7 {
		fmt.Println("Expected 7 cards in the combo generator, got:", len(cards))
		return nil
	}

	var combos []Hand

	for _, indices := range comboIndices {
		var combo []Card
		for _, index := range indices {
			combo = append(combo, cards[index])
		}

		hand := Hand{Cards: combo}
		combos = append(combos, hand)
	}

	return combos
}

var comboIndices = [][]int{
	{0, 1, 2, 3, 4}, {0, 1, 2, 3, 5}, {0, 1, 2, 3, 6},
	{0, 1, 2, 4, 5}, {0, 1, 2, 4, 6}, {0, 1, 2, 5, 6},
	{0, 1, 3, 4, 5}, {0, 1, 3, 4, 6}, {0, 1, 3, 5, 6},
	{0, 1, 4, 5, 6}, {0, 2, 3, 4, 5}, {0, 2, 3, 4, 6},
	{0, 2, 3, 5, 6}, {0, 2, 4, 5, 6}, {0, 3, 4, 5, 6},
	{1, 2, 3, 4, 5}, {1, 2, 3, 4, 6}, {1, 2, 3, 5, 6},
	{1, 2, 4, 5, 6}, {1, 3, 4, 5, 6}, {2, 3, 4, 5, 6},
}

func FindBestHand(combinations []Hand) Hand {
	// Assuming that combinations is not empty, otherwise you would want to check that first
	if len(combinations) != 21 {
		log.Fatalf("Expected length of 21 combinations got %d", len(combinations))
	}
	resHand := combinations[0]

	for i := 1; i < len(combinations); i++ {
		hand := combinations[i]
		fmt.Printf("resHand: %v\n", resHand)
		fmt.Printf("hand: %v\n", hand)

		if curStrength := hand.EvaluateHand(); curStrength > resHand.EvaluateHand() {
			resHand = hand
		}
	}

	return resHand
}

func EvalHand(hand []Card) (Hand, int) {
	combinations := GenerateCombinations(hand)

	bestHand := FindBestHand(combinations)

	return bestHand, bestHand.EvaluateHand()
}

func HandRank(val int16) int {
	if val > 6185 {
		return (HighCard)
	} // 1277 high card
	if val > 3325 {
		return (OnePair)
	} // 2860 one pair
	if val > 2467 {
		return (TwoPair)
	} //  858 two pair
	if val > 1609 {
		return (Trips)
	} //  858 three-kind
	if val > 1599 {
		return (Straight)
	} //   10 straights
	if val > 322 {
		return (Flush)
	} // 1277 flushes
	if val > 166 {
		return (FullHouse)
	} //  156 full house
	if val > 10 {
		return (Quads)
	} //  156 four-kind
	return (StraightFlush) //   10 straight-flushes
}

func (hand *Hand) getFlushStraightIndex() int16 {
	return int16((hand.Cards[0] | hand.Cards[1] | hand.Cards[2] | hand.Cards[3] | hand.Cards[4]) >> 16)
}

func (hand *Hand) getPrime() uint {
	return uint((hand.Cards[0] & primeMask) * (hand.Cards[1] & primeMask) * (hand.Cards[2] & primeMask) * (hand.Cards[3] & primeMask) * (hand.Cards[4] & primeMask))
}

func (hand *Hand) EvaluateHand() int {
	q := hand.getFlushStraightIndex()

	if hand.isFlush() {
		return HandRank(Flushes[q])
	}
	if s := Unique5[q]; s != 0 {
		return HandRank(s)
	}

	fmt.Println("here")
	fmt.Printf("hand.Cards: %v\n", hand.Cards)
	return HandRank(HashValues[hashIndex(hand.getPrime())])
}

func hashIndex(prime uint) uint {
	var a, b uint
	prime += 0xe91aaa35
	prime ^= prime >> 16
	prime += prime << 8
	prime ^= prime >> 4
	b = (prime >> 8) & 0x1ff
	a = (prime + (prime << 2)) >> 19
	fmt.Printf("a: %v\n", a)
	fmt.Printf("b: %v\n", b)
	return a ^ uint(HashAdjust[b])
}

func (hand *Hand) isFlush() bool {
	res := hand.Cards[0]
	for _, card := range hand.Cards[1:] {
		res &= card
	}
	return res&suitMask != 0
}

func (hand *Hand) SortHand() {
	sortedCards := make([]Card, len(hand.Cards))
	copy(sortedCards, hand.Cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i]&rankMask < sortedCards[j]&rankMask
	})
	hand.Cards = sortedCards
}
