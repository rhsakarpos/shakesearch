package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type AllSearchResults struct {
	Results []SearchResult
}

type SearchResult struct {
	Title   string
	Matches []string
}

func main() {

	searcher := Searcher{}

	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	_ = buildSearchArraysByTitle("works//completeworks.txt", "works//workslist.txt", searcher)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type Searcher struct {
	MapTitleVsCompleteWorks map[string]string
	MapTitleVsSuffixArray   map[string]*suffixarray.Index
	//SuffixArray   *suffixarray.Index
	SearchSet string
}

func buildSearchArraysByTitle(filename string, titlesfilename string, searcher Searcher) error {

	// load the full file
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	fullText := string(dat)

	fullArray := suffixarray.New(dat)
	// load the titles
	linesArr, _ := readLines(titlesfilename)

	nextTitleIndex := 1
	for tIndex, _ := range linesArr {
		title := strings.TrimSpace(linesArr[tIndex])
		if len(title) != 0 {
			if nextTitleIndex < len(linesArr) {
				nextTitle := strings.TrimSpace(linesArr[nextTitleIndex])

				beginregex := regexp.MustCompile(title)
				idxs := fullArray.FindAllIndex(beginregex, -1)
				beginIndex := 0
				for _, ele := range idxs {
					beginIndex = ele[0]
				}

				endregex := regexp.MustCompile(nextTitle)
				idxs = fullArray.FindAllIndex(endregex, -1)
				endIndex := 0
				for _, ele := range idxs {
					endIndex = ele[0]
				}
				fmt.Println(title)
				fmt.Println("\t" + nextTitle)
				searcher.MapTitleVsSuffixArray[title] = suffixarray.New([]byte(fullText[beginIndex:endIndex]))
				searcher.MapTitleVsCompleteWorks[title] = fullText[beginIndex:endIndex]
				nextTitleIndex++
			}

		}
	}

	return nil

}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := AllSearchResults{}
		results.Results = make([]SearchResult, 0)
		searcher.SearchAllCase(query[0], &results)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
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
	/*dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SearchSet = s.CompleteWorks*/
	s.MapTitleVsSuffixArray = make(map[string]*suffixarray.Index)
	s.MapTitleVsCompleteWorks = make(map[string]string)
	//s.SuffixArray = suffixarray.New([]byte(s.SearchSet))
	return nil
}

func (s *Searcher) reloadSearchSet() error {
	//s.SuffixArray = suffixarray.New([]byte(s.SearchSet))
	return nil
}

func (s *Searcher) SearchAllCase(query string, allResult *AllSearchResults) {
	//allPerms := generateCasePerms(query)
	var allPerms []string
	allPerms = append(allPerms, query)
	for _, searchVal := range allPerms {
		_ = s.Search(searchVal, allResult)
	}
}

func (s *Searcher) Search(query string, allResult *AllSearchResults) []string {
	stringResults := []string{}
	for title, val := range s.MapTitleVsSuffixArray {
		result := s.SearchInSuffixArray(query, title, val)
		if len(result) > 0 {
			stringResults = append(stringResults, result...)
			sr := SearchResult{Title: title, Matches: stringResults}
			allResult.Results = append(allResult.Results, sr)
		}
	}
	return stringResults
}

func (s *Searcher) SearchInSuffixArray(query string, title string, sArray *suffixarray.Index) []string {
	idxs := sArray.Lookup([]byte(query), -1)
	results := []string{}
	for _, idx := range idxs {
		orig := s.MapTitleVsCompleteWorks[title][idx-250 : idx+250]
		s.SearchSet = strings.ReplaceAll(s.SearchSet, orig, "")
		formatted := strings.ReplaceAll(orig, query, "<span style='font-weight:bold;'>"+query+"</span>")
		results = append(results, formatted)
	}
	fmt.Println(results)
	return results
}
