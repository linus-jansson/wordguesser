package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// After a successful game, we send the results to a webhook for easier reading
// DiscordEmbedPayload represents the payload structure for Discord webhook
type DiscordEmbedPayload struct {
	Embeds []Embed `json:"embeds"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// SendFileToDiscordWebhook reads a file and sends its content to a Discord webhook as an embed
func SendFileToDiscordWebhook(filePath, webhookURL string) error {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	payload := DiscordEmbedPayload{
		Embeds: []Embed{
			{
				Title:       "Här är dagens ordel!",
				Description: string(contentBytes),
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func getDiscordWebhookURL() string {
	return os.Getenv("DISCORD_WEBHOOK_URL")
}
