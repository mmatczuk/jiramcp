package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mmatczuk/jiramcp/internal/jirahttp"
	"github.com/mmatczuk/jiramcp/internal/jiramcp"
)

func main() {
	client, err := jirahttp.New(jirahttp.Config{
		URL:        requireEnv("JIRA_URL"),
		Email:      requireEnv("JIRA_EMAIL"),
		APIToken:   requireEnv("JIRA_API_TOKEN"),
		MaxRetries: 3,
		BaseDelay:  time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create JIRA client: %v", err)
	}

	s := jiramcp.NewServer(client)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "Required environment variable %s is not set\n", key)
		os.Exit(1)
	}
	return v
}
