package payload

import (
	"strconv"

	"github.com/alxmsl/cpn"
)

// PayloadIntString returns a token payload as an int value. Function expects that payload is string.
func PayloadIntString(token cpn.Token) (int, error) {
	var payload, ok = token.Payload().(string)
	if !ok {
		return 0, cpn.ErrorWrongPayloadType
	}
	value, err := strconv.Atoi(payload)
	if err != nil {
		return 0, err
	}
	return value, nil
}
