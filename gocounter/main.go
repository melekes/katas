// Reads URLs from stdin. For each URL it performs a GET request and counts all
// the occurrences of string "Go" in response body.
//
// Usage:
//   echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func main() {
	total := 0

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("Total: %d\n", total)
			os.Exit(0)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		count, err := GetAndCountGoAt(s)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Count for %s: %d\n", s, count)
			total += count
		}
	}

	fmt.Printf("Total: %d\n", total)
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
