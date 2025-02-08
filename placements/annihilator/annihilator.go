// Package annihilator
//
// An annihilator is a cpn.ProcessFunc which receives cpn.Token from the input channel and has neven sent cpn.Token to
// the output channel. So, annihilator removes token from the Petri-Net.
package annihilator

import (
	"context"

	"github.com/alxmsl/cpn"
)

// Annihilator defines a default annihilator. It just reads cpn.Token from an input channel and does nothing.
var Annihilator = func(ctx context.Context, input <-chan cpn.Token, _ chan<- cpn.Token) error {
	for {
		select {
		case <-input:
		case <-ctx.Done():
			return nil
		}
	}
}
