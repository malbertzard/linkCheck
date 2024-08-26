package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
)

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
