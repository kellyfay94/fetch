package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// OutputMsg - Format for pushing messages logger
type OutputMsg struct {
	Msg string
	Err error
}

// Default Values
const (
	defaultRequestTimeoutMS = 2500
	defaultDialTimeoutMS    = 3000
)

// FUTURE IMPROVEMENTS:
//		- Custom Timeouts
//		- Proxy Settings
//		- Custom Output Destinations

// Fetcher - Client for fetching pages and loading metadata while communicating to console
type Fetcher struct {
	client *http.Client

	msgs   chan OutputMsg
	msgsMu sync.Mutex
}

// NewFetcher - Initializes a Fetcher struct with a default HTTP Client with tunable timeouts
func NewFetcher() *Fetcher {
	// FUTURE IMPROVEMENTS: Read from "config", utilize optional proxies

	dialer := net.Dialer{
		Timeout: time.Duration(defaultDialTimeoutMS) * time.Millisecond,
	}

	httpclient := &http.Client{
		Timeout: time.Duration(defaultRequestTimeoutMS) * time.Millisecond,
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}

	return &Fetcher{
		client: httpclient,
		msgs:   make(chan OutputMsg),
	}
}

// Msgs - Returns the msgs table
func (fetcher *Fetcher) Msgs() <-chan OutputMsg {
	return fetcher.msgs
}

// Shutdown - Shutsdown the Fetcher, should be called when all operations are done
func (fetcher *Fetcher) Shutdown() {
	fetcher.msgsMu.Lock()
	defer fetcher.msgsMu.Unlock()

	select {
	case <-fetcher.msgs:
		// if this returns, then the 'fetcher' channel is already closed
	default:
		close(fetcher.msgs)
	}
}

// FetchPage - Given a URL, fetches the page
func (fetcher *Fetcher) FetchPage(url string) (Page, OutputMsg) {
	// FUTURE IMPROVEMENT: Add a listener for os.SIGTERM to safely cancel the request

	// Build the Request
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return Page{}, OutputMsg{"Error occurred during Fetch Request generation", err}
	}

	// Fetch The Page
	resp, err := fetcher.client.Do(req)
	if err != nil {
		return Page{}, OutputMsg{"Error occurred while Fetching the Page", err}
	}
	// Store the approximate completed time that we received the page
	fetchedTime := time.Now().UTC()
	defer resp.Body.Close()

	// Parse the Body of the Page
	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Page{}, OutputMsg{"Error occurred while Parsing the Page", err}
	}

	// Build and return the Page object
	return Page{
		PageBodyBytes: bodyB,
		PageMetadata: PageMetadata{
			URL:           url,
			LastFetchTime: fetchedTime,
		},
	}, OutputMsg{"", nil}
}

// SavePage - Saves the page to the local disk
func (fetcher *Fetcher) SavePage(p Page) OutputMsg {
	// Based on the URL, generate the base of the filename
	// NOTE: As the filename schema is not completely defined in the assignment document,
	//	 	 it is assumed that if the page's full path contains "/", that it can be replaced
	//		 with another character and should still be stored
	filename := parseFilename(p.URL)

	// Write the Page to a File
	err := os.WriteFile(filename+".html", p.PageBodyBytes, 0777)
	if err != nil {
		return OutputMsg{fmt.Sprintf("Error storing HTML of Page: %s", p.URL), err}
	}

	// Convert the Page's Metadata into a Marshalled JSON
	metadataBytes, err := json.Marshal(p.PageMetadata)
	if err != nil {
		return OutputMsg{fmt.Sprintf("Error rendering the Metadata for Page: %s", p.URL), err}
	}

	// Write the Page's Metadata to a JSON
	err = os.WriteFile(filename+"-metadata.json", metadataBytes, 0777)
	if err != nil {
		return OutputMsg{fmt.Sprintf("Error storing the Metadata for Page: %s", p.URL), err}
	}

	return OutputMsg{"", nil}
}

// LoadPageMetadata - Loads the Page's Metadata
func (fetcher *Fetcher) LoadPageMetadata(url string) OutputMsg {
	filename := parseFilename(url) + "-metadata.json"
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return OutputMsg{fmt.Sprintf("Metadata does not exist for Page: %s", url), nil}
	} else if err != nil {
		return OutputMsg{fmt.Sprintf("Error opening Metadata for Page: %s", url), err}
	}

	fileB, err := ioutil.ReadAll(file)
	if err != nil {
		return OutputMsg{fmt.Sprintf("Error reading Metadata storage for Page: %s", url), err}
	}

	pM := &PageMetadata{}
	err = json.Unmarshal(fileB, pM)
	if err != nil {
		return OutputMsg{fmt.Sprintf("Error parsing Metadata storage for Page: %s", url), err}
	}

	return OutputMsg{pM.String(), nil}
}

// FetchURLs - Given a list of URLs, fetches each Page and stores its HTML and related metadata to Disk; errors are logged to the Fetcher's Msg channel
func (fetcher *Fetcher) FetchURLs(urls []string) {
	for _, url := range urls {
		// Fetch Page
		//	if error occurs, proceed to the next URL - do not store
		page, msg := fetcher.FetchPage(url)
		if msg.Err != nil {
			fetcher.msgs <- msg
			continue
		}

		err := page.ExtractMetadata()
		if err != nil {
			fetcher.msgs <- OutputMsg{"Error when extracting metadata", err}
		}

		// Store Page
		//	if error occurs, push error message to fetcher channel; continue
		msg = fetcher.SavePage(page)
		if msg.Err != nil {
			fetcher.msgs <- msg
			continue
		}
	}
}

// ReturnMetadata - Provided a list of URLs, loads the metadata for each and pushes the metadata to the Fetcher's Msg channel
func (fetcher *Fetcher) ReturnMetadata(urls []string) {
	for _, url := range urls {
		// Load the Page Metadata, and push the msg to the outMsg
		msg := fetcher.LoadPageMetadata(url)
		fetcher.msgs <- msg
	}
}
