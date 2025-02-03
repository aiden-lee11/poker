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
type Card int32
type HandStrength int

const (
	HighCard HandStrength = iota
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

type Hand struct {
	Cards    []Card
	Strength HandStrength
}

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
			hand.sortHand()
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

func findBestHand(combinations []Hand) Hand {
	var resHand Hand

	for _, hand := range combinations {
		hand.EvaluateHand()
		if hand.Strength > resHand.Strength {
			resHand = hand
		}
	}

	return resHand
}

func (hand *Hand) sortHand() {
	sort.Slice(hand.Cards, func(i, j int) bool {
		return hand.Cards[i]&0x00000F00 < hand.Cards[j]&0x00000F00
	})
}

// func (hand *Hand) EvaluateEqualStrength(other Hand) Hand {
// }

func (hand *Hand) getFlushStraightIndex() int16 {
	return int16((hand.Cards[0] | hand.Cards[1] | hand.Cards[2] | hand.Cards[3] | hand.Cards[4]) >> 16)
}

func (hand *Hand) getPrime() uint {
	return uint((hand.Cards[0] & 0xFF) * (hand.Cards[1] & 0xFF) * (hand.Cards[2] & 0xFF) * (hand.Cards[3] & 0xFF) * (hand.Cards[4] & 0xFF))
}

func (hand *Hand) EvaluateHand() int16 {
	q := hand.getFlushStraightIndex()

	if hand.isFlush() {
		return Flushes[q]
	}
	if s := Unique5[q]; s != 0 {
		return s
	}

	return HashValues[hashIndex(hand.getPrime())]
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
	return res&0x0000F000 != 0
}
