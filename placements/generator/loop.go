package generator

import (
	"context"

	"github.com/alxmsl/cpn"
)

func For(from, to, step int, valueBuilder func(idx int) any) cpn.ProcessFunc {
	return func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
		for i := from; i < to; i += step {
			t := cpn.NewToken(valueBuilder(i))
			output <- *t
		}
		return nil
	}
}
