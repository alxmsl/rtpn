package cpn_test

import (
	"context"
	"testing"

	"github.com/alxmsl/cpn"
	"github.com/alxmsl/cpn/placements/annihilator"
	"github.com/alxmsl/cpn/placements/mediator"
	"github.com/alxmsl/prmtvs/plexus"
)

//	digraph ptp {
//		p0[label="p0"]
//		p1[label="p1"]
//
//		t0[shape="box"]
//
//		p0 -> t0
//		t0 -> p1
//	}
func BenchmarkPTP(b *testing.B) {
	net := cpn.NewNet("ptp")
	_ = net.AddPlace("p0", mediator.Mediator)
	_ = net.AddPlace("p1", annihilator.Annihilator)
	_ = net.AddTransition("t0", []string{"p0"}, []string{"p1"})

	var ctx, cancel = context.WithCancel(context.Background())
	go func() {
		var err = net.Run(ctx)
		if err != nil {
			b.Error(err)
			return
		}
	}()

	var token = cpn.NewToken("some payload")
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		var err = net.PutToken("p0", *token)
		if err != nil {
			b.Fatal(err)
			return
		}
	}
	cancel()
}

//	digraph ptptp {
//		p0[label="p0"]
//		p1[label="p1"]
//		p2[label="p2"]
//
//		t0[shape="box"]
//		t1[shape="box"]
//
//		p0 -> t0
//		t0 -> p1
//		p1 -> t1
//		t1 -> p2
//	}
func BenchmarkPTPTP(b *testing.B) {
	net := cpn.NewNet("ptptp")
	_ = net.AddPlace("p0", mediator.Mediator)
	_ = net.AddPlace("p1", mediator.Mediator)
	_ = net.AddPlace("p2", annihilator.Annihilator)
	_ = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
	_ = net.AddTransition("t1", []string{"p1"}, []string{"p2"})

	var ctx, cancel = context.WithCancel(context.Background())
	go func() {
		var err = net.Run(ctx)
		if err != nil {
			b.Error(err)
			return
		}
	}()

	var token = cpn.NewToken("some payload")
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		var err = net.PutToken("p0", *token)
		if err != nil {
			b.Fatal(err)
			return
		}
	}
	cancel()
}

//	digraph ptptptp {
//		p0[label="p0"]
//		p1[label="p1"]
//		p2[label="p2"]
//		p3[label="p3"]
//
//		t0[shape="box"]
//		t1[shape="box"]
//		t2[shape="box"]
//
//		p0 -> t0
//		t0 -> p1
//		p1 -> t1
//		t1 -> p2
//		p2 -> t2
//		t2 -> p3
//	}
func BenchmarkPTPTPTP(b *testing.B) {
	net := cpn.NewNet("ptptptp")
	_ = net.AddPlace("p0", mediator.Mediator)
	_ = net.AddPlace("p1", mediator.Mediator)
	_ = net.AddPlace("p2", mediator.Mediator)
	_ = net.AddPlace("p3", annihilator.Annihilator)
	_ = net.AddTransition("t0", []string{"p0"}, []string{"p1"})
	_ = net.AddTransition("t1", []string{"p1"}, []string{"p2"})
	_ = net.AddTransition("t2", []string{"p2"}, []string{"p3"})

	var ctx, cancel = context.WithCancel(context.Background())
	go func() {
		var err = net.Run(ctx)
		if err != nil {
			b.Error(err)
			return
		}
	}()

	var token = cpn.NewToken("some payload")
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		var err = net.PutToken("p0", *token)
		if err != nil {
			b.Fatal(err)
			return
		}
	}
	cancel()
}

//	digraph ppttpp {
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
func BenchmarkPPTTPP(b *testing.B) {
	net := cpn.NewNet("ppttpp")
	_ = net.AddPlace("p0", mediator.Mediator)
	_ = net.AddPlace("p1", mediator.Mediator)
	_ = net.AddPlace("p2", annihilator.Annihilator)
	_ = net.AddPlace("p3", annihilator.Annihilator)
	_ = net.AddTransition("t0", []string{"p0", "p1"}, []string{"p2", "p3"})
	_ = net.AddTransition("t1", []string{"p0", "p1"}, []string{"p2", "p3"})

	var ctx, cancel = context.WithCancel(context.Background())
	go func() {
		var err = net.Run(ctx)
		if err != nil {
			b.Error(err)
			return
		}
	}()

	var token = cpn.NewToken(plexus.Counter(1))
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		var err error
		err = net.PutToken("p0", *token)
		if err != nil {
			b.Fatal(err)
			return
		}
		err = net.PutToken("p1", *token)
		if err != nil {
			b.Fatal(err)
			return
		}
	}
	cancel()
}
