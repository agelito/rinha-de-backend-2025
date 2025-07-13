package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/agelito/rinha-de-backend-2025/messages/model/servers"
	"github.com/agelito/rinha-de-backend-2025/messages/subjects"
	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type ServerScore struct {
	Server string
	Score  int32
}

type ServerHealthResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

const (
	CalculateScoreFrequency = 5
	MaxServerResponseTime   = 10_000
	ScoreMultiplier         = 1_000
)

var (
	client = &http.Client{
		Timeout: CalculateScoreFrequency * time.Second,
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

	log.Info("service-health", "server", server, "response", healthResponse)

	return &healthResponse, err
}

func calculateServerScore(srv string) int32 {
	res, err := serverHealthRequest(srv)

	if err != nil {
		log.Error("service-health request failed", "server", srv, "error", err)
		return 0
	}

	score := int32(0)

	if !res.Failing {
		minResponseTime := int(math.Max(0, math.Min(float64(res.MinResponseTime), float64(MaxServerResponseTime))))
		mul := 1.0 - float32(minResponseTime)/float32(MaxServerResponseTime)

		score = int32(mul * ScoreMultiplier)
	}

	log.Info("server score", "server", srv, "score", score)

	return score
}

func checkServerScoreRepeating(nc *nats.Conn, srv string, freq time.Duration) {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()

	for {
		score := calculateServerScore(srv)

		scoreMsg := servers.ServerScore{
			Server: srv,
			Score:  score,
		}

		msgBytes, err := proto.Marshal(&scoreMsg)

		if err != nil {
			log.Error("could not serialize score message", "server", srv, "error", err)
			continue
		}

		if err := nc.Publish(subjects.SubjectServerScore, msgBytes); err != nil {
			log.Error("could not publish server score message", "server", srv, "error", err)
		}

		<-ticker.C
	}
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	ticker := time.NewTicker(CalculateScoreFrequency * time.Second)
	defer ticker.Stop()

	servers := []string{
		"http://localhost:8001",
		"http://localhost:8002",
	}

	for _, srv := range servers {
		go checkServerScoreRepeating(nc, srv, CalculateScoreFrequency*time.Second)
	}

	select {}
}
