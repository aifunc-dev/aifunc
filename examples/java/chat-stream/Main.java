import aifunc.AIFuncConfig;
import aifunc.chat_stream.ChatStream;
import aifunc.chat_stream.ChatStreamTypes.ChatStreamInput;

public class Main {

    // AIFuncConfig config = AIFuncConfig.builder()
    //         .baseUrl("https://your-api-endpoint/v1")
    //         .model("your-model-name")
    //         .apiKey("your-api-key")
    //         .maxRetries(3)
    //         .build();

    // To use a real model, replace the line below with the commented config above.
    static final AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

    static final ChatStreamInput inputShort = new ChatStreamInput(
            "What is the difference between a process and a thread? Answer in 3 sentences.",
            null
    );

    static final ChatStreamInput inputWithContext = new ChatStreamInput(
            "Should I prefer threads or processes for CPU-bound work on multi-core machines?",
            "Conversation history:\n"
                    + "User: What is the difference between a process and a thread?\n"
                    + "Assistant: Processes have separate memory; threads share an address space."
    );

    static final ChatStreamInput inputLong = new ChatStreamInput(
            "Explain the entire history of the internet from ARPANET to today, in detail.",
            null
    );

    public static void main(String[] args) {
        // Short reply — run to completion
        System.out.println("--- short reply (run to completion) ---");
        try (var tokens = ChatStream.chatStream(config, inputShort)) {
            while (tokens.hasNext()) {
                System.out.print(tokens.next());
            }
        }
        System.out.print("\n\n");

        // Follow-up with context
        System.out.println("--- reply with context ---");
        try (var tokens = ChatStream.chatStream(config, inputWithContext)) {
            while (tokens.hasNext()) {
                System.out.print(tokens.next());
            }
        }
        System.out.print("\n\n");

        // Long reply — cancel after 500 characters
        System.out.println("--- long reply (cancel after 500 chars) ---");
        try (var tokens = ChatStream.chatStream(config, inputLong)) {
            int chars = 0;
            while (tokens.hasNext()) {
                String token = tokens.next();
                System.out.print(token);
                chars += token.length();
                if (chars >= 500) { break; }
            }
        }
        System.out.print("\n[cancelled]\n");
    }
}
