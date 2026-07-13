using Aifunc;
using Aifunc.ExtractKeywords;
using Aifunc.GenerateReply;
using Aifunc.RecognizeIntent;
using Aifunc.Summarize;

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
    Console.WriteLine(
        "This example requires a real LLM to produce meaningful results.\n" +
        "Mock mode cannot simulate intent-aware replies grounded in conversation history.\n" +
        "\n" +
        "To run this example, replace the line:\n" +
        "\n" +
        "  var config = new AIFuncConfig { Mock = true };\n" +
        "\n" +
        "with:\n" +
        "\n" +
        "  var config = new AIFuncConfig\n" +
        "  {\n" +
        "      BaseUrl = \"https://your-api-endpoint/v1\",\n" +
        "      Model   = \"your-model-name\",\n" +
        "      ApiKey  = \"your-api-key\",\n" +
        "  };\n"
    );
    return;
}

// Conversation window: keep the last N turns in context before compressing
const int Window        = 4;
const int CompressAfter = 6;

var intentsAvailable = new List<string>
{
    "ask_recommendation", "ask_logistics", "ask_budget",
    "share_preference", "confirm", "other",
};

// Simulated user turns about planning a trip to Europe
var userMessages = new[]
{
    "I'm planning a three-week trip across Europe in September. Where should I start?",
    "I enjoy hiking and local markets. I'd rather skip the big tourist traps.",
    "What's the best way to get from Paris to Barcelona — high-speed train or budget flight?",
    "How much should I budget per day for food and transport in Western Europe?",
    "I've heard the Dolomites are stunning in early autumn. Is it worth a detour from Venice?",
    "Alright, I think I'll do Paris → Barcelona → Rome → Venice → Dolomites. Does that route make sense?",
};

var history       = new List<(string Role, string Text)>();
var topics        = new List<string>();
var intentsLog    = new List<string>();
var memorySummary = "";

async Task MaybeCompress()
{
    if (history.Count <= CompressAfter) return;

    var older     = history[..^Window];
    var olderText = string.Join(" ", older.Select(m => m.Text));

    var sumResult = await Summarize.SummarizeAsync(config,
        new SummarizeTypes.SummarizeInput(olderText, maxLength: 40));
    memorySummary = sumResult.Summary;

    history = history[^Window..].ToList();
    Console.WriteLine($"  [memory compressed → \"{memorySummary}\"]\n");
}

string BuildContext()
{
    var parts = new List<string>();

    if (!string.IsNullOrEmpty(memorySummary))
        parts.Add($"Earlier in this conversation: {memorySummary}");

    var recent = history.TakeLast(Window).ToList();
    if (recent.Count > 0)
    {
        var dialogue = string.Join("\n", recent.Select(m =>
            $"  {(m.Role == "user" ? "User" : "Assistant")}: {m.Text}"));
        parts.Add($"Recent conversation:\n{dialogue}");
    }

    if (topics.Count > 0)
        parts.Add($"Topics discussed so far: {string.Join(", ", topics)}");

    if (intentsLog.Count > 0)
        parts.Add($"User intent pattern: {string.Join(" → ", intentsLog.TakeLast(4))}");

    return string.Join("\n\n", parts);
}

for (var i = 0; i < userMessages.Length; i++)
{
    var userMsg = userMessages[i];
    Console.WriteLine($"[Turn {i + 1}] User: {userMsg}");

    // 1. Classify the user's intent
    var intentResult = await RecognizeIntent.RecognizeIntentAsync(config,
        new RecognizeIntentTypes.RecognizeIntentInput(
            intents: intentsAvailable,
            text:    userMsg));
    intentsLog.Add(intentResult.Intent);

    // 2. Extract keywords and accumulate into the topics list
    var kwResult = await ExtractKeywords.ExtractKeywordsAsync(config,
        new ExtractKeywordsTypes.ExtractKeywordsInput(text: userMsg, maxKeywords: 3));
    foreach (var kw in kwResult.Keywords)
    {
        if (kw.TryGetValue("word", out var wordVal) && wordVal is string word && !topics.Contains(word))
            topics.Add(word);
    }

    // 3. Append user turn to history
    history.Add(("user", userMsg));

    // 4. Compress old history before replying if it has grown too long
    await MaybeCompress();

    // 5. Build context from memory state and generate a reply
    var ctx = BuildContext();
    var replyResult = await GenerateReply.GenerateReplyAsync(config,
        new GenerateReplyTypes.GenerateReplyInput(
            message: userMsg,
            tone:    "friendly",
            context: ctx));

    // 6. Append assistant reply to history
    history.Add(("assistant", replyResult.Reply));

    Console.WriteLine($"         Intent  : {intentResult.Intent} ({intentResult.Confidence:P0})");
    Console.WriteLine($"         Topics  : [{string.Join(", ", topics)}]");
    Console.WriteLine($"         Reply   : {replyResult.Reply}");
    Console.WriteLine();
}

// Final memory state
Console.WriteLine(new string('=', 60));
Console.WriteLine("Final memory state");
Console.WriteLine($"  topics  : [{string.Join(", ", topics)}]");
Console.WriteLine($"  intents : [{string.Join(", ", intentsLog)}]");
if (!string.IsNullOrEmpty(memorySummary))
    Console.WriteLine($"  summary : {memorySummary}");
Console.WriteLine($"  history : {history.Count} turns in window");
