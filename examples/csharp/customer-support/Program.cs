using Aifunc;
using Aifunc.AnalyzeSentiment;
using Aifunc.ExtractJson;
using Aifunc.RecognizeIntent;

// var config = new AIFuncConfig
// {
//     BaseUrl = "https://your-api-endpoint/v1",
//     Model = "your-model-name",
//     ApiKey = "your-api-key",
//     MaxRetries = 3,
// };

// To run this example, replace the mock config below with real credentials:
var config = new AIFuncConfig { Mock = true };

if (config.Mock)
{
    Console.WriteLine(
        "This example requires a real LLM to produce meaningful results.\n" +
        "Mock mode cannot simulate multi-step reasoning (sentiment → intent → extraction).\n" +
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

var messages = new[]
{
    "What the hell?! I ordered this a WEEK ago and it still hasn't shipped! I want my money back NOW!",
    "Hi, I'd like to check on my order #ORD-20240601-123. It's been three days with no shipping update.",
    "Your stupid app crashed again and I lost all my data! Fix it or I'm leaving!",
    "I was charged twice this month. Transaction IDs: TXN-88201 and TXN-88202. Please help.",
    "It would be cool if you added a dark mode. The bright screen hurts my eyes at night.",
    "How do I export my purchase history to CSV? I can't find the option.",
    "I am SO FURIOUS! Your delivery guy threw my package over the fence and it's destroyed! I want a manager NOW!",
    "Any ongoing promotions for loyal customers? I've been a member for 2 years.",
};

var rng     = new Random();
var message = messages[rng.Next(messages.Length)];
Console.WriteLine($"Customer: {message}\n");

// Step 1: Sentiment analysis
var sentiment = await AnalyzeSentiment.AnalyzeSentimentAsync(config,
    new AnalyzeSentimentTypes.AnalyzeSentimentInput(
        text:   message,
        labels: ["angry", "frustrated", "neutral", "happy", "other"]));

Console.WriteLine($"Sentiment: {sentiment.Label} ({sentiment.Confidence:P0})");

if (sentiment.Label == "angry" && sentiment.Confidence > 0.7)
{
    Console.WriteLine("\n=> call_human_agent(message, priority=\"HIGH\")");
    return;
}

// Step 2: Intent recognition
var intents = new List<string>
{
    "query_order", "request_refund", "technical_support",
    "billing_issue", "feature_request", "general_inquiry",
};
var intentResult = await RecognizeIntent.RecognizeIntentAsync(config,
    new RecognizeIntentTypes.RecognizeIntentInput(
        intents: intents,
        text:    message));

Console.WriteLine($"Intent: {intentResult.Intent} ({intentResult.Confidence:P0})");

// Step 3: Route by intent and extract structured data
var fields = intentResult.Intent switch
{
    "query_order" =>
        new List<Dictionary<string, object?>>
        {
            new() { ["name"] = "order_id", ["description"] = "Order number",                    ["type"] = "string" },
            new() { ["name"] = "issue",    ["description"] = "What the customer wants to know", ["type"] = "string" },
        },
    "request_refund" =>
        new List<Dictionary<string, object?>>
        {
            new() { ["name"] = "order_id", ["description"] = "Order number",       ["type"] = "string" },
            new() { ["name"] = "reason",   ["description"] = "Reason for refund",  ["type"] = "string" },
        },
    "technical_support" =>
        new List<Dictionary<string, object?>>
        {
            new() { ["name"] = "issue",    ["description"] = "Technical problem",  ["type"] = "string" },
            new() { ["name"] = "platform", ["description"] = "Device or platform", ["type"] = "string" },
        },
    "billing_issue" =>
        new List<Dictionary<string, object?>>
        {
            new() { ["name"] = "transaction_id", ["description"] = "Transaction ID",    ["type"] = "string" },
            new() { ["name"] = "problem",        ["description"] = "Billing problem",   ["type"] = "string" },
        },
    "feature_request" =>
        new List<Dictionary<string, object?>>
        {
            new() { ["name"] = "feature", ["description"] = "Requested feature", ["type"] = "string" },
        },
    _ =>
        new List<Dictionary<string, object?>>
        {
            new() { ["name"] = "question", ["description"] = "Customer's question", ["type"] = "string" },
        },
};

var info = await ExtractJson.ExtractJsonAsync(config,
    new ExtractJsonTypes.ExtractJsonInput(fields: fields, text: message));

var e = info.Extracted;
Console.WriteLine(intentResult.Intent switch
{
    "query_order"       => $"\n=> query_order_system(order_id=\"{e.GetValueOrDefault("order_id")}\", issue=\"{e.GetValueOrDefault("issue")}\")",
    "request_refund"    => $"\n=> submit_refund(order_id=\"{e.GetValueOrDefault("order_id")}\", reason=\"{e.GetValueOrDefault("reason")}\")",
    "technical_support" => $"\n=> create_tech_ticket(issue=\"{e.GetValueOrDefault("issue")}\", platform=\"{e.GetValueOrDefault("platform")}\")",
    "billing_issue"     => $"\n=> flag_billing_dispute(transaction_id=\"{e.GetValueOrDefault("transaction_id")}\", problem=\"{e.GetValueOrDefault("problem")}\")",
    "feature_request"   => $"\n=> log_feature_request(feature=\"{e.GetValueOrDefault("feature")}\")",
    _                   => $"\n=> send_to_faq_system(question=\"{e.GetValueOrDefault("question")}\")",
});
