package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var MandatoryArgs = []string{"--csv"}
var SupportedArgs = []string{"--num", "--goroutines"}

func main() {
	inArgs := readInputArgs(os.Args[1:])
	var err error
	if inArgs == nil {
		log.Fatal("Bad input arguments")
	}
	log.Println(inArgs)

	// read input arguments to variables
	var max int
	if val, ok := inArgs["--num"]; ok {
		max, err = strconv.Atoi(val)
		if err != nil || max < 1 {
			log.Fatalf("Bad --num, %v, %v\n", val, err)
		}
	}

	goroutines := 1
	if val, ok := inArgs["--goroutines"]; ok {
		goroutines, err = strconv.Atoi(val)
		if err != nil || goroutines < 1 {
			log.Fatalf("Bad --num, %v, %v\n", val, err)
		}
	}

	// read list of URLs to crawl
	listURLs, err := readCSVFile(inArgs["--csv"])
	if err != nil {
		log.Fatal(err)
	}

	// begin crawling
	// show a breakdown of URLs per goroutine
	if max == 0 || len(listURLs) < max {
		max = len(listURLs)
	}
	var wg sync.WaitGroup                // used to check when ALL goroutines are done
	urlsPerGoroutine := max / goroutines // compute minimum number of URLs for each goroutine

	log.Printf("urlsPerGoroutine: %d\n", urlsPerGoroutine)
	log.Printf("Crawling %d URLs with %d goroutines\n", max, goroutines)
	for i := 0; i < goroutines; i++ {
		start := urlsPerGoroutine * i
		end := start + urlsPerGoroutine
		if i == goroutines-1 {
			end = max
		}

		log.Printf("Starting goroutine %d, from %d to %d\n", i, start, end-1)
		wg.Add(1)
		go DoCrawl(listURLs[start:end], i, &wg)

	}

	log.Println("Waiting for goroutines to finish ...")
	wg.Wait()
	log.Println("DONE")

}

func readInputArgs(a []string) map[string]string {
	// just read os.Args in a map
	tmpMap := make(map[string]string)

	if len(a)%2 != 0 || len(a) == 0 {
		return nil
	}

	for i := 0; i < len(a); i += 2 {
		if Contains(MandatoryArgs, a[i]) || Contains(SupportedArgs, a[i]) {
			tmpMap[a[i]] = a[i+1]
		} else {
			// unrecognised input argument
			log.Printf("Unrecognised input argument %s\n", a[i])
			return nil
		}
	}

	return tmpMap
}

func readCSVFile(fileLoc string, entries ...int) ([]string, error) {
	log.Printf("Reading CSV file at %s\n", fileLoc)
	max := -1
	if len(entries) == 1 {
		max = entries[0]
	}

	csvFile, err := os.Open(fileLoc)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var list []string
	counter := 0
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		}

		if len(line[1]) != 0 && counter > 0 {
			list = append(list, line[1])
		}
		if max != -1 && len(list) >= max {
			break
		}
		counter++
	}

	return list, nil
}

func DoCrawl(listURLs []string, routineID int, wg *sync.WaitGroup) {
	// start running HTTP GETs to the URLs provided in listURLs
	// log.Printf("%d, %v\n", routineID, listURLs) // DEBUG
	defer wg.Done()
	client := &http.Client{
		// do not follow any redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for i := 0; i < len(listURLs); i++ {
		tmpURL := "https://" + listURLs[i]
		req, err := http.NewRequest("GET", tmpURL, nil)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		log.Printf("%d %d - %s\n", routineID, i, showReqInfo(req, resp))

	}
}

func showReqInfo(r *http.Request, a *http.Response) string {
	// show a summary of the HTTP response received as a string
	str := fmt.Sprintf("HTTP GET %s, Status: %s, ContentLength: %d", r.URL.String(), a.Status, a.ContentLength)

	return str
}

// helper functions
func Contains(a []string, x string) bool {
	// Contains tells whether a contains x.
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsAt(a []string, x string) int {
	// similar to Contains, returns index of first occurrence
	for i := 0; i < len(a); i++ {
		if a[i] == x {
			return i
		}
	}
	return -1
}
