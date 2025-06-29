package fetchers

import (
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/common/logger"
	"io"
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
		logger.Error(err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to fetch bus stations: status code %d", resp.StatusCode)
		logger.Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return html, nil
}
