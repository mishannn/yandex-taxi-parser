package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getProxies(url string, proxyType string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server sent status code %d, but expected %d", resp.StatusCode, 200)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	proxies := make([]string, 0)

	for _, rawLine := range strings.Split(string(respBody), "\n") {
		proxies = append(proxies, fmt.Sprintf("%s://%s", proxyType, strings.TrimSpace(rawLine)))
	}

	return proxies, nil
}
