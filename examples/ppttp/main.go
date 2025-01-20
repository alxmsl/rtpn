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
	ppttp()
}

//	digraph ppttp {
//		p0[label="p0"]
//		p1[label="p1"]
//		p2[label="p2"]
//
//		t0[shape="box"]
//		t1[shape="box"]
//
//		p0 -> t0
//		p1 -> t1
//		t0 -> p2
//		t1 -> p2
//	}
func ppttp() {
	var valueBuilder = func(idx int) any { return plexus.Counter(idx) }

	net := cpn.NewNet("ppttp")
	_ = net.AddPlace("p0", generator.For(1, 10, 2, valueBuilder))
	_ = net.AddPlace("p1", generator.For(0, 10, 2, valueBuilder))
	_ = net.AddPlace("p2", annihilator.PrintToken("p2"))
	_ = net.AddTransition("t0", []string{"p0"}, []string{"p2"})
	_ = net.AddTransition("t1", []string{"p1"}, []string{"p2"})

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
