package payload

import (
	"fmt"
	"reflect"

	"github.com/alxmsl/cpn"
)

// Erroneous interface describes an object which stores error.
type Erroneous interface {
	AddError(err error)
	LastError() error
}

// ErrorListPayload
type ErrorListPayload struct {
	errors []error
}

// AddError stores a given error to the object.
func (p *ErrorListPayload) AddError(err error) {
	p.errors = append(p.errors, err)
}

// LastError returns a previously stored error from the object. If there is no error, then function returns nil.
func (p *ErrorListPayload) LastError() error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors[len(p.errors)-1]
}

// ErroneousPayload returns a token payload as an Erroneous value.
func ErroneousPayload(token cpn.Token) (Erroneous, error) {
	var payload, ok = token.Payload().(Erroneous)
	if !ok {
		return nil, fmt.Errorf("payload type %s: %w", reflect.TypeOf(token.Payload()), cpn.ErrorWrongPayloadType)
	}
	return payload, nil
}
