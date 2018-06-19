# JSON Spilt testing

This repository is used for testing a split of the Helm index with all metadata into multiple files with a file per chart and an index listing all charts with the metadata for the latest chart.

There are three pieces of information we are currently looking at:

1. The size of files. This is due to the size being transferred over the wire
1. The time to generate the different forms of files. To understand if one method or another has advantages
1. The performance of reading and working with files


## Findings Summary

For the findings we look at 2 different cases

1. A large index with many charts. This is useful to know if we are working with a central service that aggregates charts and knows about a large number of them.
1. An individual repository using a devops workflow that causing many chart versions to be published.

### File Sizes

* Spec v1: An index.json file with 100 charts and 5,000 chart versions per chart was:
  * 348.8mb in size for a pretty/readable format
  * 216.3mb in size normally
* Proposed spec v2:
  * A repo with 100 charts and 5,000 chart versions per chart was:
    * 45kb index.json file
    * Each chart json file was 2.2mb
  * A repo with 5,000 charts and 100 chart versions per chart was:
    * 2.3mb index.json file
    * Each chart json file was 43kb

### Generation Time

* Spec v1: An index.json file with 5,000 charts and 100 chart versions per chart was: 41.34s
* Spec v2: An repo with 5,000 charts and 100 chart versions per chart was: 43.38s

### File Read Time

In these cases this is for a repo with 5,000 charts and 100 chart versions each

```
BenchmarkJson-8             	       1	5289763888 ns/op	788211456 B/op	12620257 allocs/op
BenchmarkSplitJsonIndex-8   	      30	  50010804 ns/op	 8083786 B/op	  140155 allocs/op
BenchmarkSplitJsonChart-8   	    2000	    904676 ns/op	  242488 B/op	    2329 allocs/op
```

The initial Json one is the v1 spec while the SplitJson ones are the proposed v2 ones. The first v2 one is the index with 5,000 charts and the latest chart version details included. The second split one is an individual chart file with 100 chart versions.

For the next set of cases this is for a repo with 100 charts and 5,000 chart versions each

```
BenchmarkJson-8             	       1	5209633527 ns/op	782795552 B/op	12504360 allocs/op
BenchmarkSplitJsonIndex-8   	    2000	    972878 ns/op	  159004 B/op	    2817 allocs/op
BenchmarkSplitJsonChart-8   	      30	  45574826 ns/op	 9508206 B/op	  115047 allocs/op
```

The gist is that reading an index file with all the charts and version in 1 file takes over 5 seconds and uses 782MB of memory while the case of an index with just the latest chart version and separate files for each chart containing all the version causes worst case file read times of 0.5 seconds and 9.5mb of memory.

## Test Environment

The environment the tests were performed in was on a 2017 MacBook Pro with an SSD

Note, the main.go file was manually altered to run different scenarios.