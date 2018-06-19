// The gen command builds an example yaml file that's fairly massive to use in
// benchmarking.
//
// Some quick design notes...
// - Instead of using yaml tooling using templates
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/mattfarina/yamlbench"

	v2 "github.com/mattfarina/helm-index-json"
)

var header = `apiVersion: v1
entries:
`

var footer = "generated: 2018-04-11T16:56:56.656249201Z"

var chart = `  dummy-chart-{{.Num}}:
`

var release = `  - created: 2017-07-06T01:33:50.952906435Z
    description: Example description
    digest: 249e27501dbfe1bd93d4039b04440f0ff19c707ba720540f391b5aefa3571455
    home: https://example.com
    icon: https://example.com/foo.png
    keywords:
    - A
    - B
    maintainers:
    - email: bar@example.com
      name: Bar
    name: dummy-chart-{{.Num}}
    sources:
    - https://example.com
    - https://example.com
    urls:
    - https://example.com
    version: 1.2.{{.Num2}}
`

type wrapper struct {
	Num  int
	Num2 int
}

// TODO: Rewrite this so it's much more efficient... if it ever matters
func main() {

	// Generate a YAML file for testing
	genYaml()

	// Also generate a json version of the same file content for testing
	genJson()
	genSplitJson()

	// Generate split json

	fmt.Println("Done generating testing files")
}

func genYaml() {
	// TODO(mattfarina): create an argument to capture the file name.

	fmt.Println("Generating index.yaml for testing")

	os.MkdirAll("./testfiles/yaml", os.ModePerm)
	f, err := os.Create("./testfiles/yaml/index.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.WriteString(header)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	charttmpl, err := template.New("chart").Parse(chart)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reltmpl, err := template.New("release").Parse(release)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var w wrapper
	for i := 0; i < 100; i++ {
		w = wrapper{Num: i}
		err = charttmpl.Execute(f, w)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for j := 0; j < 5000; j++ {
			w.Num2 = j
			err = reltmpl.Execute(f, w)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	_, err = f.WriteString(footer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genJson() {
	yml, err := ioutil.ReadFile("./testfiles/yaml/index.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Generating json files")

	in := &yamlbench.IndexFile{}
	err = yaml.Unmarshal(yml, &in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println("Generating index.pretty.json for testing")
	// out, err := json.MarshalIndent(in, "", "  ")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = ioutil.WriteFile("./index.pretty.json", out, 0644)

	fmt.Println("Generating index.json for testing")
	out, err := json.Marshal(in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.MkdirAll("./testfiles/json", os.ModePerm)
	err = ioutil.WriteFile("./testfiles/json/index.json", out, 0644)
}

func genSplitJson() {
	yml, err := ioutil.ReadFile("./testfiles/yaml/index.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Generating split json files")

	in := &yamlbench.IndexFile{}
	err = yaml.Unmarshal(yml, &in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.MkdirAll("./testfiles/splitjson", os.ModePerm)

	v2ft := v2.IndexFile{}
	v2ft.APIVersion = "v2"
	v2ft.Generated = in.Generated
	v2ft.Entries = make(map[string]yamlbench.ChartVersion)

	for k, v := range in.Entries {
		i := 0

		// start the chartfile
		curChart := v2.ChartFile{
			Generated:  in.Generated,
			APIVersion: "v2",
			Name:       k,
		}

		// iterate over the versions and add them
		for kk, vv := range v {
			curChart.Versions = append(curChart.Versions, vv)
			if kk > i {
				i = kk
			}
		}

		// write the chartfile to disk
		out, err := json.Marshal(curChart)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = ioutil.WriteFile("./testfiles/splitjson/"+k+".json", out, 0644)

		// add the last version to the index
		v2ft.Entries[k] = in.Entries[k][i]
	}

	// write the index to disk
	out, err := json.Marshal(v2ft)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile("./testfiles/splitjson/index.json", out, 0644)

}
