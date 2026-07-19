using Aifunc;
using Aifunc.ChatStream;

// var config = new AIFuncConfig
// {
//     BaseUrl = "https://your-api-endpoint/v1",
//     Model = "your-model-name",
//     ApiKey = "your-api-key",
//     MaxRetries = 3,
// };

// To use a real model, replace the line below with the commented config above.
var config = new AIFuncConfig { Mock = true };

var inputShort = new ChatStreamTypes.ChatStreamInput(
    message: "What is the difference between a process and a thread? Answer in 3 sentences.");

var inputWithContext = new ChatStreamTypes.ChatStreamInput(
    message: "Should I prefer threads or processes for CPU-bound work on multi-core machines?",
    context: "Conversation history:\n"
           + "User: What is the difference between a process and a thread?\n"
           + "Assistant: Processes have separate memory; threads share an address space.");

var inputLong = new ChatStreamTypes.ChatStreamInput(
    message: "Explain the entire history of the internet from ARPANET to today, in detail.");

// ── Short reply: run to completion ──────────────────────────────────────────
Console.WriteLine("--- short reply (run to completion) ---");
await foreach (var token in ChatStream.ChatStreamAsync(config, inputShort))
{
    Console.Write(token);
}
Console.Write("\n\n");

// ── Follow-up with context ──────────────────────────────────────────────────
Console.WriteLine("--- reply with context ---");
await foreach (var token in ChatStream.ChatStreamAsync(config, inputWithContext))
{
    Console.Write(token);
}
Console.Write("\n\n");

// ── Long reply: cancel after 500 characters ──────────────────────────────────
Console.WriteLine("--- long reply (cancel after 500 chars) ---");
using var cts = new CancellationTokenSource();
int chars     = 0;

await foreach (var token in ChatStream.ChatStreamAsync(config, inputLong, cts.Token))
{
    Console.Write(token);
    chars += token.Length;
    if (chars >= 500)
    {
        cts.Cancel();
        break;
    }
}
Console.Write("\n[cancelled]\n");
