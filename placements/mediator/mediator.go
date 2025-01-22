// Package mediator
//
// A mediator is a cpn.ProcessFunc which receives cpn.Token from the input channel and sends cpn.Token to the output
// channel. So, mediator just passes token forward in the Petri-Net.
package mediator

import (
	"context"

	"github.com/alxmsl/cpn"
)

// Mediator defines a default mediator. It just reads cpn.Token from an input channel and sends the same cpn.Token to
// the output channel.
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
