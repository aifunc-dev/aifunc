package main

import (
	"context"
	"fmt"
	"os"

	"hello-aifunc/aifunc/summarize"
)

func main() {

	// config := &summarize.AIFuncConfig{
	// 	BaseURL:    "https://your-api-endpoint/v1",
	// 	Model:      "your-model-name",
	// 	APIKey:     "your-api-key",
	// 	MaxRetries: 3,
	// }

	// To use a real model, replace the line below with the commented config above.
	config := &summarize.AIFuncConfig{Mock: true}

	if config.Mock {
		fmt.Println("Notice: You are using mock mode for offline testing. " +
			"Configure a real model for the full experience. Continuing with mock responses...")
	}

	text := "The James Webb Space Telescope captured its first full-color images in July 2022, " +
		"revealing thousands of galaxies in a patch of sky smaller than a grain of sand held " +
		"at arm's length. The images show galaxies as they appeared over 13 billion years ago, " +
		"providing a glimpse into the early universe shortly after the Big Bang."

	maxLen := 30
	result, err := summarize.Summarize(context.Background(), config, summarize.SummarizeInput{
		Text:      text,
		MaxLength: &maxLen,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println("Original  :", text)
	fmt.Println("Summary   :", result.Summary)
	fmt.Println("Word count:", result.WordCount)
}
