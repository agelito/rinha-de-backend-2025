package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
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

	log.Debug("service-health", "server", server, "response", healthResponse)

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

	log.Debug("server score", "server", srv, "score", score)

	return score
}

func checkServerScoreRepeating(nc *nats.Conn, srv string, freq time.Duration, ch chan *servers.ServerScore) {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()

	for {
		score := calculateServerScore(srv)

		scoreMsg := &servers.ServerScore{
			Server: srv,
			Score:  score,
		}

		msgBytes, err := proto.Marshal(scoreMsg)

		if err != nil {
			log.Error("could not serialize score message", "server", srv, "error", err)
			continue
		}

		if err := nc.Publish(subjects.SubjectServerScore, msgBytes); err != nil {
			log.Error("could not publish server score message", "server", srv, "error", err)
		}

		ch <- scoreMsg

		<-ticker.C
	}
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	serverList := []string{
		"http://localhost:8001",
		"http://localhost:8002",
	}

	ch := make(chan *servers.ServerScore)

	for _, srv := range serverList {
		go checkServerScoreRepeating(nc, srv, CalculateScoreFrequency*time.Second, ch)
	}

	scores := make(map[string]int32)

	for {
		updatedScore := <-ch

		scores[updatedScore.Server] = updatedScore.Score

		var sortedScores []*servers.ServerScore

		for server, score := range scores {
			sortedScores = append(sortedScores, &servers.ServerScore{
				Server: server,
				Score:  score,
			})
		}

		sort.Slice(sortedScores, func(i, j int) bool {
			return sortedScores[i].Score > sortedScores[j].Score
		})

		serverRanking := &servers.ServerRanking{
			Servers: sortedScores,
		}

		msgData, err := proto.Marshal(serverRanking)

		if err != nil {
			log.Error("could not serialize server ranking", "error", err)
			continue
		}

		if err := nc.Publish(subjects.SubjectServerRanking, msgData); err != nil {
			log.Error("could not publish server rankings message", "error", err)
		}

		log.Info("server rankings", "servers", sortedScores)
	}
}
