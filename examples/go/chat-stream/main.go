package main

import (
	"context"
	"fmt"
	"os"

	"chat-stream/aifunc/chat_stream"
)

func strPtr(s string) *string { return &s }

func main() {

	// config := &chat_stream.AIFuncConfig{
	// 	BaseURL:    "https://your-api-endpoint/v1",
	// 	Model:      "your-model-name",
	// 	APIKey:     "your-api-key",
	// 	MaxRetries: 3,
	// }

	// To use a real model, replace the line below with the commented config above.
	config := &chat_stream.AIFuncConfig{Mock: true}

	inputShort := chat_stream.ChatStreamInput{
		Message: "What is the difference between a process and a thread? Answer in 3 sentences.",
	}

	inputWithContext := chat_stream.ChatStreamInput{
		Message: "Should I prefer threads or processes for CPU-bound work on multi-core machines?",
		Context: strPtr(
			"Conversation history:\n" +
				"User: What is the difference between a process and a thread?\n" +
				"Assistant: Processes have separate memory; threads share an address space.",
		),
	}

	inputLong := chat_stream.ChatStreamInput{
		Message: "Explain the entire history of the internet from ARPANET to today, in detail.",
	}

	// Short reply — run to completion
	fmt.Println("--- short reply (run to completion) ---")
	tokens, errc := chat_stream.ChatStream(context.Background(), config, inputShort)
	for token := range tokens {
		fmt.Print(token)
	}
	if err := <-errc; err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print("\n\n")

	// Follow-up with context
	fmt.Println("--- reply with context ---")
	tokens, errc = chat_stream.ChatStream(context.Background(), config, inputWithContext)
	for token := range tokens {
		fmt.Print(token)
	}
	if err := <-errc; err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print("\n\n")

	// Long reply — cancel after 500 characters
	fmt.Println("--- long reply (cancel after 500 chars) ---")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tokens, errc = chat_stream.ChatStream(ctx, config, inputLong)
	chars := 0
	for token := range tokens {
		fmt.Print(token)
		chars += len(token)
		if chars >= 500 {
			cancel()
			break
		}
	}
	for range tokens {
	}
	fmt.Print("\n[cancelled]\n")
	if err := <-errc; err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
