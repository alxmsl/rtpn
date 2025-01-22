package mediator

import (
	"context"
	"fmt"

	"github.com/alxmsl/cpn"
)

// PrintToken creates a mediator which reads cpn.Token from the input channel and calls fmt.Println for it in
// concatenation for a given arguments. Then it sends the same cpn.Token to the output channel.
func PrintToken(a ...any) cpn.ProcessFunc {
	return func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
		for {
			select {
			case token := <-input:
				var a = append(a, token.String())
				fmt.Println(a...)
				output <- token
			case <-ctx.Done():
				return nil
			}
		}
	}
}
