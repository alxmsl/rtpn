package payload

import (
	"fmt"
	"reflect"

	"github.com/alxmsl/cpn"
)

// IntStorage interface describes an object which stores integer value.
type IntStorage interface {
	SetValue(value int)
	Value() int
}

// IntPayload struct describes object to store integer value.
type IntPayload struct {
	value int
}

// SetValue sets integer value.
func (p *IntPayload) SetValue(value int) {
	p.value = value
}

// Value returns a previously stored integer value.
func (p *IntPayload) Value() int {
	return p.value
}

// IntStoragePayload returns a token payload as an IntStorage value.
func IntStoragePayload(token cpn.Token) (IntStorage, error) {
	var payload, ok = token.Payload().(IntStorage)
	if !ok {
		return nil, fmt.Errorf("payload type %s: %w", reflect.TypeOf(token.Payload()), cpn.ErrorWrongPayloadType)
	}
	return payload, nil
}
