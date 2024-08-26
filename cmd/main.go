package main

import (
	"flag"
)

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
