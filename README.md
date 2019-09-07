# Fair benchmarks for golang JSON libs 
```
BenchmarkFair/stable-flavor|insane-json-4      1000	   2081607 ns/op	 617.96 MB/s	    1664 B/op	       8 allocs/op
BenchmarkFair/stable-flavor|fastjson-4         500	   2445162 ns/op	 526.08 MB/s	   20839 B/op	      20 allocs/op
BenchmarkFair/stable-flavor|jsonparser-4       50	  33479704 ns/op	  38.42 MB/s	      15 B/op	       0 allocs/op
BenchmarkFair/stable-flavor|gjson-4            50	  23755202 ns/op	  54.15 MB/s	   21223 B/op	     181 allocs/op
BenchmarkFair/stable-flavor|go-simplejson-4    30	  50403348 ns/op	  25.52 MB/s	14382211 B/op	  134479 allocs/op
BenchmarkFair/chaotic-flavor|insane-json-4     20	  72524415 ns/op	   7.56 MB/s	       0 B/op	       0 allocs/op
BenchmarkFair/chaotic-flavor|fastjson-4        20	  68833467 ns/op	   7.96 MB/s	    1401 B/op	       2 allocs/op
BenchmarkFair/chaotic-flavor|jsonparser-4      2	 933643754 ns/op	   0.59 MB/s	       0 B/op	       0 allocs/op
BenchmarkFair/chaotic-flavor|gjson-4           1	1195911707 ns/op	   0.46 MB/s	195482064 B/op	 3259151 allocs/op
BenchmarkFair/chaotic-flavor|go-simplejson-4   5	 249144470 ns/op	   2.20 MB/s	44908977 B/op	 2535674 allocs/op
```