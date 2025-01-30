package eval

import (
	"fmt"
)

// Card will be a byte with vals | cdhsrrrr |
type Card byte

// func for generating combos, evaluating a single 5 card hand strength, evaluating best 5 card hand out 7

// Takes in 7 cards and returns all possible 5 card hands, essentially 7 chose 5
func GenerateCombinations(cards []Card) [][]Card {
	if len(cards) != 7 {
		fmt.Println("Expected 7 cards in the combo generator got: ", len(cards))
		return nil
	}

	var combos [][]Card

	var helper func(start int, combo []Card)
	helper = func(start int, combo []Card) {
		if len(combo) == 5 {
			comboCopy := make([]Card, 5)
			copy(comboCopy, combo)
			combos = append(combos, comboCopy)
			return
		}

		for i := start; i < 7; i++ {
			helper(i+1, append(combo, cards[i]))
		}
	}

	helper(0, []Card{})
	return combos
}
