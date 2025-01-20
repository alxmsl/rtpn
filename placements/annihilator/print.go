package annihilator

import (
	"context"
	"fmt"

	"github.com/alxmsl/cpn"
)

func PrintToken(a ...any) cpn.ProcessFunc {
	return func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
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
