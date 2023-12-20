package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getCookies(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("server sent status code %d, but expected %d", resp.StatusCode, 200)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read response body: %w", err)
	}

	return strings.TrimSpace(string(respBody)), nil
}
