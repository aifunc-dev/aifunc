import aifunc.AIFuncConfig;
import aifunc.analyze_sentiment.AnalyzeSentiment;
import aifunc.analyze_sentiment.AnalyzeSentimentTypes.AnalyzeSentimentInput;
import aifunc.answer_question.AnswerQuestion;
import aifunc.answer_question.AnswerQuestionTypes.AnswerQuestionInput;
import aifunc.classify.Classify;
import aifunc.classify.ClassifyTypes.ClassifyInput;
import aifunc.detect_language.DetectLanguage;
import aifunc.detect_language.DetectLanguageTypes.DetectLanguageInput;
import aifunc.extract_entities.ExtractEntities;
import aifunc.extract_entities.ExtractEntitiesTypes.ExtractEntitiesInput;
import aifunc.extract_json.ExtractJson;
import aifunc.extract_json.ExtractJsonTypes.ExtractJsonInput;
import aifunc.extract_keywords.ExtractKeywords;
import aifunc.extract_keywords.ExtractKeywordsTypes.ExtractKeywordsInput;
import aifunc.generate_email.GenerateEmail;
import aifunc.generate_email.GenerateEmailTypes.GenerateEmailInput;
import aifunc.generate_post.GeneratePost;
import aifunc.generate_post.GeneratePostTypes.GeneratePostInput;
import aifunc.generate_reply.GenerateReply;
import aifunc.generate_reply.GenerateReplyTypes.GenerateReplyInput;
import aifunc.generate_slug.GenerateSlug;
import aifunc.generate_slug.GenerateSlugTypes.GenerateSlugInput;
import aifunc.generate_title.GenerateTitle;
import aifunc.generate_title.GenerateTitleTypes.GenerateTitleInput;
import aifunc.recognize_intent.RecognizeIntent;
import aifunc.recognize_intent.RecognizeIntentTypes.RecognizeIntentInput;
import aifunc.rewrite.Rewrite;
import aifunc.rewrite.RewriteTypes.RewriteInput;
import aifunc.score_quality.ScoreQuality;
import aifunc.score_quality.ScoreQualityTypes.ScoreQualityInput;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;
import aifunc.translate.Translate;
import aifunc.translate.TranslateTypes.TranslateInput;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class Main {
    // AIFuncConfig config = AIFuncConfig.builder()
    //         .baseUrl("https://your-api-endpoint/v1")
    //         .model("your-model-name")
    //         .apiKey("your-api-key")
    //         .maxRetries(3)
    //         .build();

    // To use a real model, replace the line below with the commented config above.
    static final AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

    public static void main(String[] args) {
        if (config.isMock()) {
            System.out.println("Notice: You are using mock mode for offline testing. "
                    + "Configure a real model for the full experience. Continuing with mock responses...");
        }

        String article = "In 1915, Albert Einstein published the General Theory of Relativity, "
                + "fundamentally transforming our understanding of physics. The theory posits "
                + "that gravity is not an invisible force, but rather a curvature of spacetime "
                + "caused by the presence of mass and energy. This groundbreaking framework "
                + "revolutionized modern science and introduced the famous equation E=mc².";

        // ── Easy: single input → single output ────────────────────────────

        section("1. DETECT LANGUAGE");
        for (String text : List.of(
                "The quick brown fox jumps over the lazy dog.",
                "Der schnelle braune Fuchs springt über den faulen Hund.",
                "Le renard brun rapide saute par-dessus le chien paresseux.",
                "El veloz zorro marrón salta sobre el perro perezoso."
        )) {
            var r = DetectLanguage.detectLanguage(config, new DetectLanguageInput(text)).join();
            System.out.printf("  [%s] %s (conf: %.0f%%)  \"%s\"%n",
                    r.getLanguage(), r.getLanguageName(), r.getConfidence() * 100, preview(text, 40));
        }

        section("2. GENERATE SLUG");
        var sr = GenerateSlug.generateSlug(config, new GenerateSlugInput(
                "10 Practical Tips for Writing Faster Python Code", null, "en")).join();
        System.out.println("Title : 10 Practical Tips for Writing Faster Python Code");
        System.out.printf("Slug  : %s%nMeta  : %s%nTags  : %s%n",
                sr.getSlug(), sr.getMetaDescription(), sr.getTags());

        section("3. SUMMARIZE");
        var sumR = Summarize.summarize(config, new SummarizeInput(article, 30)).join();
        System.out.printf("Summary   : %s%nWord count: %d%n", sumR.getSummary(), sumR.getWordCount());

        section("4. TRANSLATE");
        var tr = Translate.translate(config, new TranslateInput(
                "es", "The meeting has been moved to Friday at 3 PM.", null)).join();
        System.out.println("Original : The meeting has been moved to Friday at 3 PM.");
        System.out.printf("Spanish  : %s%nDetected : %s%n", tr.getTranslation(), tr.getSourceLang());

        section("5. REWRITE");
        String orig = "hey, just wanna let u know the deploy went fine, no issues at all";
        var rwr = Rewrite.rewrite(config, new RewriteInput("formal", orig, null)).join();
        System.out.printf("Casual : %s%nFormal : %s%n", orig, rwr.getRewritten());

        // ── Medium: structured output or multiple parameters ──────────────

        section("6. GENERATE TITLE");
        String titleContent = "This guide covers how to use Docker and GitHub Actions to automate "
                + "testing and deployment of a Node.js application to a cloud server.";
        var titR = GenerateTitle.generateTitle(config, new GenerateTitleInput(
                titleContent, 4, null, "seo")).join();
        System.out.printf("Content: %s%nTitles:%n", titleContent);
        List<String> titles = titR.getTitles();
        for (int i = 0; i < titles.size(); i++) {
            System.out.printf("  %d. %s%n", i + 1, titles.get(i));
        }

        section("7. EXTRACT KEYWORDS");
        var kwR = ExtractKeywords.extractKeywords(config, new ExtractKeywordsInput(article, 5)).join();
        System.out.println("Keywords from article:");
        for (Map<String, Object> kw : kwR.getKeywords()) {
            Object word = kw.get("word");
            Object relevance = kw.get("relevance");
            System.out.printf("  %-30s relevance: %s%n", word, relevance);
        }

        section("8. ANALYZE SENTIMENT");
        for (String text : List.of(
                "The product arrived on time and works perfectly. Very happy!",
                "Terrible experience. The package was damaged and support ignored my emails.",
                "Item received. Does what it says."
        )) {
            var r = AnalyzeSentiment.analyzeSentiment(config, new AnalyzeSentimentInput(
                    text, List.of("positive", "negative", "neutral"), null)).join();
            System.out.printf("  [%-8s %.0f%%] %s%n",
                    r.getLabel(), r.getConfidence() * 100, preview(text, 55));
        }

        section("9. CLASSIFY");
        List<String> cats = List.of("shipping", "technical", "feature request", "billing", "other");
        for (String ticket : List.of(
                "My order hasn't shipped after five days. Please help.",
                "The API returns a 500 error when the payload exceeds 1 MB.",
                "It would be great to have a dark mode option.",
                "I was charged twice for the same subscription this month."
        )) {
            var r = Classify.classify(config, new ClassifyInput(cats, ticket, null)).join();
            Map<String, Object> top = r.getClassifications().get(0);
            String cat = String.valueOf(top.get("category"));
            double conf = ((Number) top.get("confidence")).doubleValue();
            System.out.printf("  [%-16s %.0f%%]  %s%n", cat, conf * 100, preview(ticket, 55));
        }

        section("10. RECOGNIZE INTENT");
        String riCtx = "You are a customer support routing system for an e-commerce platform.";
        List<String> intents = List.of(
                "query_order", "request_refund", "general_inquiry", "manage_subscription");
        for (String msg : List.of(
                "Where is my order? I placed it three days ago.",
                "I want a refund for the broken item.",
                "Can you tell me your business hours?",
                "I'd like to upgrade my subscription to the pro plan."
        )) {
            var r = RecognizeIntent.recognizeIntent(config, new RecognizeIntentInput(
                    intents, msg, riCtx)).join();
            System.out.printf("  [%-20s %.0f%%]  \"%s\"%n",
                    r.getIntent(), r.getConfidence() * 100, preview(msg, 50));
        }

        // ── Advanced: complex extraction and generation ───────────────────

        section("11. EXTRACT ENTITIES");
        String entText = "On March 10, 2024, NASA astronaut Sarah Mitchell landed at "
                + "Kennedy Space Center in Florida after a six-month mission.";
        var entR = ExtractEntities.extractEntities(config, new ExtractEntitiesInput(
                entText, List.of("person", "organization", "location", "date"))).join();
        System.out.printf("Text: %s%nEntities:%n", entText);
        for (Map<String, Object> e : entR.getEntities()) {
            System.out.printf("  [%-12s] \"%s\"%n", e.get("type"), e.get("text"));
        }

        section("12. EXTRACT JSON");
        String jobPost = "We are looking for a Senior Backend Engineer in Berlin. "
                + "Requirements: 5+ years of experience, proficiency in Go or Rust, "
                + "experience with Kubernetes. Salary range: €80,000–€110,000.";
        List<Map<String, Object>> fields = new ArrayList<>();
        fields.add(field("title", "Job title", "string"));
        fields.add(field("location", "City or country", "string"));
        fields.add(field("skills", "Required technical skills", "array"));
        fields.add(field("experience_years", "Minimum years of experience", "number"));
        fields.add(field("salary_range", "Salary range", "string"));
        var ejR = ExtractJson.extractJson(config, new ExtractJsonInput(fields, jobPost)).join();
        System.out.printf("Text     : %s%nExtracted: %s%nMissing  : %s%n",
                jobPost, ejR.getExtracted(), ejR.getMissing());

        section("13. ANSWER QUESTION");
        String qaCtx = "AIFunc is a function-based AI toolkit. Developers declare the packages they need "
                + "in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. "
                + "Each package supports a mock mode for testing without consuming API credits.";
        record QaPair(String q, String ctx) {}
        for (QaPair p : List.of(
                new QaPair("Which languages does AIFunc support?", qaCtx),
                new QaPair("What is mock mode used for?", qaCtx),
                new QaPair("What is a monad in functional programming?", null)
        )) {
            var r = AnswerQuestion.answerQuestion(config, new AnswerQuestionInput(
                    p.q(), p.ctx(), null, 60)).join();
            String source = Boolean.TRUE.equals(r.getGrounded()) ? "from context" : "general knowledge";
            System.out.printf("  Q: %s%n  A: %s  [%s, conf: %.0f%%]%n%n",
                    p.q(), r.getAnswer(), source, r.getConfidence() * 100);
        }

        section("14. GENERATE REPLY");
        String replyMsg = "I placed an order three days ago but haven't received a shipping confirmation yet.";
        var grR = GenerateReply.generateReply(config, new GenerateReplyInput(
                replyMsg, "You are a customer support agent for an online store.", null, "empathetic")).join();
        System.out.printf("Customer : %s%nReply    : %s%n", replyMsg, grR.getReply());

        section("15. GENERATE POST");
        var gpR = GeneratePost.generatePost(config, new GeneratePostInput(
                "How switching to async Python cut our API response time by 60%",
                true, null, "linkedin", "professional")).join();
        List<String> tagged = new ArrayList<>();
        if (gpR.getHashtags() != null) {
            for (String t : gpR.getHashtags()) {
                tagged.add("#" + t);
            }
        }
        System.out.printf("Post     : %s%nHashtags : %s%n", gpR.getPost(), tagged);

        section("16. GENERATE EMAIL");
        var geR = GenerateEmail.generateEmail(config, new GenerateEmailInput(
                "Apologize to a customer for a billing error and explain the resolution",
                List.of(
                        "An incorrect charge of $29.99 was applied on June 1st",
                        "The charge has been fully refunded and will appear within 3–5 business days",
                        "We have applied a 20% discount to the next invoice as compensation"
                ),
                "English",
                "Alex",
                "Billing Support Team",
                "formal"
        )).join();
        System.out.printf("Subject: %s%nBody:%n%s%n", geR.getSubject(), geR.getBody());

        section("17. SCORE QUALITY");
        record QualSample(String text, String audience, String purpose) {}
        for (QualSample s : List.of(
                new QualSample(
                        "Our product is good. It has many features. Users like it.",
                        "customers", "marketing"),
                new QualSample(
                        "To set up the CI pipeline: 1) Install Docker and the GitHub CLI. "
                                + "2) Create .github/workflows/deploy.yml. 3) Push to main to trigger deployment.",
                        "developers", "explanation")
        )) {
            var r = ScoreQuality.scoreQuality(config, new ScoreQualityInput(
                    s.text(), 3, s.purpose(), 3, s.audience())).join();
            System.out.printf("Text       : %s...%nScore      : %d/100  [%s]%nSummary    : %s%nSuggestions:%n",
                    preview(s.text(), 55), r.getOverallScore(), r.getLevel(), r.getSummary());
            if (r.getSuggestions() != null) {
                for (String sg : r.getSuggestions()) {
                    System.out.println("  - " + sg);
                }
            }
            System.out.println();
        }

        if (config.isMock()) {
            System.out.println("Notice: You are using mock mode for offline testing. "
                    + "Configure a real model for the full experience.");
        }
    }

    private static void section(String title) {
        System.out.println();
        System.out.println("=".repeat(60));
        System.out.println("  " + title);
        System.out.println("=".repeat(60));
    }

    private static String preview(String text, int max) {
        if (text == null) return "";
        return text.length() <= max ? text : text.substring(0, max);
    }

    private static Map<String, Object> field(String name, String description, String type) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("name", name);
        m.put("description", description);
        m.put("type", type);
        return m;
    }
}
