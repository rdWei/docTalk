package main 

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type Message struct {
	Timestamp time.Time
	Sender    string
	Color     string
	Message   string
}

func ReadDocMessages(documentId, serviceAccountFile string) ([]Message, error) {
	ctx := context.Background()

	srv, err := docs.NewService(ctx, option.WithCredentialsFile(serviceAccountFile))
	if err != nil {
		return nil, fmt.Errorf("unable to create docs service: %v", err)
	}

	doc, err := srv.Documents.Get(documentId).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve document: %v", err)
	}

	var lines []string
	for _, c := range doc.Body.Content {
		if c.Paragraph != nil {
			var line strings.Builder
			for _, elem := range c.Paragraph.Elements {
				if elem.TextRun != nil {
					line.WriteString(elem.TextRun.Content)
				}
			}
			trimmed := strings.TrimSpace(line.String())
			if trimmed != "" {
				lines = append(lines, trimmed)
			}
		}
	}

	// Regex per: HH:mm:ss - Sender - Color - Message
	pattern := regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}) - ([^ -]+) - ([^ -]+) - (.+)$`)

	var messages []Message
	for _, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		timestampStr := matches[1]
		sender := matches[2]
		color := matches[3]
		text := matches[4]

		// Parse orario (senza data)
		timestamp, err := time.Parse("15:04:05", timestampStr)
		if err != nil {
			continue
		}

		messages = append(messages, Message{
			Timestamp: timestamp,
			Sender:    sender,
			Color:     color,
			Message:   text,
		})
	}

	return messages, nil
}

