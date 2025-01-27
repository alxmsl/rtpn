package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/placements/annihilator"
	"github.com/alxmsl/cpn/placements/mediator"
	"github.com/alxmsl/prmtvs/plexus"
)

func main() {
	pptpp()
}

// pptpp runs a Petri-Net for the graph:
//
//	digraph pptpp {
//		p0[label="p0"]
//		p1[label="p1"]
//		p2[label="p2"]
//		p3[label="p3"]
//
//		t0[shape="box"]
//
//		p0 -> t0
//		p1 -> t0
//		t0 -> p2
//		t0 -> p3
//	}
func pptpp() {
	net := cpn.NewNet("pptpp")
	_ = net.AddPlace("p0", mediator.Mediator)
	_ = net.AddPlace("p1", mediator.Mediator)
	_ = net.AddPlace("p2", annihilator.PrintToken("p2"))
	_ = net.AddPlace("p3", annihilator.PrintToken("p3"))
	_ = net.AddTransition("t0", []string{"p0", "p1"}, []string{"p2", "p3"})

	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	go func() {
		var (
			t   = cpn.NewToken(plexus.Counter(1))
			err = net.Markup("p0", *t)
		)
		if err != nil {
			panic(err)
		}

		t = cpn.NewToken(plexus.Counter(2))
		err = net.Markup("p1", *t)
		if err != nil {
			panic(err)
		}
	}()
	var err = net.Run(ctx)
	if err != nil {
		panic(err)
	}

	net.Stats().Over(func(idx int, placeName string, v interface{}) bool {
		var statsValue = v.(cpn.Stats)
		fmt.Println(placeName, statsValue)
		return true
	})
}
