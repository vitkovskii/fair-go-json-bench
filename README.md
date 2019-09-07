# Fair benchmarks for golang JSON libs 
```
stable-flavor|insane-json-4     1000	   2031230 ns/op	 633.29 MB/s
stable-flavor|fastjson-4        500	       2308785 ns/op	 557.15 MB/s
stable-flavor|jsonparser-4      50	       32628135 ns/op	 39.42 MB/s
stable-flavor|gjson-4           50	       23047809 ns/op	 55.81 MB/s
stable-flavor|go-simplejson-4   30	       47524008 ns/op	 27.07 MB/s
chaotic-flavor|insane-json-4    3000	   521958 ns/op	     159.21 MB/s
chaotic-flavor|fastjson-4       2000	   586540 ns/op	     141.68 MB/s
chaotic-flavor|jsonparser-4     100	       12841475 ns/op	 6.47 MB/s
chaotic-flavor|gjson-4          100	       14691591 ns/op	 5.66 MB/s
chaotic-flavor|go-simplejson-4  200	       5597591 ns/op	 14.85 MB/s
```