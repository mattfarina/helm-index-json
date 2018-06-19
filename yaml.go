package indexjson

import (
	"time"

	"github.com/mattfarina/yamlbench"
)

type IndexFile struct {
	APIVersion string    `json:"apiVersion"`
	Generated  time.Time `json:"generated"`

	// Difference from v1 is that entries only has one ChartVersion, the latest, rather than all of them.
	Entries map[string]yamlbench.ChartVersion `json:"entries"`
}

// ChartFile is a new file in v2 and contains all the versions of a chart
type ChartFile struct {
	APIVersion string                  `json:"apiVersion"`
	Generated  time.Time               `json:"generated"`
	Name       string                  `json:"name"`
	Versions   yamlbench.ChartVersions `json:"versions"`
}
