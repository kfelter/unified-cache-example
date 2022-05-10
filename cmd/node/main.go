package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type cache struct {
	sync.Mutex
	origin string
	codes  map[string]int
	data   map[string][]byte
	hits   int
	misses int
}

func (c *cache) Size() int {
	c.Lock()
	defer c.Unlock()
	size := 0
	for _, data := range c.data {
		size += len(data)
	}
	return size
}

func (c *cache) Metrics() []byte {
	m := map[string]interface{}{
		"hits":   c.hits,
		"misses": c.misses,
		"size":   c.Size(),
	}
	data, _ := json.Marshal(m)
	return data
}

func main() {
	portFlag := flag.String("p", os.Getenv("NODE_PORT"), "listen port")
	flag.Parse()
	// start a pull through cache node
	local := cache{
		Mutex:  sync.Mutex{},
		origin: "http://localhost:3004",
		codes:  make(map[string]int),
		data:   make(map[string][]byte),
	}

	http.HandleFunc("/_/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(local.Metrics())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		local.Lock()
		defer local.Unlock()

		cacheKey := r.URL.String()

		// respond from cache
		if data, ok := local.data[cacheKey]; ok {
			w.Header().Add("x-cache", "HIT")
			local.hits++
			w.WriteHeader(local.codes[cacheKey])
			w.Write(data)
			fmt.Println("served cached response for", cacheKey)
			return
		}

		// pull from origin
		res, err := http.Get(local.origin + r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.StatusCode != http.StatusOK {
			http.Error(w, "status="+fmt.Sprintf("%d", res.StatusCode)+" cannot cache: "+r.URL.String(), http.StatusInternalServerError)
			return
		}

		// save response
		local.codes[cacheKey] = res.StatusCode
		local.data[cacheKey] = data

		// send response to client
		w.Header().Add("x-cache", "MISS")
		local.misses++
		w.WriteHeader(res.StatusCode)
		w.Write(data)
		fmt.Println("served origin response for", cacheKey)
	})

	fmt.Println("cache node running on :" + *portFlag)
	if err := http.ListenAndServe(":"+*portFlag, nil); err != nil {
		panic(err)
	}
}
