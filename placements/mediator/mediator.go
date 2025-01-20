package mediator

import (
	"context"

	"github.com/alxmsl/cpn"
)

var Mediator = func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
	for {
		select {
		case token := <-input:
			output <- token
		case <-ctx.Done():
			return nil
		}
	}
}
