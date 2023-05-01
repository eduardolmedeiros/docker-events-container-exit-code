package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var exitCodes = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "docker_container_exit_codes",
		Help: "Exit codes of Docker containers",
	},
	[]string{"container_id", "container_name"},
)

func init() {
	prometheus.MustRegister(exitCodes)
}

func listenToDockerEvents() {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	events, errs := dockerClient.Events(ctx, types.EventsOptions{})
	for {
		select {
		case event := <-events:
			if event.Type == "container" && event.Action == "die" {
				exitCodeStr := event.Actor.Attributes["exitCode"]
				containerID := event.Actor.ID
				containerName := event.Actor.Attributes["name"]
				exitCode, err := strconv.ParseFloat(exitCodeStr, 64)
				if err != nil {
					log.Printf("Failed to parse exit code for container %s (%s): %s", containerID, containerName, err.Error())
				} else {
					exitCodes.WithLabelValues(containerID, containerName).Set(exitCode)
				}
			}
		case err := <-errs:
			log.Fatal(err)
		}
	}
}

func main() {
	go listenToDockerEvents()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
