package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"all-packages-stream/aifunc/answer_stream"
	"all-packages-stream/aifunc/article_stream"
	"all-packages-stream/aifunc/chat_stream"
	"all-packages-stream/aifunc/explain_stream"
	"all-packages-stream/aifunc/review_stream"
	"all-packages-stream/aifunc/translate_stream"
	"all-packages-stream/aifunc/write_stream"
)

// var config = &answer_stream.AIFuncConfig{
// 	BaseURL:    "https://your-api-endpoint/v1",
// 	Model:      "your-model-name",
// 	APIKey:     "your-api-key",
// 	MaxRetries: 3,
// }

// To use a real model, replace the line below with the commented config above.
var config = &answer_stream.AIFuncConfig{Mock: true}

func section(title string) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("  %s\n", title)
	fmt.Printf("%s\n", strings.Repeat("=", 60))
}

func strPtr(s string) *string { return &s }
func intPtr(n int) *int       { return &n }

// streamPrint consumes a token channel and writes each chunk to stdout immediately.
func streamPrint(tokens <-chan string, errc <-chan error) error {
	for token := range tokens {
		fmt.Print(token)
	}
	fmt.Println()
	return <-errc
}

func main() {
	if config.Mock {
		fmt.Println("Notice: You are using mock mode for offline testing. " +
			"Configure a real model for the full experience. Continuing with mock responses...")
	}

	ctx := context.Background()

	article := "In 1915, Albert Einstein published the General Theory of Relativity, " +
		"fundamentally transforming our understanding of physics. The theory posits " +
		"that gravity is not an invisible force, but rather a curvature of spacetime " +
		"caused by the presence of mass and energy. This groundbreaking framework " +
		"revolutionized modern science and introduced the famous equation E=mc²."

	codeSnippet := `def fetch_user(user_id):
    conn = get_connection()
    result = conn.execute(f"SELECT * FROM users WHERE id = {user_id}")
    return result.fetchone()
`

	// ─── Conversational & Q&A ─────────────────────────────────────────

	section("1. CHAT STREAM")
	fmt.Println("User: Explain goroutines and channels in Go in 3 sentences.")
	fmt.Println()
	fmt.Print("Assistant: ")
	if err := streamPrint(chat_stream.ChatStream(ctx, config, chat_stream.ChatStreamInput{
		Messages: []chat_stream.Message{
			{Role: "user", Content: "Explain goroutines and channels in Go in 3 sentences."},
		},
	})); err != nil {
		fmt.Fprintln(os.Stderr, "chat-stream:", err)
		os.Exit(1)
	}

	section("2. ANSWER STREAM (with context / RAG)")
	contextText := "AIFunc is a function-based AI toolkit. Developers declare the packages they need " +
		"in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. " +
		"Each package supports a mock mode for testing without consuming API credits. " +
		"Streaming packages return tokens incrementally via channels."
	question := "How does AIFunc support offline testing, and what do streaming packages return?"
	fmt.Printf("Q: %s\n\n", question)
	fmt.Print("A: ")
	if err := streamPrint(answer_stream.AnswerStream(ctx, config, answer_stream.AnswerStreamInput{
		Question: question,
		Context:  strPtr(contextText),
		Depth:    strPtr("concise"),
		Audience: strPtr("technical"),
	})); err != nil {
		fmt.Fprintln(os.Stderr, "answer-stream:", err)
		os.Exit(1)
	}

	section("3. EXPLAIN STREAM")
	fmt.Println("Topic: Go's garbage collector")
	fmt.Println()
	if err := streamPrint(explain_stream.ExplainStream(ctx, config, explain_stream.ExplainStreamInput{
		Topic:    "Go's garbage collector",
		Audience: strPtr("intermediate"),
		Depth:    strPtr("standard"),
	})); err != nil {
		fmt.Fprintln(os.Stderr, "explain-stream:", err)
		os.Exit(1)
	}

	// ─── Long-form writing ────────────────────────────────────────────

	section("4. ARTICLE STREAM")
	title := "Why Typed AI Functions Beat Ad-Hoc Prompt Scripts"
	outline := "- The cost of untyped prompt glue code\n" +
		"- How function-shaped AI APIs improve testability\n" +
		"- Streaming vs batch for product UX\n" +
		"- Practical adoption tips"
	fmt.Printf("Title  : %s\n", title)
	fmt.Printf("Outline: %s\n\n", outline)
	if err := streamPrint(article_stream.ArticleStream(ctx, config, article_stream.ArticleStreamInput{
		Title:     title,
		Outline:   strPtr(outline),
		Style:     strPtr("informational"),
		Audience:  strPtr("developers"),
		WordCount: intPtr(250),
	})); err != nil {
		fmt.Fprintln(os.Stderr, "article-stream:", err)
		os.Exit(1)
	}

	section("5. WRITE STREAM")
	prompt := "Write a short internal proposal recommending that our team adopt AIFunc " +
		"for customer-support reply generation."
	structure := "1. Problem\n" +
		"2. Proposed approach\n" +
		"3. Expected benefits\n" +
		"4. Next steps"
	fmt.Printf("Prompt   : %s\n", prompt)
	fmt.Printf("Structure: %s\n\n", structure)
	if err := streamPrint(write_stream.WriteStream(ctx, config, write_stream.WriteStreamInput{
		Prompt:    prompt,
		Format:    strPtr("proposal"),
		Structure: strPtr(structure),
		Tone:      strPtr("professional"),
		Audience:  strPtr("engineers"),
		WordCount: intPtr(300),
	})); err != nil {
		fmt.Fprintln(os.Stderr, "write-stream:", err)
		os.Exit(1)
	}

	// ─── Translation & review ─────────────────────────────────────────

	section("6. TRANSLATE STREAM")
	fmt.Printf("Original (EN):\n%s\n\n", article)
	fmt.Println("Translation (zh-CN):")
	fmt.Println()
	if err := streamPrint(translate_stream.TranslateStream(ctx, config, translate_stream.TranslateStreamInput{
		Text:       article,
		TargetLang: "zh-CN",
		Style:      strPtr("natural"),
		Domain:     strPtr("technical"),
	})); err != nil {
		fmt.Fprintln(os.Stderr, "translate-stream:", err)
		os.Exit(1)
	}

	section("7. REVIEW STREAM")
	fmt.Printf("Code under review:\n%s\n", codeSnippet)
	fmt.Println("Findings:")
	fmt.Println()
	if err := streamPrint(review_stream.ReviewStream(ctx, config, review_stream.ReviewStreamInput{
		Content:        codeSnippet,
		Type:           strPtr("code"),
		Language:       strPtr("Python"),
		Focus:          strPtr("correctness, security"),
		Context:        strPtr("Simple data-access helper in a web API."),
		Severity:       strPtr("all"),
		OutputLanguage: strPtr("English"),
	})); err != nil {
		fmt.Fprintln(os.Stderr, "review-stream:", err)
		os.Exit(1)
	}

	if config.Mock {
		fmt.Println("Notice: You are using mock mode for offline testing. " +
			"Configure a real model for the full experience.")
	}
}
