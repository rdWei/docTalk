package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// AppendMessageToDoc aggiunge un messaggio formattato come
// "HH:mm:ss - sender - color - message" al Google Doc specificato.
func AppendMessageToDoc(documentId, serviceAccountFile, sender, color, message string) error {
	ctx := context.Background()

	srv, err := docs.NewService(ctx, option.WithCredentialsFile(serviceAccountFile))
	if err != nil {
		return fmt.Errorf("unable to create docs service: %w", err)
	}

	doc, err := srv.Documents.Get(documentId).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve document: %w", err)
	}

	content := doc.Body.Content
	if len(content) == 0 {
		return fmt.Errorf("document is empty")
	}

	endIndex := content[len(content)-1].EndIndex
	if endIndex == 0 {
		return fmt.Errorf("invalid end index")
	}

	// Orario nel formato HH:mm:ss
	now := time.Now().Format("15:04:05")

	textToAppend := fmt.Sprintf("%s - %s - %s - %s\n", now, sender, color, message)

	requests := []*docs.Request{
		{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{
					Index: endIndex - 1,
				},
				Text: textToAppend,
			},
		},
	}

	_, err = srv.Documents.BatchUpdate(documentId, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to append message: %w", err)
	}

	return nil
}

