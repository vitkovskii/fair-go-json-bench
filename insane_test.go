package insane_json_bench

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/valyala/fastjson"
	"github.com/vitkovskii/insane-json"
	"github.com/bitly/go-simplejson"
)

type test struct {
	json []byte
	name string

	digFields []string
}

func loadTest(name string, getFields []string) *test {
	content, err := ioutil.ReadFile(fmt.Sprintf("benchdata/%s.json", name))
	if err != nil {
		panic(err.Error())
	}

	return &test{json: content, name: name, digFields: getFields}
}

func getCompetitions() []*test {
	tests := make([]*test, 0, 0)
	tests = append(tests, loadTest("light-ws", []string{"about"}))
	tests = append(tests, loadTest("many-objects", []string{"somefield", "somefield", "somefield", "somefield", "somefield", "somefield", "somefield"}))
	tests = append(tests, loadTest("heavy", []string{"first", "second", "third", "fourth", "fifth"}))
	tests = append(tests, loadTest("many-fields", []string{"compfanvy"}))
	tests = append(tests, loadTest("few-fields", []string{"compfanvy"}))
	tests = append(tests, loadTest("insane", []string{"statuses", "2", "user", "entities", "url", "urls", "0", "expanded_url"}))
	return tests
}

func BenchmarkDecode(b *testing.B) {
	pretenders := []struct {
		name string
		fn   func(b *testing.B, json []byte)
	}{
		{
			name: "insane-json",
			fn: func(b *testing.B, json []byte) {
				root := insaneJSON.Spawn()
				for i := 0; i < b.N; i++ {
					_ = insaneJSON.DecodeBytesReusing(root, json)
				}
				insaneJSON.Release(root)
			},
		},
		{
			name: "fastjson",
			fn: func(b *testing.B, json []byte) {
				parser := fastjson.Parser{}
				for i := 0; i < b.N; i++ {
					_, _ = parser.ParseBytes(json)
				}
			},
		},
	}

	for _, competition := range getCompetitions() {
		for _, pretender := range pretenders {
			b.Run(competition.name+"|"+pretender.name, func(b *testing.B) {
				b.SetBytes(int64(len(competition.json)))
				b.ResetTimer()
				pretender.fn(b, competition.json)
			})
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, benchmark := range getCompetitions() {
		b.Run("insane-"+benchmark.name, func(b *testing.B) {
			root, _ := insaneJSON.DecodeBytes(benchmark.json)
			s := make([]byte, 0, 500000)
			b.SetBytes(int64(len(benchmark.json)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s = root.EncodeNoAlloc(s[:0])
			}
			insaneJSON.Release(root)
		})

		b.Run("fastjson", func(b *testing.B) {
			parser := fastjson.Parser{}
			c, _ := parser.ParseBytes(benchmark.json)
			s := make([]byte, 0, 500000)
			b.SetBytes(int64(len(benchmark.json)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s = c.MarshalTo(s[:0])
			}
		})
	}
}

func BenchmarkDig(b *testing.B) {
	for _, benchmark := range getCompetitions() {
		b.Run("insane-"+benchmark.name, func(b *testing.B) {
			root, _ := insaneJSON.DecodeBytes(benchmark.json)
			b.SetBytes(1)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				root.Dig(benchmark.digFields...)
			}
			insaneJSON.Release(root)
		})

		b.Run("fastjson", func(b *testing.B) {
			parser := fastjson.Parser{}
			c, _ := parser.ParseBytes(benchmark.json)
			b.SetBytes(1)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Get(benchmark.digFields...)
			}
		})
	}
}
