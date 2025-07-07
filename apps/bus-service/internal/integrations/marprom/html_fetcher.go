package marprom

import (
	"io"
	"log"
	"net/http"
)

type HTMLFetcher struct {
	client *http.Client
}

func NewHTMLFetcher() *HTMLFetcher {
	return &HTMLFetcher{
		client: &http.Client{},
	}
}

type FetchOptions struct {
	URL string
}

func (f *HTMLFetcher) FetchHTML(opts *FetchOptions) ([]byte, error) {
	resp, err := f.client.Get(opts.URL)
	if err != nil {
		log.Fatalf("Error fetching URL %s: %s", opts.URL, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching URL %s: %s", opts.URL, resp.Status)
		return nil, err
	}

	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
		return nil, err
	}

	return html, nil
}
