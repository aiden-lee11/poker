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

func TestHands(t *testing.T) {

	hand1 := eval.Hand{Cards: []eval.Card{0x80, 0x80, 0x80, 0x80, 0x80}}

	res := hand1.EvaluateHand()

	assert.Equal(t, res, eval.Flush, "Expected flush got %v", res)
}
