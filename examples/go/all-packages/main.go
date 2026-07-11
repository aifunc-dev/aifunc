package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"all-packages/aifunc/analyze_sentiment"
	"all-packages/aifunc/answer_question"
	"all-packages/aifunc/classify"
	"all-packages/aifunc/detect_language"
	"all-packages/aifunc/extract_entities"
	"all-packages/aifunc/extract_json"
	"all-packages/aifunc/extract_keywords"
	"all-packages/aifunc/generate_email"
	"all-packages/aifunc/generate_post"
	"all-packages/aifunc/generate_reply"
	"all-packages/aifunc/generate_slug"
	"all-packages/aifunc/generate_title"
	"all-packages/aifunc/recognize_intent"
	"all-packages/aifunc/rewrite"
	"all-packages/aifunc/score_quality"
	"all-packages/aifunc/summarize"
	"all-packages/aifunc/translate"
)

// var config = &summarize.AIFuncConfig{
// 	BaseURL:    "https://your-api-endpoint/v1",
// 	Model:      "your-model-name",
// 	APIKey:     "your-api-key",
// 	MaxRetries: 3,
// }

// To use a real model, replace the line below with the commented config above.
var config = &summarize.AIFuncConfig{Mock: true}

func section(title string) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("  %s\n", title)
	fmt.Printf("%s\n", strings.Repeat("=", 60))
}

func strPtr(s string) *string { return &s }
func intPtr(n int) *int       { return &n }
func boolPtr(b bool) *bool    { return &b }

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

	// ── Easy: single input → single output ────────────────────────────

	section("1. DETECT LANGUAGE")
	for _, text := range []string{
		"The quick brown fox jumps over the lazy dog.",
		"Der schnelle braune Fuchs springt über den faulen Hund.",
		"Le renard brun rapide saute par-dessus le chien paresseux.",
		"El veloz zorro marrón salta sobre el perro perezoso.",
	} {
		r, err := detect_language.DetectLanguage(ctx, config, detect_language.DetectLanguageInput{Text: text})
		if err != nil {
			fmt.Fprintln(os.Stderr, "detect-language:", err)
			os.Exit(1)
		}
		preview := text
		if len(preview) > 40 {
			preview = preview[:40]
		}
		fmt.Printf("  [%s] %s (conf: %.0f%%)  \"%s\"\n", r.Language, r.LanguageName, r.Confidence*100, preview)
	}

	section("2. GENERATE SLUG")
	sr, err := generate_slug.GenerateSlug(ctx, config, generate_slug.GenerateSlugInput{
		Title: "10 Practical Tips for Writing Faster Python Code", Language: strPtr("en"),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate-slug:", err)
		os.Exit(1)
	}
	fmt.Println("Title : 10 Practical Tips for Writing Faster Python Code")
	fmt.Printf("Slug  : %s\nMeta  : %s\nTags  : %v\n", sr.Slug, sr.MetaDescription, sr.Tags)

	section("3. SUMMARIZE")
	sumR, err := summarize.Summarize(ctx, config, summarize.SummarizeInput{Text: article, MaxLength: intPtr(30)})
	if err != nil {
		fmt.Fprintln(os.Stderr, "summarize:", err)
		os.Exit(1)
	}
	fmt.Printf("Summary   : %s\nWord count: %d\n", sumR.Summary, sumR.WordCount)

	section("4. TRANSLATE")
	tr, err := translate.Translate(ctx, config, translate.TranslateInput{
		Text: "The meeting has been moved to Friday at 3 PM.", TargetLang: "es",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "translate:", err)
		os.Exit(1)
	}
	fmt.Println("Original : The meeting has been moved to Friday at 3 PM.")
	fmt.Printf("Spanish  : %s\nDetected : %s\n", tr.Translation, tr.SourceLang)

	section("5. REWRITE")
	orig := "hey, just wanna let u know the deploy went fine, no issues at all"
	rwr, err := rewrite.Rewrite(ctx, config, rewrite.RewriteInput{Text: orig, Style: "formal"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "rewrite:", err)
		os.Exit(1)
	}
	fmt.Printf("Casual : %s\nFormal : %s\n", orig, rwr.Rewritten)

	// ── Medium: structured output or multiple parameters ──────────────

	section("6. GENERATE TITLE")
	titleContent := "This guide covers how to use Docker and GitHub Actions to automate " +
		"testing and deployment of a Node.js application to a cloud server."
	titR, err := generate_title.GenerateTitle(ctx, config, generate_title.GenerateTitleInput{
		Content: titleContent, Style: strPtr("seo"), Count: intPtr(4),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate-title:", err)
		os.Exit(1)
	}
	fmt.Printf("Content: %s\nTitles:\n", titleContent)
	for i, t := range titR.Titles {
		fmt.Printf("  %d. %s\n", i+1, t)
	}

	section("7. EXTRACT KEYWORDS")
	kwR, err := extract_keywords.ExtractKeywords(ctx, config, extract_keywords.ExtractKeywordsInput{
		Text: article, MaxKeywords: intPtr(5),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "extract-keywords:", err)
		os.Exit(1)
	}
	fmt.Println("Keywords from article:")
	for _, kw := range kwR.Keywords {
		word, _ := kw["word"].(string)
		relevance, _ := kw["relevance"].(float64)
		fmt.Printf("  %-30s relevance: %.2f\n", word, relevance)
	}

	section("8. ANALYZE SENTIMENT")
	for _, text := range []string{
		"The product arrived on time and works perfectly. Very happy!",
		"Terrible experience. The package was damaged and support ignored my emails.",
		"Item received. Does what it says.",
	} {
		r, err := analyze_sentiment.AnalyzeSentiment(ctx, config, analyze_sentiment.AnalyzeSentimentInput{
			Text: text, Labels: []string{"positive", "negative", "neutral"},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "analyze-sentiment:", err)
			os.Exit(1)
		}
		preview := text
		if len(preview) > 55 {
			preview = preview[:55]
		}
		fmt.Printf("  [%-8s %.0f%%] %s\n", r.Label, r.Confidence*100, preview)
	}

	section("9. CLASSIFY")
	cats := []string{"shipping", "technical", "feature request", "billing", "other"}
	for _, ticket := range []string{
		"My order hasn't shipped after five days. Please help.",
		"The API returns a 500 error when the payload exceeds 1 MB.",
		"It would be great to have a dark mode option.",
		"I was charged twice for the same subscription this month.",
	} {
		r, err := classify.Classify(ctx, config, classify.ClassifyInput{Text: ticket, Categories: cats})
		if err != nil {
			fmt.Fprintln(os.Stderr, "classify:", err)
			os.Exit(1)
		}
		top := r.Classifications[0]
		cat, _ := top["category"].(string)
		conf, _ := top["confidence"].(float64)
		preview := ticket
		if len(preview) > 55 {
			preview = preview[:55]
		}
		fmt.Printf("  [%-16s %.0f%%]  %s\n", cat, conf*100, preview)
	}

	section("10. RECOGNIZE INTENT")
	riCtx := "You are a customer support routing system for an e-commerce platform."
	for _, msg := range []string{
		"Where is my order? I placed it three days ago.",
		"I want a refund for the broken item.",
		"Can you tell me your business hours?",
		"I'd like to upgrade my subscription to the pro plan.",
	} {
		r, err := recognize_intent.RecognizeIntent(ctx, config, recognize_intent.RecognizeIntentInput{
			Text:    msg,
			Intents: []string{"query_order", "request_refund", "general_inquiry", "manage_subscription"},
			Context: &riCtx,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "recognize-intent:", err)
			os.Exit(1)
		}
		preview := msg
		if len(preview) > 50 {
			preview = preview[:50]
		}
		fmt.Printf("  [%-20s %.0f%%]  \"%s\"\n", r.Intent, r.Confidence*100, preview)
	}

	// ── Advanced: complex extraction and generation ───────────────────

	section("11. EXTRACT ENTITIES")
	entText := "On March 10, 2024, NASA astronaut Sarah Mitchell landed at Kennedy Space Center in Florida after a six-month mission."
	entR, err := extract_entities.ExtractEntities(ctx, config, extract_entities.ExtractEntitiesInput{
		Text:        entText,
		EntityTypes: []string{"person", "organization", "location", "date"},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "extract-entities:", err)
		os.Exit(1)
	}
	fmt.Printf("Text: %s\nEntities:\n", entText)
	for _, e := range entR.Entities {
		etype, _ := e["type"].(string)
		etext, _ := e["text"].(string)
		fmt.Printf("  [%-12s] \"%s\"\n", etype, etext)
	}

	section("12. EXTRACT JSON")
	jobPost := "We are looking for a Senior Backend Engineer in Berlin. " +
		"Requirements: 5+ years of experience, proficiency in Go or Rust, " +
		"experience with Kubernetes. Salary range: €80,000–€110,000."
	ejR, err := extract_json.ExtractJson(ctx, config, extract_json.ExtractJsonInput{
		Text: jobPost,
		Fields: []map[string]any{
			{"name": "title", "description": "Job title", "type": "string"},
			{"name": "location", "description": "City or country", "type": "string"},
			{"name": "skills", "description": "Required technical skills", "type": "array"},
			{"name": "experience_years", "description": "Minimum years of experience", "type": "number"},
			{"name": "salary_range", "description": "Salary range", "type": "string"},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "extract-json:", err)
		os.Exit(1)
	}
	fmt.Printf("Text     : %s\nExtracted: %v\nMissing  : %v\n", jobPost, ejR.Extracted, ejR.Missing)

	section("13. ANSWER QUESTION")
	qaCtx := "AIFunc is a function-based AI toolkit. Developers declare the packages they need " +
		"in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. " +
		"Each package supports a mock mode for testing without consuming API credits."
	type qaPair struct {
		q   string
		ctx *string
	}
	for _, p := range []qaPair{
		{"Which languages does AIFunc support?", &qaCtx},
		{"What is mock mode used for?", &qaCtx},
		{"What is a monad in functional programming?", nil},
	} {
		r, err := answer_question.AnswerQuestion(ctx, config, answer_question.AnswerQuestionInput{
			Question: p.q, Context: p.ctx, MaxLength: intPtr(60),
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "answer-question:", err)
			os.Exit(1)
		}
		source := "general knowledge"
		if r.Grounded {
			source = "from context"
		}
		fmt.Printf("  Q: %s\n  A: %s  [%s, conf: %.0f%%]\n\n", p.q, r.Answer, source, r.Confidence*100)
	}

	section("14. GENERATE REPLY")
	replyMsg := "I placed an order three days ago but haven't received a shipping confirmation yet."
	grR, err := generate_reply.GenerateReply(ctx, config, generate_reply.GenerateReplyInput{
		Message: replyMsg,
		Tone:    strPtr("empathetic"),
		Context: strPtr("You are a customer support agent for an online store."),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate-reply:", err)
		os.Exit(1)
	}
	fmt.Printf("Customer : %s\nReply    : %s\n", replyMsg, grR.Reply)

	section("15. GENERATE POST")
	gpR, err := generate_post.GeneratePost(ctx, config, generate_post.GeneratePostInput{
		Topic:           "How switching to async Python cut our API response time by 60%",
		Platform:        strPtr("linkedin"),
		Tone:            strPtr("professional"),
		IncludeHashtags: boolPtr(true),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate-post:", err)
		os.Exit(1)
	}
	tagged := make([]string, len(gpR.Hashtags))
	for i, t := range gpR.Hashtags {
		tagged[i] = "#" + t
	}
	fmt.Printf("Post     : %s\nHashtags : %v\n", gpR.Post, tagged)

	section("16. GENERATE EMAIL")
	geR, err := generate_email.GenerateEmail(ctx, config, generate_email.GenerateEmailInput{
		Intent:        "Apologize to a customer for a billing error and explain the resolution",
		Tone:          strPtr("formal"),
		SenderName:    strPtr("Billing Support Team"),
		RecipientName: strPtr("Alex"),
		KeyPoints: []string{
			"An incorrect charge of $29.99 was applied on June 1st",
			"The charge has been fully refunded and will appear within 3–5 business days",
			"We have applied a 20% discount to the next invoice as compensation",
		},
		Language: strPtr("English"),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate-email:", err)
		os.Exit(1)
	}
	fmt.Printf("Subject: %s\nBody:\n%s\n", geR.Subject, geR.Body)

	section("17. SCORE QUALITY")
	type qualSample struct{ text, audience, purpose string }
	for _, s := range []qualSample{
		{"Our product is good. It has many features. Users like it.", "customers", "marketing"},
		{"To set up the CI pipeline: 1) Install Docker and the GitHub CLI. " +
			"2) Create .github/workflows/deploy.yml. 3) Push to main to trigger deployment.",
			"developers", "explanation"},
	} {
		r, err := score_quality.ScoreQuality(ctx, config, score_quality.ScoreQualityInput{
			Text:           s.text,
			TargetAudience: &s.audience,
			Purpose:        &s.purpose,
			MaxSuggestions: intPtr(3),
			Strictness:     intPtr(3),
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "score-quality:", err)
			os.Exit(1)
		}
		preview := s.text
		if len(preview) > 55 {
			preview = preview[:55]
		}
		fmt.Printf("Text       : %s...\nScore      : %d/100  [%s]\nSummary    : %s\nSuggestions:\n",
			preview, r.OverallScore, r.Level, r.Summary)
		for _, sg := range r.Suggestions {
			fmt.Printf("  - %s\n", sg)
		}
		fmt.Println()
	}

	if config.Mock {
		fmt.Println("Notice: You are using mock mode for offline testing. " +
			"Configure a real model for the full experience.")
	}
}
