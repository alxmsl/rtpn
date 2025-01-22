// Package generator
//
// A generator is a cpn.ProcessFunc which sends cpn.Token to the output channel and has neven reads cpn.Token from
// the input channel. So, generator creates token for the Petri-Net.
package generator

import (
	"context"

	"github.com/alxmsl/cpn"
)

type ValueBuilder func(idx int) any

// For creates a generator which sends cpn.Token for a given range of integer values. It uses a given ValueBuilder to
// get payload for a cpn.Token.
func For(from, to, step int, valueBuilder ValueBuilder) cpn.ProcessFunc {
	return func(ctx context.Context, _ <-chan cpn.Token, output chan<- cpn.Token) error {
		for i := from; i < to; i += step {
			t := cpn.NewToken(valueBuilder(i))
			output <- *t
		}
		return nil
	}
}
