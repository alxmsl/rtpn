package cpn_test

import (
	. "gopkg.in/check.v1"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/payload"
)

type TokenSuite struct{}

var _ = Suite(&TokenSuite{})

// TestPayloadInt tests PayloadInt function.
func (s *TokenSuite) TestPayloadInt(c *C) {
	var testsData = []struct {
		token         cpn.Token
		expectedValue int
	}{
		{
			token:         *cpn.NewToken("1"),
			expectedValue: 1,
		},
	}
	for _, testData := range testsData {
		var value, err = payload.PayloadIntString(testData.token)
		c.Assert(err, IsNil, Commentf("token: %v", testData.token))
		c.Assert(value, Equals, testData.expectedValue, Commentf("token: %v", testData.token))
	}
}
