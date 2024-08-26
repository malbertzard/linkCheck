package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

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

	writer.Write([]string{"Link", "Status Code", "Found At"})

	for _, bl := range brokenLinks {
		writer.Write([]string{bl.Link, fmt.Sprintf("%d", bl.Status), bl.FoundAt})
	}

	fmt.Printf("Broken links exported to %s\n", filename)
}
