package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

const urlRegexp = `(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		count, err := GetAndCountLinksAt(s)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Count for %s: %d\n", s, count)
		}
	}
}

// GetAndCountLinksAt makes an HTTP GET request to a given url and counts all
// links in a response body.
func GetAndCountLinksAt(url string) (int, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(urlRegexp)
	urls := re.FindAll(b, -1)
	if urls != nil {
		return len(urls), nil
	}
	return 0, nil
}
