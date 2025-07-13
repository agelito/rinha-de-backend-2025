package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

type ServerScore struct {
	Server string
	Score  int
}

type ServerHealthResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

const (
	MaxServerResponseTime = 15_000
	ScoreMultiplier       = 1000
)

var (
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	scores = []ServerScore{
		{
			Server: "http://localhost:8001",
			Score:  0,
		},
		{
			Server: "http://localhost:8002",
			Score:  0,
		},
	}
)

func serverHealthRequest(server string) (*ServerHealthResponse, error) {
	url := fmt.Sprintf("%s/payments/service-health", server)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &ServerHealthResponse{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return &ServerHealthResponse{}, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return &ServerHealthResponse{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	var healthResponse ServerHealthResponse
	if err := json.NewDecoder(res.Body).Decode(&healthResponse); err != nil {
		return &ServerHealthResponse{}, err
	}

	return &healthResponse, err
}

func calculateServerScores() {
	for _, srv := range scores {
		res, err := serverHealthRequest(srv.Server)

		if err != nil {
			log.Printf("server health request failed: %v\n", err)
			srv.Score = 0
			continue
		}

		if !res.Failing {
			minResponseTime := int(math.Max(0, math.Min(float64(res.MinResponseTime), float64(MaxServerResponseTime))))
			mul := 1.0 - float32(minResponseTime)/float32(MaxServerResponseTime)

			srv.Score = int(mul * ScoreMultiplier)
		}
	}
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		calculateServerScores()

		fmt.Printf("scores: %v\n", scores)

		<-ticker.C
	}
}
