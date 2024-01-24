package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const MAX_RESULTS = 20

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		batch, _ := r.URL.Query()["b"]
		batchIdx := 0
		var err error
		if len(batch) > 0 {
			batchIdx, err = strconv.Atoi(batch[0])
			if err != nil {
				batchIdx = 0
			}
		}

		results := searcher.Search(query[0], batchIdx)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err = enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(string(dat))))
	return nil
}

func (s *Searcher) Search(query string, batch int) []string {
	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	sort.Ints(idxs)
	results := []string{}
	for i, idx := range idxs {
		if i >= batch*MAX_RESULTS && i < min(len(idxs), (batch+1)*MAX_RESULTS) {
			results = append(results, s.CompleteWorks[idx-250:idx+250])
		}
	}
	return results
}
