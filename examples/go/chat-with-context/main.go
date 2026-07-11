package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"chat-with-context/aifunc/extract_keywords"
	"chat-with-context/aifunc/generate_reply"
	"chat-with-context/aifunc/recognize_intent"
	"chat-with-context/aifunc/summarize"
)

// var config = &summarize.AIFuncConfig{
// 	BaseURL:    "https://your-api-endpoint/v1",
// 	Model:      "your-model-name",
// 	APIKey:     "your-api-key",
// 	MaxRetries: 3,
// }

// To use a real model, replace the line below with the commented config above.
var config = &summarize.AIFuncConfig{Mock: true}

const (
	window        = 4
	compressAfter = 6
)

var intentLabels = []string{
	"ask_recommendation",
	"ask_logistics",
	"ask_budget",
	"share_preference",
	"confirm",
	"other",
}

// Simulated user turns about planning a trip to Europe
var messages = []string{
	"I'm planning a three-week trip across Europe in September. Where should I start?",
	"I enjoy hiking and local markets. I'd rather skip the big tourist traps.",
	"What's the best way to get from Paris to Barcelona — high-speed train or budget flight?",
	"How much should I budget per day for food and transport in Western Europe?",
	"I've heard the Dolomites are stunning in early autumn. Is it worth a detour from Venice?",
	"Alright, I think I'll do Paris → Barcelona → Rome → Venice → Dolomites. Does that route make sense?",
}

type turn struct {
	role string
	text string
}

// Memory is just plain slices — no special memory object needed.
// history : { role, text } tuples accumulated across turns
// topics  : deduplicated keywords accumulated across all turns
// intents : intent label per user turn, in order
var history []turn
var topics []string
var intents []string
var memorySummary string

func buildContext() string {
	var parts []string

	if memorySummary != "" {
		parts = append(parts, "Earlier in this conversation: "+memorySummary)
	}

	recent := history
	if len(recent) > window {
		recent = recent[len(recent)-window:]
	}
	if len(recent) > 0 {
		var lines []string
		for _, t := range recent {
			role := "Assistant"
			if t.role == "user" {
				role = "User"
			}
			lines = append(lines, "  "+role+": "+t.text)
		}
		parts = append(parts, "Recent conversation:\n"+strings.Join(lines, "\n"))
	}

	if len(topics) > 0 {
		parts = append(parts, "Topics discussed so far: "+strings.Join(topics, ", "))
	}

	if len(intents) > 0 {
		tail := intents
		if len(tail) > 4 {
			tail = tail[len(tail)-4:]
		}
		parts = append(parts, "User intent pattern: "+strings.Join(tail, " → "))
	}

	return strings.Join(parts, "\n\n")
}

func maybeCompress(ctx context.Context) {
	if len(history) <= compressAfter {
		return
	}

	older := history[:len(history)-window]
	var texts []string
	for _, t := range older {
		texts = append(texts, t.text)
	}
	olderText := strings.Join(texts, " ")

	maxLen := 40
	result, err := summarize.Summarize(ctx, config, summarize.SummarizeInput{
		Text:      olderText,
		MaxLength: &maxLen,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "compress error:", err)
		return
	}

	memorySummary = result.Summary
	history = history[len(history)-window:]
	fmt.Printf("  [memory compressed → \"%s\"]\n\n", memorySummary)
}

func containsTopic(word string) bool {
	for _, t := range topics {
		if t == word {
			return true
		}
	}
	return false
}

func main() {
	if config.Mock {
		fmt.Println(
			"This example requires a real LLM to produce meaningful results.\n" +
				"Mock mode cannot simulate intent-aware replies grounded in conversation history.\n" +
				"\n" +
				"To run this example, replace the line:\n" +
				"\n" +
				"  // Mock: true,\n" +
				"\n" +
				"with real credentials:\n" +
				"\n" +
				"  config = &summarize.AIFuncConfig{\n" +
				"      BaseURL: \"https://your-api-endpoint/v1\",\n" +
				"      Model:   \"your-model-name\",\n" +
				"      APIKey:  \"your-api-key\",\n" +
				"  }\n",
		)
		os.Exit(0)
	}

	ctx := context.Background()

	for i, userMsg := range messages {
		fmt.Printf("[Turn %d] User: %s\n", i+1, userMsg)

		// 1. Classify the user's intent
		intentResult, err := recognize_intent.RecognizeIntent(ctx, config, recognize_intent.RecognizeIntentInput{
			Text:    userMsg,
			Intents: intentLabels,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "recognize-intent error:", err)
			os.Exit(1)
		}
		intents = append(intents, intentResult.Intent)

		// 2. Extract keywords and accumulate into the topics slice
		maxKw := 3
		kwResult, err := extract_keywords.ExtractKeywords(ctx, config, extract_keywords.ExtractKeywordsInput{
			Text:        userMsg,
			MaxKeywords: &maxKw,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-keywords error:", err)
			os.Exit(1)
		}
		for _, kw := range kwResult.Keywords {
			if word, ok := kw["word"].(string); ok && !containsTopic(word) {
				topics = append(topics, word)
			}
		}

		// 3. Append user turn to history slice
		history = append(history, turn{role: "user", text: userMsg})

		// 4. Compress old history before replying if it has grown too long
		maybeCompress(ctx)

		// 5. Build context from memory slices and generate a reply
		ctxStr := buildContext()
		tone := "friendly"
		replyResult, err := generate_reply.GenerateReply(ctx, config, generate_reply.GenerateReplyInput{
			Message: userMsg,
			Tone:    &tone,
			Context: &ctxStr,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "generate-reply error:", err)
			os.Exit(1)
		}

		// 6. Append assistant reply to history slice
		history = append(history, turn{role: "assistant", text: replyResult.Reply})

		fmt.Printf("         Intent  : %s (%.0f%%)\n", intentResult.Intent, intentResult.Confidence*100)
		fmt.Printf("         Topics  : %v\n", topics)
		fmt.Printf("         Reply   : %s\n", replyResult.Reply)
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Final memory state")
	fmt.Printf("  topics  : %v\n", topics)
	fmt.Printf("  intents : %v\n", intents)
	if memorySummary != "" {
		fmt.Printf("  summary : %s\n", memorySummary)
	}
	fmt.Printf("  history : %d turns in window\n", len(history))
}
