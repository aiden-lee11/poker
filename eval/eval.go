package eval

import (
	"fmt"
	"sort"
)

// Card will be a byte with vals | cdhsrrrr |
type Card byte
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
		// } else if hand.Strength == resHand.Strength {
		// 	resHand = hand.EvaluateEqualStrength(resHand)
		// }
	}

	return resHand
}

func (hand *Hand) sortHand() {
	sort.Slice(hand.Cards, func(i, j int) bool {
		return hand.Cards[i]&0x0F < hand.Cards[j]&0x0F
	})
}

// func (hand *Hand) EvaluateEqualStrength(other Hand) Hand {
// 	switch hand.Strength {
// 	case HighCard:
// 		for i := len(hand.Cards); i >= 0; i-- {
// 			if hand.Cards[i] > other.Cards[i] {
// 				return *hand
// 			} else if other.Cards[i] > hand.Cards[i] {
// 				return other
// 			}
// 		}
// 	case OnePair:

// 	}
// }

func (hand *Hand) EvaluateHand() {
	if hand.isFlush() {
		if !hand.isStraight() {
			hand.Strength = Flush
			return
		}

		if hand.Cards[0]&0x09 != 0 {
			hand.Strength = RoyalFlush
			return
		}
		hand.Strength = StraightFlush
		return
	}

	if hand.isNumOfAKind(4) {
		hand.Strength = Quads
		return
	}

	if hand.isFullHouse() {
		hand.Strength = FullHouse
		return
	}

	if hand.isStraight() {
		hand.Strength = Straight
		return
	}

	if hand.isNumOfAKind(3) {
		hand.Strength = Trips
		return
	}

	if hand.isTwoPair() {
		hand.Strength = TwoPair
		return
	}

	if hand.isNumOfAKind(2) {
		hand.Strength = OnePair
		return
	}

	hand.Strength = HighCard
	return
}

func (hand *Hand) isFlush() bool {
	res := hand.Cards[0]
	for _, card := range hand.Cards[1:] {
		res &= card
	}
	return res&0xF0 != 0
}

func (hand *Hand) isStraight() bool {
	for i := 0; i < len(hand.Cards)-1; i++ {
		if hand.Cards[i] != hand.Cards[i+1]-1 {
			return false
		}
	}
	return true
}

// Need to fully calculate hand first to ensure that trips dont count as 2 of a kind
func (hand *Hand) isNumOfAKind(num int) bool {

	rankCount := make(map[Card]int)

	for _, card := range hand.Cards {
		rank := card & 0x0F
		rankCount[rank]++
	}

	for _, card := range hand.Cards {
		rank := card & 0x0F
		if rankCount[rank] == num {
			return true
		}
	}

	return false
}

func (hand *Hand) isFullHouse() bool {
	return hand.isNumOfAKind(3) && hand.isNumOfAKind(2)
}

// Ensure that the two pairs are distinct and not the same ie 77 not double counted would need 77 and 88 etc
func (hand *Hand) isTwoPair() bool {
	rankCount := make(map[Card]int)

	for _, card := range hand.Cards {
		rank := card & 0x0F
		rankCount[rank]++
	}

	res := 0
	for _, card := range hand.Cards {
		rank := card & 0x0F
		if rankCount[rank] == 2 {
			res += 1
			rankCount[rank] = 0
		}
	}

	return res == 2
}
