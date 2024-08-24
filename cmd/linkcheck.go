package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"golang.org/x/net/html"
)

func main() {
	startURL := "https://exactcode.com"
	domain := getDomain(startURL)

	visited := map[string]bool{}
	brokenLinks := map[string]int{}
	toVisit := []string{startURL}

	for len(toVisit) > 0 {
		// Dequeue a URL to visit
		currentURL := toVisit[0]
		toVisit = toVisit[1:]

		if visited[currentURL] {
			continue
		}
		visited[currentURL] = true

		// Get all links on the current page
		links := getAllLinks(currentURL, domain)

		for _, link := range links {
			if !visited[link] {
				toVisit = append(toVisit, link)
			}
			// Check if the link is broken
			status := checkLink(link)
			if status >= 400 {
				brokenLinks[link] = status
			}
		}
	}

	// Output broken links
	if len(brokenLinks) > 0 {
		fmt.Println("Broken links found:")
		for link, status := range brokenLinks {
			fmt.Printf("%s: %d\n", link, status)
		}
	} else {
		fmt.Println("No broken links found.")
	}
}

func getDomain(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

func getAllLinks(pageURL string, domain string) []string {
	links := []string{}
	resp, err := http.Get(pageURL)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", pageURL, err)
		return links
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return links
			}
			fmt.Printf("Error parsing HTML: %v\n", tokenizer.Err())
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						href := attr.Val
						fullURL := resolveURL(pageURL, href)
						if isSameDomain(fullURL, domain) {
							links = append(links, fullURL)
						}
					}
				}
			}
		}
	}
}

func resolveURL(base string, href string) string {
	parsedBase, err := url.Parse(base)
	if err != nil {
		return href
	}
	parsedHref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return parsedBase.ResolveReference(parsedHref).String()
}

func isSameDomain(link string, domain string) bool {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return false
	}
	return parsedLink.Host == domain
}

func checkLink(link string) int {
	resp, err := http.Head(link)
	if err != nil {
		return 500 // Consider it as server error
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
