package main

import (
	"flag"
	"fmt"

	"github.com/kellyfay94/fetch/fetcher"
)

// ParseCommands - Parses the Command-Line Arguments and Flags
func ParseCommands() ([]string, bool, bool) {
	verbose := false
	flag.BoolVar(&verbose, "v", false, "Prints the ")

	metadataRequest := false
	flag.BoolVar(&metadataRequest, "metadata", false, "Returns the metadata of previous requests to URL(s)")
	flag.Parse()

	urls := flag.Args()

	if verbose {
		fmt.Printf("URLs: %+v\n", urls)
		action := ""
		switch metadataRequest {
		case true:
			action = "Loading cached metadata..."
		default:
			action = "Fetching new page(s)..."
		}
		fmt.Println(action)
	}

	return urls, metadataRequest, verbose
}

func main() {
	urls, metadataRequest, _ := ParseCommands()
	f := fetcher.NewFetcher()
	if metadataRequest {
		go func() {
			f.ReturnMetadata(urls)
			f.Shutdown()
		}()
	} else {
		go func() {
			f.FetchURLs(urls)
			f.Shutdown()
		}()
	}
	for msg := range f.Msgs() {
		fmt.Println(msg.Msg)
		if msg.Err != nil {
			fmt.Printf("\tErr Details: %+v", msg.Err)
		}
	}
}
