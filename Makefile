.PHONY: bench
bench:
	go test  . -benchmem -bench BenchmarkDecode -count 1 -run _
	go test  . -benchmem -bench BenchmarkEncode -count 1 -run _
	go test  . -benchmem -bench BenchmarkDig -count 1 -run _
