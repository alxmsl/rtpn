package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/placements/annihilator"
	"github.com/alxmsl/cpn/placements/generator"
	"github.com/alxmsl/prmtvs/plexus"
)

func main() {
	pttp()
}

//	digraph pttp {
//		p0[label="p0"]
//		p1[label="p1"]
//
//		t0[shape="box"]
//		t1[shape="box"]
//
//		p0 -> t0
//		p0 -> t1
//		t0 -> p1
//		t1 -> p1
//	}
func pttp() {
	var valueBuilder = func(idx int) any { return plexus.Counter(idx) }

	net := cpn.NewNet("pttp")
	_ = net.AddPlace("p0", generator.For(0, 10, 1, valueBuilder))
	_ = net.AddPlace("p1", annihilator.PrintToken("p1"))
	_ = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
	_ = net.AddTransition("t1", []string{"p0"}, []string{"p1"})

	var ctx, _ = context.WithTimeout(context.Background(), time.Second)
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
