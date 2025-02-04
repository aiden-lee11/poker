package eval

import (
	"fmt"
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

// Takes in 7 cards and returns all possible 5 card hands, essentially 7 chose 5
// Thanks cesar :)
func GenerateCombinations(cards []Card) []Hand {
	if len(cards) != 7 {
		fmt.Println("Expected 7 cards in the combo generator got: ", len(cards))
		return nil
	}

	var combos []Hand

	var helper func(start int, combo []Card)
	helper = func(start int, combo []Card) {
		if len(combo) == 5 {
			hand := Hand{Cards: combo}
			hand.SortHand()
			combos = append(combos, hand)
			return
		}

		for i := start; i < 7; i++ {
			helper(i+1, append(combo, cards[i]))
		}
	}

	helper(0, []Card{})
	return combos
}

func FindBestHand(combinations []Hand) Hand {
	resHand := combinations[0]

	for i := range combinations[1:] {
		hand := combinations[i]
		if curStrength := hand.EvaluateHand(); curStrength > resHand.EvaluateHand() {
			resHand = hand
		}
	}

	return resHand
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

func (hand *Hand) SortHand() {
	sort.Slice(hand.Cards, func(i, j int) bool {
		return hand.Cards[i]&rankMask < hand.Cards[j]&rankMask
	})
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
	return a ^ uint(HashAdjust[b])
}

func (hand *Hand) isFlush() bool {
	res := hand.Cards[0]
	for _, card := range hand.Cards[1:] {
		res &= card
	}
	return res&suitMask != 0
}
