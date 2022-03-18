package fetcher

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/html"
)

// Page - Struct to Store Page and relevant Metadata
type Page struct {
	PageBodyBytes []byte
	PageMetadata
}

// ExtractMetadata - Extracts the metadata from the page (LinkCount and Images)
func (p *Page) ExtractMetadata() error {
	// FUTURE IMPROVEMENT: Handle "broken" HTML that the HTML Parser seems unable to handle, like the format available on www.google.com/test

	// Tokenize the HTML
	token := html.NewTokenizer(bytes.NewBuffer(p.PageBodyBytes))

	// Iterate through the tokenized HTML
	for {
		tokenType := token.Next()
		switch {
		case tokenType == html.ErrorToken:
			err := token.Err()
			// `io.EOF` is a an "error" that simply indicates end of file; don't pass this to the rest of the app
			if err == io.EOF {
				return nil
			}
			return err
		case tokenType != html.EndTagToken:
			nameB, _ := token.TagName()

			// FUTURE IMPROVEMENT: While iterating through Nodes, also identify attributes to be fetched, and add to `Page` struct

			switch string(nameB) {
			case "img":
				// NOTE: As the assignment wasn't specific, we're just counting 'img' images specifically
				//		 Any item that may be a background image as part of a CSS or JS addition is not included
				p.PageMetadata.Images++
			case "a":
				// NOTE: As the assignment wasn't specific, we're just counting 'a' links specifically
				//		 Any item that may cause a redirect and is not a typical "link"
				//		 (Like a DIV with an attached redirect JS) is not counted
				p.PageMetadata.LinkCount++
			}
		}
	}
}

// PageMetadata - Struct to store / Retrieve Page Metadata
type PageMetadata struct {
	URL           string    `json:"site"`
	LinkCount     int       `json:"num_links"`
	Images        int       `json:"num_images"`
	LastFetchTime time.Time `json:"last_fetch"`
}

func (pM PageMetadata) String() string {
	return fmt.Sprintf("site: %s\n num_links: %v\n images: %v\n last_fetch: %s",
		pM.URL,
		pM.LinkCount,
		pM.Images,
		pM.LastFetchTime.Format("Mon Jan _2 15:04:05 2006 MST"), // Using the format from the `Section 2` example output
	)
}
