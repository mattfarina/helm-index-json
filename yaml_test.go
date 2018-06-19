package indexjson

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/mattfarina/yamlbench"
)

func BenchmarkJson(b *testing.B) {
	index, err := ioutil.ReadFile("./testfiles/json/index.json")
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		var in yamlbench.IndexFile2
		err := json.Unmarshal(index, &in)
		if err != nil {
			b.Errorf("json err: %s", err)
		}
	}
}

func BenchmarkSplitJsonIndex(b *testing.B) {
	index, err := ioutil.ReadFile("./testfiles/splitjson/index.json")
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		var in IndexFile
		err := json.Unmarshal(index, &in)
		if err != nil {
			b.Errorf("json err: %s", err)
		}
	}
}

func BenchmarkSplitJsonChart(b *testing.B) {
	index, err := ioutil.ReadFile("./testfiles/splitjson/dummy-chart-0.json")
	if err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		var in ChartFile
		err := json.Unmarshal(index, &in)
		if err != nil {
			b.Errorf("json err: %s", err)
		}
	}
}
