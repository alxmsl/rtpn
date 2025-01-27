package main

import (
	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/placements/annihilator"
	"github.com/alxmsl/cpn/placements/mediator"
)

// Net defines an interop symbol to be used externally.
var Net = *newNet()

// newNet construct a Petri-Net for the graph:
//
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
func newNet() *cpn.Net {
	var err error
	net := cpn.NewNet("newNet")
	err = net.AddPlace("p0", mediator.Mediator)
	if err != nil {
		panic(err)
	}
	err = net.AddPlace("p1", annihilator.PrintToken("p1"))
	if err != nil {
		panic(err)
	}
	err = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
	if err != nil {
		panic(err)
	}
	err = net.AddTransition("t1", []string{"p0"}, []string{"p1"})
	if err != nil {
		panic(err)
	}
	return net
}
