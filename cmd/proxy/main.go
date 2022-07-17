package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kfelter/unified-cache-example/internal/bucket"
	"github.com/kfelter/unified-cache-example/internal/ketama"
)

func main() {
	nodes := flag.String("n", "http://localhost:3001;http://localhost:3002;http://localhost:3003", "proxy nodes")
	flag.Parse()
	// get ketama buckets from nodes
	buckets := parseBuckets(*nodes)

	// internal bucket.Bucket implements the ketama.Bucket interface
	// with each label being the destination of the cache node that contains
	// the requested object
	continuum := ketama.New(buckets)

	http.HandleFunc("/_/metrics", collectMetrics(continuum))

	// when a request comes in to the proxy, hash it using ketama to find the cache node
	// the resource should be found in
	http.HandleFunc("/", redirectToNode(continuum))

	fmt.Printf("buckets: %+v\n", continuum.Buckets())
	fmt.Println("proxy serving on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

func parseBuckets(nodes string) []ketama.Bucket {
	buckets := []ketama.Bucket{}
	ss := strings.Split(nodes, ";")
	for i := range ss {
		buckets = append(buckets, bucket.New(ss[i]))
	}
	return buckets
}

func collectMetrics(continuum *ketama.Continuum) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responses := map[string]map[string]int{}
		for _, b := range continuum.Buckets() {
			res, err := http.Get(b.Label() + "/_/metrics")
			if err != nil {
				fmt.Println("err getting:", b.Label()+"/_/metrics", err)
				continue
			}
			data, _ := io.ReadAll(res.Body)
			m := map[string]int{}
			json.Unmarshal(data, &m)
			responses[b.Label()] = m
		}
		w.WriteHeader(http.StatusOK)
		byt, _ := json.Marshal(responses)
		w.Write(byt)
	}
}

func redirectToNode(continuum *ketama.Continuum) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlString := r.URL.String()
		nodeAddr := continuum.Hash([]byte(urlString)).Label()
		fmt.Println(urlString, "->", nodeAddr+urlString)
		http.Redirect(w, r, nodeAddr+urlString, http.StatusTemporaryRedirect)
	}
}
