using Aifunc;
using Aifunc.AnalyzeSentiment;
using Aifunc.AnswerQuestion;
using Aifunc.Classify;
using Aifunc.DetectLanguage;
using Aifunc.ExtractEntities;
using Aifunc.ExtractJson;
using Aifunc.ExtractKeywords;
using Aifunc.GenerateEmail;
using Aifunc.GeneratePost;
using Aifunc.GenerateReply;
using Aifunc.GenerateSlug;
using Aifunc.GenerateTitle;
using Aifunc.RecognizeIntent;
using Aifunc.Rewrite;
using Aifunc.ScoreQuality;
using Aifunc.Summarize;
using Aifunc.Translate;

// var config = new AIFuncConfig
// {
//     BaseUrl = "https://your-api-endpoint/v1",
//     Model = "your-model-name",
//     ApiKey = "your-api-key",
//     MaxRetries = 3,
// };

// To use a real model, replace the line below with the commented config above.
var config = new AIFuncConfig { Mock = true };

if (config.Mock)
{
    Console.WriteLine("Notice: You are using mock mode for offline testing. " +
        "Configure a real model for the full experience. Continuing with mock responses...");
}

static void Section(string title)
{
    Console.WriteLine();
    Console.WriteLine(new string('=', 60));
    Console.WriteLine($"  {title}");
    Console.WriteLine(new string('=', 60));
}

var article =
    "In 1915, Albert Einstein published the General Theory of Relativity, " +
    "fundamentally transforming our understanding of physics. The theory posits " +
    "that gravity is not an invisible force, but rather a curvature of spacetime " +
    "caused by the presence of mass and energy. This groundbreaking framework " +
    "revolutionized modern science and introduced the famous equation E=mc².";

// ── Easy: single input → single output ────────────────────────────

Section("1. DETECT LANGUAGE");
foreach (var text in new[]
{
    "The quick brown fox jumps over the lazy dog.",
    "Der schnelle braune Fuchs springt über den faulen Hund.",
    "Le renard brun rapide saute par-dessus le chien paresseux.",
    "El veloz zorro marrón salta sobre el perro perezoso.",
})
{
    var r = await DetectLanguage.DetectLanguageAsync(config, new DetectLanguageTypes.DetectLanguageInput(text));
    var preview = text.Length > 40 ? text[..40] : text;
    Console.WriteLine($"  [{r.Language}] {r.LanguageName} (conf: {r.Confidence * 100:F0}%)  \"{preview}\"");
}

Section("2. GENERATE SLUG");
var slugResult = await GenerateSlug.GenerateSlugAsync(config, new GenerateSlugTypes.GenerateSlugInput(
    title: "10 Practical Tips for Writing Faster Python Code",
    language: "en"));
Console.WriteLine("Title : 10 Practical Tips for Writing Faster Python Code");
Console.WriteLine($"Slug  : {slugResult.Slug}");
Console.WriteLine($"Meta  : {slugResult.MetaDescription}");
Console.WriteLine($"Tags  : [{string.Join(", ", slugResult.Tags)}]");

Section("3. SUMMARIZE");
var sumResult = await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput(article, maxLength: 30));
Console.WriteLine($"Summary   : {sumResult.Summary}");
Console.WriteLine($"Word count: {sumResult.WordCount}");

Section("4. TRANSLATE");
var trResult = await Translate.TranslateAsync(config, new TranslateTypes.TranslateInput(
    targetLang: "es",
    text: "The meeting has been moved to Friday at 3 PM."));
Console.WriteLine("Original : The meeting has been moved to Friday at 3 PM.");
Console.WriteLine($"Spanish  : {trResult.Translation}");
Console.WriteLine($"Detected : {trResult.SourceLang}");

Section("5. REWRITE");
var orig = "hey, just wanna let u know the deploy went fine, no issues at all";
var rwResult = await Rewrite.RewriteAsync(config, new RewriteTypes.RewriteInput(
    style: "formal",
    text: orig));
Console.WriteLine($"Casual : {orig}");
Console.WriteLine($"Formal : {rwResult.Rewritten}");

// ── Medium: structured output or multiple parameters ──────────────

Section("6. GENERATE TITLE");
var titleContent =
    "This guide covers how to use Docker and GitHub Actions to automate " +
    "testing and deployment of a Node.js application to a cloud server.";
var titResult = await GenerateTitle.GenerateTitleAsync(config, new GenerateTitleTypes.GenerateTitleInput(
    content: titleContent,
    style: "seo",
    count: 4));
Console.WriteLine($"Content: {titleContent}");
Console.WriteLine("Titles:");
for (var i = 0; i < titResult.Titles.Count; i++)
    Console.WriteLine($"  {i + 1}. {titResult.Titles[i]}");

Section("7. EXTRACT KEYWORDS");
var kwResult = await ExtractKeywords.ExtractKeywordsAsync(config, new ExtractKeywordsTypes.ExtractKeywordsInput(
    text: article,
    maxKeywords: 5));
Console.WriteLine("Keywords from article:");
foreach (var kw in kwResult.Keywords)
{
    var word = kw.TryGetValue("word", out var w) ? w?.ToString() ?? "" : "";
    var relevance = kw.TryGetValue("relevance", out var rv) && rv is IConvertible c ? Convert.ToDouble(c) : 0.0;
    Console.WriteLine($"  {word,-30} relevance: {relevance:F2}");
}

Section("8. ANALYZE SENTIMENT");
foreach (var text in new[]
{
    "The product arrived on time and works perfectly. Very happy!",
    "Terrible experience. The package was damaged and support ignored my emails.",
    "Item received. Does what it says.",
})
{
    var r = await AnalyzeSentiment.AnalyzeSentimentAsync(config, new AnalyzeSentimentTypes.AnalyzeSentimentInput(
        text: text,
        labels: ["positive", "negative", "neutral"]));
    var preview = text.Length > 55 ? text[..55] : text;
    Console.WriteLine($"  [{r.Label,-8} {r.Confidence * 100:F0}%] {preview}");
}

Section("9. CLASSIFY");
var cats = new List<string> { "shipping", "technical", "feature request", "billing", "other" };
foreach (var ticket in new[]
{
    "My order hasn't shipped after five days. Please help.",
    "The API returns a 500 error when the payload exceeds 1 MB.",
    "It would be great to have a dark mode option.",
    "I was charged twice for the same subscription this month.",
})
{
    var r = await Classify.ClassifyAsync(config, new ClassifyTypes.ClassifyInput(
        categories: cats,
        text: ticket));
    if (r.Classifications.Count == 0) continue;
    var top = r.Classifications[0];
    var cat = top.TryGetValue("category", out var cv) ? cv?.ToString() ?? "" : "";
    var conf = top.TryGetValue("confidence", out var cfv) && cfv is IConvertible c ? Convert.ToDouble(c) : 0.0;
    var preview = ticket.Length > 55 ? ticket[..55] : ticket;
    Console.WriteLine($"  [{cat,-16} {conf * 100:F0}%]  {preview}");
}

Section("10. RECOGNIZE INTENT");
var riCtx = "You are a customer support routing system for an e-commerce platform.";
foreach (var msg in new[]
{
    "Where is my order? I placed it three days ago.",
    "I want a refund for the broken item.",
    "Can you tell me your business hours?",
    "I'd like to upgrade my subscription to the pro plan.",
})
{
    var r = await RecognizeIntent.RecognizeIntentAsync(config, new RecognizeIntentTypes.RecognizeIntentInput(
        intents: ["query_order", "request_refund", "general_inquiry", "manage_subscription"],
        text: msg,
        context: riCtx));
    var preview = msg.Length > 50 ? msg[..50] : msg;
    Console.WriteLine($"  [{r.Intent,-20} {r.Confidence * 100:F0}%]  \"{preview}\"");
}

// ── Advanced: complex extraction and generation ───────────────────

Section("11. EXTRACT ENTITIES");
var entText =
    "On March 10, 2024, NASA astronaut Sarah Mitchell landed at Kennedy Space Center " +
    "in Florida after a six-month mission.";
var entResult = await ExtractEntities.ExtractEntitiesAsync(config, new ExtractEntitiesTypes.ExtractEntitiesInput(
    text: entText,
    entityTypes: ["person", "organization", "location", "date"]));
Console.WriteLine($"Text: {entText}");
Console.WriteLine("Entities:");
foreach (var e in entResult.Entities)
{
    var etype = e.TryGetValue("type", out var tv) ? tv?.ToString() ?? "" : "";
    var etext = e.TryGetValue("text", out var etv) ? etv?.ToString() ?? "" : "";
    Console.WriteLine($"  [{etype,-12}] \"{etext}\"");
}

Section("12. EXTRACT JSON");
var jobPost =
    "We are looking for a Senior Backend Engineer in Berlin. " +
    "Requirements: 5+ years of experience, proficiency in Go or Rust, " +
    "experience with Kubernetes. Salary range: €80,000–€110,000.";
var ejResult = await ExtractJson.ExtractJsonAsync(config, new ExtractJsonTypes.ExtractJsonInput(
    fields:
    [
        new() { ["name"] = "title",            ["description"] = "Job title",                          ["type"] = "string"  },
        new() { ["name"] = "location",         ["description"] = "City or country",                    ["type"] = "string"  },
        new() { ["name"] = "skills",           ["description"] = "Required technical skills",          ["type"] = "array"   },
        new() { ["name"] = "experience_years", ["description"] = "Minimum years of experience",        ["type"] = "number"  },
        new() { ["name"] = "salary_range",     ["description"] = "Salary range",                       ["type"] = "string"  },
    ],
    text: jobPost));
Console.WriteLine($"Text     : {jobPost}");
Console.WriteLine($"Extracted: {string.Join(", ", ejResult.Extracted.Select(kv => $"{kv.Key}={kv.Value}"))}");
Console.WriteLine($"Missing  : [{string.Join(", ", ejResult.Missing)}]");

Section("13. ANSWER QUESTION");
var qaCtx =
    "AIFunc is a function-based AI toolkit. Developers declare the packages they need " +
    "in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. " +
    "Each package supports a mock mode for testing without consuming API credits.";
foreach (var (q, ctx) in new (string q, string? ctx)[]
{
    ("Which languages does AIFunc support?", qaCtx),
    ("What is mock mode used for?",          qaCtx),
    ("What is a monad in functional programming?", null),
})
{
    var r = await AnswerQuestion.AnswerQuestionAsync(config, new AnswerQuestionTypes.AnswerQuestionInput(
        question: q,
        context: ctx,
        maxLength: 60));
    var source = r.Grounded ? "from context" : "general knowledge";
    Console.WriteLine($"  Q: {q}");
    Console.WriteLine($"  A: {r.Answer}  [{source}, conf: {r.Confidence * 100:F0}%]");
    Console.WriteLine();
}

Section("14. GENERATE REPLY");
var replyMsg = "I placed an order three days ago but haven't received a shipping confirmation yet.";
var grResult = await GenerateReply.GenerateReplyAsync(config, new GenerateReplyTypes.GenerateReplyInput(
    message: replyMsg,
    tone: "empathetic",
    context: "You are a customer support agent for an online store."));
Console.WriteLine($"Customer : {replyMsg}");
Console.WriteLine($"Reply    : {grResult.Reply}");

Section("15. GENERATE POST");
var gpResult = await GeneratePost.GeneratePostAsync(config, new GeneratePostTypes.GeneratePostInput(
    topic: "How switching to async Python cut our API response time by 60%",
    platform: "linkedin",
    tone: "professional",
    includeHashtags: true));
var tagged = gpResult.Hashtags.Select(t => "#" + t).ToList();
Console.WriteLine($"Post     : {gpResult.Post}");
Console.WriteLine($"Hashtags : [{string.Join(", ", tagged)}]");

Section("16. GENERATE EMAIL");
var geResult = await GenerateEmail.GenerateEmailAsync(config, new GenerateEmailTypes.GenerateEmailInput(
    intent: "Apologize to a customer for a billing error and explain the resolution",
    tone: "formal",
    senderName: "Billing Support Team",
    recipientName: "Alex",
    keyPoints:
    [
        "An incorrect charge of $29.99 was applied on June 1st",
        "The charge has been fully refunded and will appear within 3–5 business days",
        "We have applied a 20% discount to the next invoice as compensation",
    ],
    language: "English"));
Console.WriteLine($"Subject: {geResult.Subject}");
Console.WriteLine($"Body:\n{geResult.Body}");

Section("17. SCORE QUALITY");
foreach (var (text, audience, purpose) in new (string text, string audience, string purpose)[]
{
    ("Our product is good. It has many features. Users like it.",
     "customers", "marketing"),
    ("To set up the CI pipeline: 1) Install Docker and the GitHub CLI. " +
     "2) Create .github/workflows/deploy.yml. 3) Push to main to trigger deployment.",
     "developers", "explanation"),
})
{
    var r = await ScoreQuality.ScoreQualityAsync(config, new ScoreQualityTypes.ScoreQualityInput(
        text: text,
        targetAudience: audience,
        purpose: purpose,
        maxSuggestions: 3,
        strictness: 3));
    var preview = text.Length > 55 ? text[..55] : text;
    Console.WriteLine($"Text       : {preview}...");
    Console.WriteLine($"Score      : {r.OverallScore}/100  [{r.Level}]");
    Console.WriteLine($"Summary    : {r.Summary}");
    Console.WriteLine("Suggestions:");
    foreach (var sg in r.Suggestions)
        Console.WriteLine($"  - {sg}");
    Console.WriteLine();
}

if (config.Mock)
{
    Console.WriteLine("Notice: You are using mock mode for offline testing. " +
        "Configure a real model for the full experience.");
}
