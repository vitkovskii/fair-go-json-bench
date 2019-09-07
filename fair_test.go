package fair_go_json_bench

import (
	"bufio"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/tidwall/gjson"
	"github.com/valyala/fastjson"
	"github.com/vitkovskii/insane-json"
)

type workload struct {
	json []byte
	name string

	requests [][]string
}

func getStableWorkload() ([]*workload, int64) {
	workloads := make([]*workload, 0, 0)
	workloads = append(workloads, loadJSON("light-ws", [][]string{
		{"_id"},
		{"favoriteFruit"},
		{"about"},
	}))
	workloads = append(workloads, loadJSON("many-objects", [][]string{
		{"deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper", "deeper"},
	}))
	workloads = append(workloads, loadJSON("heavy", [][]string{
		{"first", "second", "third", "fourth", "fifth"},
	}))
	workloads = append(workloads, loadJSON("many-fields", [][]string{
		{"first"},
		{"middle"},
		{"last"},
	}))
	workloads = append(workloads, loadJSON("few-fields", [][]string{
		{"first"},
		{"middle"},
		{"last"},
	}))
	workloads = append(workloads, loadJSON("insane", [][]string{
		{"statuses", "2", "user", "entities", "url", "urls", "0", "expanded_url"},
		{"statuses", "36", "retweeted_status", "user", "profile", "sidebar", "fill", "color"},
		{"statuses", "75", "entities", "user_mentions", "0", "screen_name"},
		{"statuses", "99", "coordinates"},
	}))

	size := 0
	for _, workload := range workloads {
		size += len(workload.json)
	}

	return workloads, int64(size)
}

func loadJSON(name string, requests [][]string) *workload {
	content, err := ioutil.ReadFile(fmt.Sprintf("benchdata/%s.json", name))
	if err != nil {
		panic(err.Error())
	}

	return &workload{json: content, name: name, requests: requests}
}

func getChaoticWorkload() ([][]byte, [][][]string, int64) {
	lines := make([][]byte, 0, 0)
	requests := make([][][]string, 0, 0)
	file, err := os.Open("./benchdata/chaotic-workload.log")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = file.Close()
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := []byte(scanner.Text())
		lines = append(lines, bytes)
		root, err := insaneJSON.DecodeBytes(bytes)
		if err != nil {
			panic(err.Error())
		}

		requestList := make([][]string, 0, 0)
		requestCount := rand.Int() % 3
		for x := 0; x < requestCount; x++ {
			node := root.Node
			selector := make([]string, 0, 0)
			for {
				if node.Type != insaneJSON.Object {
					break
				}

				fields := node.AsFields()
				name := fields[rand.Int()%len(fields)].AsString()
				selector = append(selector, string([]byte(name)))

				node = node.Dig(name)
			}
			requestList = append(requestList, selector)
		}
		requests = append(requests, requestList)

		insaneJSON.Release(root)
	}
	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	s, _ := file.Stat()
	return lines, requests, s.Size()
}

// BenchmarkFair benchmarks overall performance of libs as fair as it can:
// * using various JSON payload
// * decoding
// * doing low and high count of search requests
// * encoding
func BenchmarkFair(b *testing.B) {

	// some big buffer to avoid allocations
	s := make([]byte, 0, 512*1024)

	// let's make it deterministic as hell
	rand.Seed(666)

	// do little and few amount of search request
	requestsCount := []int{1, 8}

	pretenders := []struct {
		name string
		fn   func(b *testing.B, jsons [][]byte, fields [][][]string, reqCount int)
	}{
		{
			name: "insane-json",
			fn: func(b *testing.B, jsons [][]byte, fields [][][]string, reqCount int) {
				root := insaneJSON.Spawn()
				for i := 0; i < b.N; i++ {
					for _, json := range jsons {
						_ = insaneJSON.DecodeBytesReusing(root, json)
						for j := 0; j < reqCount; j++ {
							for _, f := range fields {
								for _, ff := range f {
									root.Dig(ff...)
								}
							}
						}
						s = root.EncodeNoAlloc(s[:0])
					}
				}
				insaneJSON.Release(root)
			},
		},
		{
			name: "fastjson",
			fn: func(b *testing.B, jsons [][]byte, fields [][][]string, reqCount int) {
				parser := fastjson.Parser{}
				for i := 0; i < b.N; i++ {
					for _, json := range jsons {
						c, _ := parser.ParseBytes(json)
						for j := 0; j < reqCount; j++ {
							for _, f := range fields {
								for _, ff := range f {
									c.Get(ff...)
								}
							}
						}
						s = c.MarshalTo(s[:0])
					}
				}
			},
		},
		{
			name: "jsonparser",
			fn: func(b *testing.B, jsons [][]byte, fields [][][]string, reqCount int) {
				for i := 0; i < b.N; i++ {
					for _, json := range jsons {
						for j := 0; j < reqCount; j++ {
							for _, f := range fields {
								for _, ff := range f {
									_, _, _, _ = jsonparser.Get(json, ff...)
								}
							}
						}
					}
				}
			},
		},
		{
			name: "gjson",
			fn: func(b *testing.B, jsons [][]byte, fields [][][]string, reqCount int) {
				selectors := make([][]string, 0, 0)
				for _, f := range fields {
					selectorList := make([]string, 0, 0)
					for _, ff := range f {
						selectorList = append(selectorList, strings.Join(ff, "."))
					}
					selectors = append(selectors, selectorList)
				}

				for i := 0; i < b.N; i++ {
					for _, json := range jsons {
						for j := 0; j < reqCount; j++ {
							for _, selector := range selectors {
								_ = gjson.GetManyBytes(json, selector...)
							}
						}
					}
				}
			},
		},
		{
			name: "go-simplejson",
			fn: func(b *testing.B, jsons [][]byte, fields [][][]string, reqCount int) {
				for i := 0; i < b.N; i++ {
					for _, json := range jsons {
						c, _ := simplejson.NewJson(json)
						for j := 0; j < reqCount; j++ {
							for _, f := range fields {
								for _, ff := range f {
									c.GetPath(ff...)
								}
							}
						}
						s, _ = c.Encode()
					}
				}
			},
		},
	}

	workload, size := getStableWorkload()
	for _, pretender := range pretenders {
		b.Run("stable-flavor|"+pretender.name, func(b *testing.B) {
			b.SetBytes(size * int64(len(requestsCount)))
			b.ResetTimer()
			for _, reqCount := range requestsCount {
				for _, w := range workload {
					pretender.fn(b, [][]byte{w.json}, [][][]string{w.requests}, reqCount)
				}
			}
		})
	}

	//todo: we are loosing this benchmark because poor Dig() performance
	workloads, requests, size := getChaoticWorkload()
	for _, pretender := range pretenders {
		b.Run("chaotic-flavor|"+pretender.name, func(b *testing.B) {
			b.SetBytes(size * int64(len(requestsCount)))
			b.ResetTimer()
			for _, reqCount := range requestsCount {
				pretender.fn(b, workloads, requests, reqCount)
			}
		})
	}
}
