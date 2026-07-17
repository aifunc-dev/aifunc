using Aifunc;
using Aifunc.AnswerStream;
using Aifunc.ArticleStream;
using Aifunc.ChatStream;
using Aifunc.ExplainStream;
using Aifunc.ReviewStream;
using Aifunc.TranslateStream;
using Aifunc.WriteStream;

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

static async Task StreamPrintAsync(IAsyncEnumerable<string> tokens)
{
    await foreach (var token in tokens)
        Console.Write(token);
    Console.WriteLine();
}

var article =
    "In 1915, Albert Einstein published the General Theory of Relativity, " +
    "fundamentally transforming our understanding of physics. The theory posits " +
    "that gravity is not an invisible force, but rather a curvature of spacetime " +
    "caused by the presence of mass and energy. This groundbreaking framework " +
    "revolutionized modern science and introduced the famous equation E=mc².";

var codeSnippet =
    """
    def fetch_user(user_id):
        conn = get_connection()
        result = conn.execute(f"SELECT * FROM users WHERE id = {user_id}")
        return result.fetchone()
    """;

// ─── Conversational & Q&A ─────────────────────────────────────────

Section("1. CHAT STREAM");
Console.WriteLine("User: Explain async/await and IAsyncEnumerable in C# in 3 sentences.");
Console.WriteLine();
Console.Write("Assistant: ");
await StreamPrintAsync(ChatStream.ChatStreamAsync(config, new ChatStreamTypes.ChatStreamInput(
    messages:
    [
        new("user", "Explain async/await and IAsyncEnumerable in C# in 3 sentences."),
    ])));

Section("2. ANSWER STREAM (with context / RAG)");
var context =
    "AIFunc is a function-based AI toolkit. Developers declare the packages they need " +
    "in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. " +
    "Each package supports a mock mode for testing without consuming API credits. " +
    "Streaming packages return tokens incrementally via IAsyncEnumerable.";
var question = "How does AIFunc support offline testing, and what do streaming packages return?";
Console.WriteLine($"Q: {question}");
Console.WriteLine();
Console.Write("A: ");
await StreamPrintAsync(AnswerStream.AnswerStreamAsync(config, new AnswerStreamTypes.AnswerStreamInput(
    question: question,
    audience: "technical",
    context: context,
    depth: "concise")));

Section("3. EXPLAIN STREAM");
Console.WriteLine("Topic: the .NET garbage collector");
Console.WriteLine();
await StreamPrintAsync(ExplainStream.ExplainStreamAsync(config, new ExplainStreamTypes.ExplainStreamInput(
    topic: "the .NET garbage collector",
    audience: "intermediate",
    depth: "standard")));

// ─── Long-form writing ────────────────────────────────────────────

Section("4. ARTICLE STREAM");
var title = "Why Typed AI Functions Beat Ad-Hoc Prompt Scripts";
var outline =
    "- The cost of untyped prompt glue code\n" +
    "- How function-shaped AI APIs improve testability\n" +
    "- Streaming vs batch for product UX\n" +
    "- Practical adoption tips";
Console.WriteLine($"Title  : {title}");
Console.WriteLine($"Outline: {outline}");
Console.WriteLine();
await StreamPrintAsync(ArticleStream.ArticleStreamAsync(config, new ArticleStreamTypes.ArticleStreamInput(
    title: title,
    audience: "developers",
    outline: outline,
    style: "informational",
    wordCount: 250)));

Section("5. WRITE STREAM");
var prompt =
    "Write a short internal proposal recommending that our team adopt AIFunc " +
    "for customer-support reply generation.";
var structure =
    "1. Problem\n" +
    "2. Proposed approach\n" +
    "3. Expected benefits\n" +
    "4. Next steps";
Console.WriteLine($"Prompt   : {prompt}");
Console.WriteLine($"Structure: {structure}");
Console.WriteLine();
await StreamPrintAsync(WriteStream.WriteStreamAsync(config, new WriteStreamTypes.WriteStreamInput(
    prompt: prompt,
    audience: "engineers",
    format: "proposal",
    structure: structure,
    tone: "professional",
    wordCount: 300)));

// ─── Translation & review ─────────────────────────────────────────

Section("6. TRANSLATE STREAM");
Console.WriteLine($"Original (EN):\n{article}");
Console.WriteLine();
Console.WriteLine("Translation (zh-CN):");
Console.WriteLine();
await StreamPrintAsync(TranslateStream.TranslateStreamAsync(config, new TranslateStreamTypes.TranslateStreamInput(
    text: article,
    targetLang: "zh-CN",
    domain: "technical",
    style: "natural")));

Section("7. REVIEW STREAM");
Console.WriteLine($"Code under review:\n{codeSnippet}");
Console.WriteLine("Findings:");
Console.WriteLine();
await StreamPrintAsync(ReviewStream.ReviewStreamAsync(config, new ReviewStreamTypes.ReviewStreamInput(
    content: codeSnippet,
    context: "Simple data-access helper in a web API.",
    focus: "correctness, security",
    language: "Python",
    outputLanguage: "English",
    severity: "all",
    type: "code")));

if (config.Mock)
{
    Console.WriteLine("Notice: You are using mock mode for offline testing. " +
        "Configure a real model for the full experience.");
}
