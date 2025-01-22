package annihilator

import (
	"context"
	"fmt"

	"github.com/alxmsl/cpn"
)

// PrintToken creates an annihilator which reads cpn.Token from the input channel and calls fmt.Println for it in
// concatenation for a given arguments.
func PrintToken(a ...any) cpn.ProcessFunc {
	return func(ctx context.Context, input <-chan cpn.Token, _ chan<- cpn.Token) error {
		for {
			select {
			case token := <-input:
				var a = append(a, token.String())
				fmt.Println(a...)
			case <-ctx.Done():
				return nil
			}
		}
	}
}
