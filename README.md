# Concurrent Petri Net

Petri Net implementation based on golang concurrency patterns such as goroutines and channels

## Benchmark

Solution overhead is about 3-5μs per transition

```
BenchmarkPTP-12        	  760107	      1362 ns/op	     200 B/op	       4 allocs/op
BenchmarkPTPTP-12      	  540799	      2213 ns/op	     432 B/op	       8 allocs/op
BenchmarkPTPTPTP-12    	  457405	      2478 ns/op	     744 B/op	      12 allocs/op
BenchmarkPPTTPP-12     	  461224	      2435 ns/op	     552 B/op	      10 allocs/op
```
