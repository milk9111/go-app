package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func crawl(url string, tag string, ch chan int, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	var amount int
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			ch <- amount
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isCorrectTag := t.Data == tag
			if !isCorrectTag {
				continue
			}

			amount++
		}
	}
}

func Scrape(url string, tag string) int {
	var amount int

	// Channels
	chUrls := make(chan int)
	chFinished := make(chan bool)

	go crawl(url, tag, chUrls, chFinished)

	// Subscribe to both channels
	for c := 0; c < 1; {
		select {
		case url := <-chUrls:
			amount = url
		case <-chFinished:
			c++
		}
	}

	// We're done! Print the results...

	fmt.Println("\nFound", amount, "tags")

	close(chUrls)

	return amount
}