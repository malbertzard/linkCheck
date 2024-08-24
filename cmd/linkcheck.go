package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"io"
	"golang.org/x/net/html"
)

type BrokenLink struct {
	Link     string
	Status   int
	FoundAt  string
}

func main() {
	startURL := flag.String("url", "https://example.com", "The URL to start crawling from")
	outputFile := flag.String("output", "", "The file path to export broken links as CSV (optional)")

	flag.Parse()

	domain := getDomain(*startURL)
	visited := map[string]bool{}
	brokenLinks := []BrokenLink{}
	toVisit := []string{*startURL}

	for len(toVisit) > 0 {
		currentURL := toVisit[0]
		toVisit = toVisit[1:]

		if visited[currentURL] {
			continue
		}
		visited[currentURL] = true

		links := getAllLinks(currentURL, domain)

		for _, link := range links {
			if !visited[link] {
				toVisit = append(toVisit, link)
			}

			status := checkLink(link)
			if status >= 400 {
				brokenLinks = append(brokenLinks, BrokenLink{
					Link:    link,
					Status:  status,
					FoundAt: currentURL,
				})
			}
		}
	}

	if *outputFile != "" {
		exportCSV(*outputFile, brokenLinks)
	} else {
		printBrokenLinks(brokenLinks)
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

func printBrokenLinks(brokenLinks []BrokenLink) {
	if len(brokenLinks) > 0 {
		fmt.Println("Broken links found:")
		for _, bl := range brokenLinks {
			fmt.Printf("%s: %d (found at: %s)\n", bl.Link, bl.Status, bl.FoundAt)
		}
	} else {
		fmt.Println("No broken links found.")
	}
}

func exportCSV(filename string, brokenLinks []BrokenLink) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Link", "Status Code", "Found At"})

	// Write broken links data
	for _, bl := range brokenLinks {
		writer.Write([]string{bl.Link, fmt.Sprintf("%d", bl.Status), bl.FoundAt})
	}

	fmt.Printf("Broken links exported to %s\n", filename)
}
