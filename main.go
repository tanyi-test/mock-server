package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	RequestMap map[string]string
)

func handler(w http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf("%s %s", req.Method, req.URL)
	fmt.Printf("[DEBUG] [URL] %s", s)
	if respStr, ok := RequestMap[s]; ok {
		resp := []byte(respStr)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(resp)))
		w.WriteHeader(200)

		total := 0
		for total < len(resp) {
			count, err := w.Write(resp[total:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] [%s] %v", s, err)
			}
			total += count
		}
	} else {
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(404)
	}
}

func main() {
	reqMapBytes, err := ioutil.ReadFile("request.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(reqMapBytes, &RequestMap)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
