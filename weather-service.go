package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"
)

type Cache = map[string]*ServiceResult

type WeatherService struct {
    cache Cache
}

func MakeWeatherService() *WeatherService {
    return &WeatherService{ cache: Cache{} }
}

var WeatherCodes = map[int]string{
    0: "Clear sky",
    1: "Mainly clear",
    2: "Partly cloudy",
    3: "Overcast",
    45: "Fog",
    48: "Depositing rime fog",
    51: "Light drizzle",
    53: "Moderate drizzle",
    55: "Dense drizzle",
    56: "Light freezing drizzle",
    57: "Dense freezing drizzle",
    61: "Slight rain",
    63: "Moderate rain",
    65: "Heavy rain",
    66: "Light freezing rain",
    67: "Heavy freezing rain",
    71: "Slight snow fall",
    73: "Moderate snow fall",
    75: "Heavy snow fall",
    77: "Snow grains",
    80: "Slight rain showers",
    81: "Moderate rain showers",
    82: "Violent rain showers",
    85: "Slight snow showers",
    86: "Heavy snow showers",
    95: "Thunderstorm",
    96: "Thunderstorm with slight hail",
    99: "Thunderstorm with heavy hail",
}

type WeatherRequest struct {
	Latitude  string `json:"Latitude"`
	Longitude string `json:"Longitude"`
}

func (r WeatherRequest) hash() string {
    return r.Latitude + r.Longitude
}

func (ws *WeatherService) getRandomKey() string {
    n := rand.IntN(len(ws.cache))
    i := 0
    for key := range ws.cache {
        if i == n {
            return key
        }
        i++
    }
    return ""
}

func (ws *WeatherService) cacheRequest(req WeatherRequest, res *ServiceResult) {
    const MAX_LEN = 100
    if len(ws.cache) == MAX_LEN {
        delete(ws.cache, ws.getRandomKey())
    }
    ws.cache[req.hash()] = res
}

type MeteoHourly struct {
	Temperature   []float64 `json:"temperature_2m"`
	Humidity      []int     `json:"relative_humidity_2m"`
	WindSpeed     []float64 `json:"windspeed_10m"`
	WindDirection []int     `json:"winddirection_10m"`
	WeatherCode   []int     `json:"weathercode"`
}

type MeteoResponse struct {
	Latitude        float64     `json:"latitude"`
	Longitude       float64     `json:"longitude"`
	WeatherElements MeteoHourly `json:"hourly"`
}

type CustomWeatherElement struct {
	Humidity      int     `json:"Humidity"`
	Weather       string  `json:"Weather"`
	WindDirection string  `json:"WindDirection"`
	WindSpeed     float64 `json:"WindSpeed"`
	TempCelcius   float64 `json:"Temperature"`
}

func parseWindDirection(windDir int) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	i := int(float64(windDir / 45.0)) % 8
	return dirs[i]
}

func transformMeteoElements(meteoResp MeteoResponse) []CustomWeatherElement {
	arrayLen := len(meteoResp.WeatherElements.Humidity)

	ret := make([]CustomWeatherElement, 0, arrayLen)
    
	for i := range arrayLen {
		customElem := CustomWeatherElement{
			Humidity:      meteoResp.WeatherElements.Humidity[i],
			Weather:       WeatherCodes[meteoResp.WeatherElements.WeatherCode[i]],
			WindDirection: parseWindDirection(meteoResp.WeatherElements.WindDirection[i]),
			WindSpeed:     meteoResp.WeatherElements.WindSpeed[i],
			TempCelcius:   meteoResp.WeatherElements.Temperature[i],
		}
		ret = append(ret, customElem)
	}
	return ret
}

type CustomWeatherResponse struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64  `json:"Longitude"`
	WeatherElements []CustomWeatherElement `json:"TimeSeries"` 
}

func (ws *WeatherService) Run(jsonReq RequestJson) (ret ServiceResult) {
    start := time.Now()
    defer func() {
        fmt.Println("Request took ", time.Since(start))
    }()

	var req WeatherRequest
	json.Unmarshal(jsonReq.JsonBody, &req)

    cachedRes, present := ws.cache[req.hash()]
    if present {
        return *cachedRes
    } else {
        defer func() {
            ws.cacheRequest(req, &ret)
        }()
    }

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%v&longitude=%v&hourly=temperature_2m,relative_humidity_2m,windspeed_10m,winddirection_10m,weathercode",
		req.Latitude, req.Longitude)
    
	resp, err := http.Get(url) 
	if err != nil {
		ret.Code = 500
		ret.Error = "Failed to access Meteo API"
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ret.Code = 500
		ret.Error = "Failed to read HTTP body from Meteo"
		return
	}

	var meteoResp MeteoResponse
	err = json.Unmarshal(body, &meteoResp)
	if err != nil {
		ret.Code = 500
		ret.Error = "Failed to parse Meteo Response"
		return
	}

    ret.Code = 200
    ret.Error = ""
	ret.Yield = CustomWeatherResponse{
		Latitude:        meteoResp.Latitude,
		Longitude:       meteoResp.Longitude,
		WeatherElements: transformMeteoElements(meteoResp),
	}
	return
}