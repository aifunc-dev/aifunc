package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"customer-support/aifunc/analyze_sentiment"
	"customer-support/aifunc/extract_json"
	"customer-support/aifunc/recognize_intent"
)

// config := &analyze_sentiment.AIFuncConfig{
// 	BaseURL: "https://your-api-endpoint/v1",
// 	Model:   "your-model-name",
// 	APIKey:  "your-api-key",
// 	max_retries=3,
// }

// To run this example, replace the mock config below with real credentials.
var config = &analyze_sentiment.AIFuncConfig{Mock: true}

// Tip: Each call accepts its own config — use cheaper models for simple tasks to save cost.
//
// cheap := &analyze_sentiment.AIFuncConfig{BaseURL: "...", Model: "...", APIKey: "..."}
// strong := &analyze_sentiment.AIFuncConfig{BaseURL: "...", Model: "...", APIKey: "..."}
//
// analyze_sentiment.AnalyzeSentiment(ctx, cheap, ...)   // classification is simple, cheap model is fine
// extract_json.ExtractJson(ctx, strong, ...)            // extraction needs accuracy, use a stronger model

var customerMessages = []string{
	"What the hell?! I ordered this a WEEK ago and it still hasn't shipped! I want my money back NOW!",
	"Hi, I'd like to check on my order #ORD-20240601-123. It's been three days with no shipping update.",
	"Your stupid app crashed again and I lost all my data! Fix it or I'm leaving!",
	"I was charged twice this month. Transaction IDs: TXN-88201 and TXN-88202. Please help.",
	"It would be cool if you added a dark mode. The bright screen hurts my eyes at night.",
	"How do I export my purchase history to CSV? I can't find the option.",
	"I am SO FURIOUS! Your delivery guy threw my package over the fence and it's destroyed! I want a manager NOW!",
	"Any ongoing promotions for loyal customers? I've been a member for 2 years.",
}

var intentLabels = []string{
	"query_order",
	"request_refund",
	"technical_support",
	"billing_issue",
	"feature_request",
	"general_inquiry",
}

func extractField(result extract_json.ExtractJsonOutput, name string) string {
	if v, ok := result.Extracted[name]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func main() {
	if config.Mock {
		fmt.Println(
			"This example requires a real LLM to produce meaningful results.\n" +
				"Mock mode cannot simulate multi-step reasoning (sentiment → intent → extraction).\n" +
				"\n" +
				"To run this example, replace the line:\n" +
				"\n" +
				"  // Mock: true,\n" +
				"\n" +
				"with real credentials:\n" +
				"\n" +
				"  config = &AIFuncConfig{\n" +
				"      BaseURL: \"https://your-api-endpoint/v1\",\n" +
				"      Model:   \"your-model-name\",\n" +
				"      APIKey:  \"your-api-key\",\n" +
				"  }\n",
		)
		os.Exit(0)
	}

	ctx := context.Background()
	message := customerMessages[rand.Intn(len(customerMessages))]
	fmt.Printf("Customer: %s\n\n", message)

	// Step 1: Sentiment analysis
	sentiment, err := analyze_sentiment.AnalyzeSentiment(ctx, config, analyze_sentiment.AnalyzeSentimentInput{
		Text:   message,
		Labels: []string{"angry", "frustrated", "neutral", "happy", "other"},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "analyze-sentiment error:", err)
		os.Exit(1)
	}
	fmt.Printf("Sentiment: %s (%.0f%%)\n", sentiment.Label, sentiment.Confidence*100)

	if sentiment.Label == "angry" && sentiment.Confidence > 0.7 {
		fmt.Println("\n=> call_human_agent(message, priority=\"HIGH\")")
		return
	}

	// Step 2: Intent recognition
	intentResult, err := recognize_intent.RecognizeIntent(ctx, config, recognize_intent.RecognizeIntentInput{
		Text:    message,
		Intents: intentLabels,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "recognize-intent error:", err)
		os.Exit(1)
	}
	intent := intentResult.Intent
	fmt.Printf("Intent: %s (%.0f%%)\n", intent, intentResult.Confidence*100)

	// Step 3: Route by intent
	switch intent {
	case "query_order":
		info, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
			Text: message,
			Fields: []map[string]any{
				{"name": "order_id", "description": "Order number", "type": "string"},
				{"name": "issue", "description": "What the customer wants to know", "type": "string"},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-json error:", err)
			os.Exit(1)
		}
		fmt.Printf("\n=> query_order_system(order_id=\"%s\", issue=\"%s\")\n",
			extractField(info, "order_id"), extractField(info, "issue"))

	case "request_refund":
		info, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
			Text: message,
			Fields: []map[string]any{
				{"name": "order_id", "description": "Order number", "type": "string"},
				{"name": "reason", "description": "Reason for refund", "type": "string"},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-json error:", err)
			os.Exit(1)
		}
		fmt.Printf("\n=> submit_refund(order_id=\"%s\", reason=\"%s\")\n",
			extractField(info, "order_id"), extractField(info, "reason"))

	case "technical_support":
		info, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
			Text: message,
			Fields: []map[string]any{
				{"name": "issue", "description": "Technical problem", "type": "string"},
				{"name": "platform", "description": "Device or platform", "type": "string"},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-json error:", err)
			os.Exit(1)
		}
		fmt.Printf("\n=> create_tech_ticket(issue=\"%s\", platform=\"%s\")\n",
			extractField(info, "issue"), extractField(info, "platform"))

	case "billing_issue":
		info, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
			Text: message,
			Fields: []map[string]any{
				{"name": "transaction_id", "description": "Transaction ID", "type": "string"},
				{"name": "problem", "description": "Billing problem", "type": "string"},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-json error:", err)
			os.Exit(1)
		}
		fmt.Printf("\n=> flag_billing_dispute(transaction_id=\"%s\", problem=\"%s\")\n",
			extractField(info, "transaction_id"), extractField(info, "problem"))

	case "feature_request":
		info, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
			Text: message,
			Fields: []map[string]any{
				{"name": "feature", "description": "Requested feature", "type": "string"},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-json error:", err)
			os.Exit(1)
		}
		fmt.Printf("\n=> log_feature_request(feature=\"%s\")\n", extractField(info, "feature"))

	default:
		info, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
			Text: message,
			Fields: []map[string]any{
				{"name": "question", "description": "Customer's question", "type": "string"},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "extract-json error:", err)
			os.Exit(1)
		}
		fmt.Printf("\n=> send_to_faq_system(question=\"%s\")\n", extractField(info, "question"))
	}
}
