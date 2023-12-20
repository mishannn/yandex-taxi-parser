package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeatherErrorResponse struct {
	Cod     int64  `json:"cod"`
	Message string `json:"message"`
}

type WeatherResponse struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int64   `json:"timezone_offset"`
	Current        Current `json:"current"`
}

type Current struct {
	Dt        int64         `json:"dt"`
	Temp      float64       `json:"temp"`
	FeelsLike float64       `json:"feels_like"`
	Pressure  float64       `json:"pressure"`
	Humidity  float64       `json:"humidity"`
	WindSpeed float64       `json:"wind_speed"`
	WindGust  float64       `json:"wind_gust,omitempty"`
	WindDeg   float64       `json:"wind_deg"`
	Weather   []Weather     `json:"weather"`
	Rain      Precipitation `json:"rain,omitempty"`
	Snow      Precipitation `json:"snow,omitempty"`
}

type Precipitation struct {
	The1H float64 `json:"1h"`
}

type Weather struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func GetWeather(apikey string, lat float64, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f6&lon=%f&appid=%s&exclude=minutely,hourly,daily,alerts", lat, lon, apikey)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	if res.StatusCode != 200 {
		var data WeatherErrorResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, fmt.Errorf("can't parse response error body: %w", err)
		}

		return nil, fmt.Errorf("server returned error: %s", data.Message)
	}

	var data WeatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("can't parse response body: %w", err)
	}

	return &data, nil
}

func ConvertKelvinToCelsius(deg float64) float64 {
	return deg - 273.15
}
