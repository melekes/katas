// Reads URLs from stdin. For each URL it performs a GET request and counts all
// the occurrences of string "Go" in response body.
//
// Usage:
//   echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	var goroutines, maxGoroutines int
	flag.IntVar(&maxGoroutines, "k", 4, "maximum number of concurrent goroutines")
	var bufferSize int
	flag.IntVar(&bufferSize, "b", 100, "size of the internal buffer for urls")

	total := 0

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("Total: %d\n", total)
			os.Exit(0)
		}
	}()

	urlsCh := make(chan string, bufferSize)
	resCh := make(chan int)

	var wg sync.WaitGroup
	quitCh := make(chan struct{})

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			s := scanner.Text()
			select {
			case urlsCh <- s:
				wg.Add(1)
			default:
				fmt.Printf("Urls buffer is over capacity. Dropping %s ...\n", s)
			}
		}
		close(quitCh)
	}()

	for {
		select {
		case url := <-urlsCh:
			if goroutines < maxGoroutines {
				goroutines++
				go func(url string, resCh chan<- int) {
					count, err := GetAndCountGoAt(url)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						fmt.Printf("Count for %s: %d\n", url, count)
						resCh <- count
					}
				}(url, resCh)
			}
		case count := <-resCh:
			total += count
			goroutines--
			wg.Done()
		case <-quitCh:
			go func() {
				wg.Wait()
				fmt.Printf("Total: %d\n", total)
				os.Exit(0)
			}()
		}
	}
}

// GetAndCountGoAt makes an HTTP GET request to a given url and counts all the
// occurrences of string "Go".
func GetAndCountGoAt(url string) (int, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	return strings.Count(string(b), "Go"), nil
}
