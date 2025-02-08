package cpn_test

import (
	. "gopkg.in/check.v1"

	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/placements/annihilator"
	"github.com/alxmsl/cpn/placements/mediator"
	"github.com/alxmsl/prmtvs/plexus"
)

func Test(t *testing.T) {
	TestingT(t)
}

type CpnSuite struct{}

var _ = Suite(&CpnSuite{})

const uintZero = uint64(0)

var testsData = []plexus.Counter{
	plexus.Counter(0),
	plexus.Counter(1),
	plexus.Counter(2),
}

var (
	testDataGenerator = func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
		for _, testData := range testsData {
			t := cpn.NewToken(testData)
			output <- *t
		}
		return nil
	}
	testDataAnnihilator = func(ctx context.Context, input <-chan cpn.Token, output chan<- cpn.Token) error {
		for _, testData := range testsData {
			token := <-input
			if token.Payload() != testData {
				return fmt.Errorf("expected %v, got %v", testData, token.Payload())
			}
		}
		return nil
	}
)

func stats(c *C, net *cpn.Net, placeName string) cpn.Stats {
	var v, ok = net.Stats().GetByKey(placeName)
	c.Assert(ok, Equals, true)
	statsValue, ok := v.(cpn.Stats)
	c.Assert(ok, Equals, true)
	return statsValue
}

// TestPTP runs a Net which is described as the graph:
//
//	digraph TestPTP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//
//	   t0[shape="box"]
//
//	   p0 -> t0
//	   t0 -> p1
//	}
func (s *CpnSuite) TestPTP(c *C) {
	var builder = func(c *C, p0, p1 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPTP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, testDataAnnihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)

	c.Assert(stats(c, net, "p1").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Sent, Equals, uintZero)
}

// TestPTTP runs a Net which is described as the graph:
//
//	digraph TestPTTP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//
//	   t0[shape="box"]
//	   t1[shape="box"]
//
//	   p0 -> t0
//	   p0 -> t1
//	   t0 -> p1
//	   t1 -> p1
//	}
func (s *CpnSuite) TestPTTP(c *C) {
	var builder = func(c *C, p0, p1 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPTTP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
		c.Assert(err, IsNil)
		err = net.AddTransition("t1", []string{"p0"}, []string{"p1"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, annihilator.Annihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)

	c.Assert(stats(c, net, "p1").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Sent, Equals, uintZero)
}

// TestPTPP runs a Net which is described as the graph:
//
//	digraph TestPTPP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//	   p2[label="p2"]
//
//	   t0[shape="box"]
//
//	   p0 -> t0
//	   t0 -> p1
//	   t0 -> p2
//	}
func (s *CpnSuite) TestPTPP(c *C) {
	var builder = func(c *C, p0, p1, p2 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPTPP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddPlace("p2", p2)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0"}, []string{"p1", "p2"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, testDataAnnihilator, testDataAnnihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)

	// Both placements get the same amount of testsData.
	c.Assert(stats(c, net, "p1").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Sent, Equals, uintZero)
	c.Assert(stats(c, net, "p2").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Sent, Equals, uintZero)
}

// TestPPTP runs a Net which is described as the graph:
//
//	digraph TestPPTP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//	   p2[label="p2"]
//
//	   t0[shape="box"]
//
//	   p0 -> t0
//	   p1 -> t0
//	   t0 -> p2
//	}
func (s *CpnSuite) TestPPTP(c *C) {
	var builder = func(c *C, p0, p1, p2 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPPTP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddPlace("p2", p2)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0", "p1"}, []string{"p2"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, testDataGenerator, annihilator.Annihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))

	// Both placements sent the same amount of testsData.
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p1").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p1").Sent, Equals, expectedValue)

	// Receiver placement gets the same amount of testsData. Transition passes value only when all placements are ready.
	c.Assert(stats(c, net, "p2").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Sent, Equals, uintZero)
}

// TestPPTTP runs a Net which is described as the graph:
//
//	digraph TestPPTTP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//	   p2[label="p2"]
//
//	   t0[shape="box"]
//	   t1[shape="box"]
//
//	   p0 -> t0
//	   p1 -> t1
//	   t0 -> p2
//	   t1 -> p2
//	}
func (s *CpnSuite) TestPPTTP(c *C) {
	var builder = func(c *C, p0, p1, p2 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPPTTP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddPlace("p2", p2)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0"}, []string{"p2"})
		c.Assert(err, IsNil)
		err = net.AddTransition("t1", []string{"p1"}, []string{"p2"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, testDataGenerator, annihilator.Annihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))

	// Both placements sent the same amount of testsData.
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p1").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p1").Sent, Equals, expectedValue)

	// Receiver placement gets twice more amount of testsData.
	c.Assert(stats(c, net, "p2").Accepted, Equals, 2*expectedValue)
	c.Assert(stats(c, net, "p2").Processing, Equals, 2*expectedValue)
	c.Assert(stats(c, net, "p2").Sent, Equals, uintZero)
}

// TestPPTTPP runs a Net which is described as the graph:
//
//	digraph TestPPTTPP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//	   p2[label="p2"]
//	   p3[label="p3"]
//
//	   t0[shape="box"]
//	   t1[shape="box"]
//
//	   p0 -> t0
//	   p0 -> t1
//	   p1 -> t0
//	   p1 -> t1
//	   t0 -> p2
//	   t0 -> p3
//	   t1 -> p2
//	   t1 -> p3
//	}
func (s *CpnSuite) TestPPTTPP(c *C) {
	var builder = func(c *C, p0, p1, p2, p3 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPPTTPP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddPlace("p2", p2)
		c.Assert(err, IsNil)
		err = net.AddPlace("p3", p3)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0", "p1"}, []string{"p2", "p3"})
		c.Assert(err, IsNil)
		err = net.AddTransition("t1", []string{"p0", "p1"}, []string{"p2", "p3"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, testDataGenerator, annihilator.Annihilator, annihilator.Annihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))

	// Both placements sent the same amount of testsData.
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p1").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p1").Sent, Equals, expectedValue)

	// Both receiver placements get the same amount of testsData.
	// Transition passes value only when all placements are ready.
	c.Assert(stats(c, net, "p2").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Sent, Equals, uintZero)
	c.Assert(stats(c, net, "p3").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p3").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p3").Sent, Equals, uintZero)
}

// TestPTPTP runs a Net which is described as the graph:
//
//	digraph TestPTPTP {
//	   p0[label="p0"]
//	   p1[label="p1"]
//	   p2[label="p2"]
//
//	   t0[shape="box"]
//	   t1[shape="box"]
//
//	   p0 -> t0
//	   t0 -> p1
//	   p1 -> t1
//	   t1 -> p2
//	}
func (s *CpnSuite) TestPTPTP(c *C) {
	var builder = func(c *C, p0, p1, p2 cpn.ProcessFunc) *cpn.Net {
		var net = cpn.NewNet("TestPTPTP")
		var err error

		err = net.AddPlace("p0", p0)
		c.Assert(err, IsNil)
		err = net.AddPlace("p1", p1)
		c.Assert(err, IsNil)
		err = net.AddPlace("p2", p2)
		c.Assert(err, IsNil)
		err = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
		err = net.AddTransition("t1", []string{"p1"}, []string{"p2"})
		c.Assert(err, IsNil)
		return net
	}

	var net = builder(c, testDataGenerator, mediator.Mediator, testDataAnnihilator)
	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	var err = net.Run(ctx)
	c.Assert(err, IsNil)

	var expectedValue = uint64(len(testsData))
	c.Assert(stats(c, net, "p0").Accepted, Equals, uintZero)
	c.Assert(stats(c, net, "p0").Processing, DeepEquals, uintZero)
	c.Assert(stats(c, net, "p0").Sent, Equals, expectedValue)

	c.Assert(stats(c, net, "p1").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p1").Sent, Equals, expectedValue)

	c.Assert(stats(c, net, "p2").Accepted, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Processing, Equals, expectedValue)
	c.Assert(stats(c, net, "p2").Sent, Equals, uintZero)
}
