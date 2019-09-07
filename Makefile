.PHONY: bench
bench:
	go test  . -benchmem -bench Benchmark -count 1 -run _

