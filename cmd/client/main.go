package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	filePath := flag.String("f", "*", "file path to get")
	proxyAddr := flag.String("a", "http://localhost:3000", "proxy addr")
	flag.Parse()

	if *filePath == "*" {
		for i := 0; i < 255; i++ {
			res, err := http.Get(*proxyAddr + fmt.Sprintf("/img/%d.png", i))
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			f, _ := os.Create("client_img.png")
			io.Copy(f, res.Body)
			data, _ := httputil.DumpResponse(res, false)
			fmt.Println(string(data))
		}
		return
	}

	res, err := http.Get(*proxyAddr + *filePath)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	f, _ := os.Create("client_img.png")
	io.Copy(f, res.Body)
	data, _ := httputil.DumpResponse(res, false)
	fmt.Println(string(data))

}
