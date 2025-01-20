package annihilator

import (
	"context"

	"github.com/alxmsl/cpn"
)

var Annihilator = func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
	for {
		select {
		case <-input:
		case <-ctx.Done():
			return nil
		}
	}
}
