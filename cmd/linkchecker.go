package main

import (
	"net/http"
)

// checkLink sends a HEAD request to the link and returns the status code.
func checkLink(link string) int {
	resp, err := http.Head(link)
	if err != nil {
		return 500 // Consider it as server error
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
