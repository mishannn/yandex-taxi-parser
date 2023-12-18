package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type StateView struct {
	Stack []Stack `json:"stack"`
}

type RouterResponse struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	TaxiInfo TaxiInfo `json:"taxiInfo"`
}

type TaxiInfo struct {
	CurrencyCode string `json:"currencyCode"`
	Price        int    `json:"price"`
	Time         int    `json:"time"`
	WaitingTime  int    `json:"waitingTime"`
	PriceText    string `json:"priceText"`
	IsSurge      bool   `json:"isSurge"`
	Link         string `json:"link"`
}

type Stack struct {
	RouterResponse RouterResponse `json:"routerResponse"`
}

type Coordinates struct {
	Lat float64
	Lng float64
}

func getMoscowTaxiRoute(ctx context.Context, cli *http.Client, cookie string, from Coordinates, to Coordinates) (*TaxiInfo, error) {
	qs := url.Values{}
	qs.Set("mode", "routes")
	qs.Set("rtt", "taxi")
	qs.Set("rtext", fmt.Sprintf("%f,%f~%f,%f", from.Lat, from.Lng, to.Lat, to.Lng))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://yandex.ru/maps/", nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	req.URL.RawQuery = qs.Encode()

	req.Header.Set("cookie", cookie)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server sent status code %d, but expected 200", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't parse response body: %w", err)
	}

	stateViewHTML, err := doc.Find("script.state-view").Html()
	if err != nil {
		return nil, fmt.Errorf("can't get state view content: %w", err)
	}

	stateViewJSON := html.UnescapeString(stateViewHTML)

	var stateView StateView
	err = json.Unmarshal([]byte(stateViewJSON), &stateView)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal state view: %w", err)
	}

	if len(stateView.Stack) != 1 {
		return nil, fmt.Errorf("unexpected state view stack size: %d", len(stateView.Stack))
	}

	stackItem := stateView.Stack[0]

	if len(stackItem.RouterResponse.Routes) != 1 {
		return nil, fmt.Errorf("unexpected state view routes count: %d", len(stackItem.RouterResponse.Routes))
	}

	return &stackItem.RouterResponse.Routes[0].TaxiInfo, nil
}

func getMoscowTaxiRouteWithProxies(proxyURLs []string, cookie string, from Coordinates, to Coordinates) (*TaxiInfo, error) {
	resultCh := make(chan *TaxiInfo)

	ctx, cancel := context.WithCancel(context.TODO())
	defer func() {
		cancel()
	}()

	go func() {
		var wg sync.WaitGroup

		for _, proxyURL := range proxyURLs {
			proxy := proxyURL

			wg.Add(1)
			go func() {
				defer wg.Done()

				cli := &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
					Transport: &http.Transport{
						Proxy: func(r *http.Request) (*url.URL, error) {
							return url.Parse(proxy)
						},
					},
					Timeout: 30 * time.Second,
				}

				taxiInfo, err := getMoscowTaxiRoute(ctx, cli, cookie, from, to)
				if err != nil {
					return
				}

				resultCh <- taxiInfo
			}()
		}

		wg.Wait()

		close(resultCh)
	}()

	taxiInfo, ok := <-resultCh
	if !ok {
		return nil, fmt.Errorf("no successful requests")
	}

	return taxiInfo, nil
}
