package eval_test

import (
	"fmt"
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
