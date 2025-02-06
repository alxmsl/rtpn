package cpn

import (
	"fmt"
	"strings"
	"time"

	"github.com/alxmsl/prmtvs/plexus"
)

// Token struct defines a Petri-Net token
type Token struct {
	payload any
	word    checkpoints
}

// NewToken creates a Token with a given payload.
func NewToken(payload any) *Token {
	return &Token{
		payload: payload,
	}
}

// Payload returns current payload for the Token
func (t Token) Payload() any {
	return t.payload
}

func (t Token) String() string {
	return fmt.Sprintf("%v: %s", t.payload, t.word)
}

// Merge implements a plexus.Mergeable interface for a Token.
// Payload has to implement plexus.Mergeable interface as well. Otherwise, function panics.
func (a Token) Merge(b plexus.Mergeable) plexus.Mergeable {
	if _, ok := a.payload.(plexus.Mergeable); !ok {
		panic(fmt.Errorf("payload from a is not mergeable: %w", ErrorWrongTokenType))
	}
	if _, ok := b.(Token).payload.(plexus.Mergeable); !ok {
		panic(fmt.Errorf("payload from b is not mergeable: %w", ErrorWrongTokenType))
	}
	var (
		v1 = a.payload.(plexus.Mergeable)
		v2 = b.(Token).payload.(plexus.Mergeable)
	)
	a.payload = v1.Merge(v2)
	return a
}

type checkpoint struct {
	name string
	when time.Time
}

func (cp checkpoint) String() string {
	return fmt.Sprintf("%s[%s]", cp.name, cp.when.Format(time.RFC3339Nano))
}

type checkpoints []checkpoint

func (cps checkpoints) String() string {
	var result = make([]string, 0, len(cps))
	for _, cp := range cps {
		result = append(result, cp.String())
	}
	return strings.Join(result, ",")
}

func (cps checkpoints) Word() string {
	var ww = make([]string, 0, len(cps))
	for _, cp := range cps {
		ww = append(ww, cp.name)
	}
	return strings.Join(ww, ",")
}
