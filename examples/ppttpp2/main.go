package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/placements/annihilator"
	"github.com/alxmsl/cpn/placements/generator"
	"github.com/alxmsl/cpn/placements/mediator"
	"github.com/alxmsl/prmtvs/plexus"
)

func main() {
	ppttpp2()
}

//	digraph ppttpp2 {
//		p0[label="p0"]
//		p1[label="p1"]
//		p2[label="p2"]
//		p3[label="p3"]
//
//		t0[shape="box"]
//		t1[shape="box"]
//
//		p0 -> t0
//		p0 -> t1
//		p1 -> t0
//		p1 -> t1
//		t0 -> p2
//		t1 -> p2
//		t0 -> p3
//		t1 -> p3
//	}
func ppttpp2() {
	var valueBuilder = func(idx int) any { return plexus.Counter(idx) }

	net := cpn.NewNet("ppttpp2")
	_ = net.AddPlace("p0", generator.For(0, 5, 1, valueBuilder))
	_ = net.AddPlace("p1", mediator.Mediator)
	_ = net.AddPlace("p2", annihilator.PrintToken("p2"))
	_ = net.AddPlace("p3", annihilator.PrintToken("p3"))
	_ = net.AddTransition("t0", []string{"p0", "p1"}, []string{"p2", "p3"})
	_ = net.AddTransition("t1", []string{"p0", "p1"}, []string{"p2", "p3"})

	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
	go func() {
		var t = cpn.NewToken(plexus.Counter(1))
		_ = net.PutToken("p1", *t)
		_ = net.PutToken("p1", *t)
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
