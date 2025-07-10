package marprom

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
}

type HTMLFetcher struct {
	client *http.Client
}

func NewHTMLFetcher() *HTMLFetcher {
	return &HTMLFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type FetchOptions struct {
	URL string
}

func (f *HTMLFetcher) FetchHTML(opts *FetchOptions) ([]byte, error) {
	req, err := http.NewRequest("GET", opts.URL, nil)
	if err != nil {
		log.Printf("Failed to create request: %s", err)
		return nil, err
	}

	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set("User-Agent", ua)

	resp, err := f.client.Do(req)
	if err != nil {
		log.Printf("Error fetching URL %s: %s", opts.URL, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response from %s: %s", opts.URL, resp.Status)
		return nil, err
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return nil, err
	}

	return html, nil
}
