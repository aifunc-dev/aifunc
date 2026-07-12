import aifunc.AIFuncConfig;
import aifunc.extract_keywords.ExtractKeywords;
import aifunc.extract_keywords.ExtractKeywordsTypes.ExtractKeywordsInput;
import aifunc.generate_reply.GenerateReply;
import aifunc.generate_reply.GenerateReplyTypes.GenerateReplyInput;
import aifunc.recognize_intent.RecognizeIntent;
import aifunc.recognize_intent.RecognizeIntentTypes.RecognizeIntentInput;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;

import java.util.ArrayList;
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

    static final int WINDOW = 4;
    static final int COMPRESS_AFTER = 6;

    static final List<String> INTENTS = List.of(
            "ask_recommendation",
            "ask_logistics",
            "ask_budget",
            "share_preference",
            "confirm",
            "other"
    );

    static final List<String> MESSAGES = List.of(
            "I'm planning a three-week trip across Europe in September. Where should I start?",
            "I enjoy hiking and local markets. I'd rather skip the big tourist traps.",
            "What's the best way to get from Paris to Barcelona — high-speed train or budget flight?",
            "How much should I budget per day for food and transport in Western Europe?",
            "I've heard the Dolomites are stunning in early autumn. Is it worth a detour from Venice?",
            "Alright, I think I'll do Paris → Barcelona → Rome → Venice → Dolomites. Does that route make sense?"
    );

    record Message(String role, String text) {}

    static final List<Message> history = new ArrayList<>();
    static final List<String> topics = new ArrayList<>();
    static final List<String> intents = new ArrayList<>();
    static String memorySummary = "";

    public static void main(String[] args) {
        if (config.isMock()) {
            System.out.println(
                    "This example requires a real LLM to produce meaningful results.\n"
                            + "Mock mode cannot simulate intent-aware replies grounded in conversation history.\n"
                            + "\n"
                            + "To run this example, replace the line:\n"
                            + "\n"
                            + "  static final AIFuncConfig config = AIFuncConfig.builder().mock(true).build();\n"
                            + "\n"
                            + "with the commented real-credentials config above.\n"
            );
            return;
        }

        for (int i = 0; i < MESSAGES.size(); i++) {
            String userMsg = MESSAGES.get(i);
            System.out.println("[Turn " + (i + 1) + "] User: " + userMsg);

            // 1. Classify the user's intent
            var intentResult = RecognizeIntent.recognizeIntent(config,
                    new RecognizeIntentInput(INTENTS, userMsg, null)).join();
            intents.add(intentResult.getIntent());

            // 2. Extract keywords and accumulate into the topics array
            var kwResult = ExtractKeywords.extractKeywords(config,
                    new ExtractKeywordsInput(userMsg, 3)).join();
            for (Map<String, Object> kw : kwResult.getKeywords()) {
                Object wordObj = kw.get("word");
                if (wordObj == null) continue;
                String word = String.valueOf(wordObj);
                if (!topics.contains(word)) {
                    topics.add(word);
                }
            }

            // 3. Append user turn to history array
            history.add(new Message("user", userMsg));

            // 4. Compress old history before replying if it has grown too long
            maybeCompress();

            // 5. Build context from memory arrays and generate a reply
            String ctx = buildContext();
            var replyResult = GenerateReply.generateReply(config,
                    new GenerateReplyInput(userMsg, ctx, null, "friendly")).join();

            // 6. Append assistant reply to history array
            history.add(new Message("assistant", replyResult.getReply()));

            System.out.printf("         Intent  : %s (%.0f%%)%n",
                    intentResult.getIntent(), intentResult.getConfidence() * 100);
            System.out.println("         Topics  : " + topics);
            System.out.println("         Reply   : " + replyResult.getReply());
            System.out.println();
        }

        System.out.println("=".repeat(60));
        System.out.println("Final memory state");
        System.out.println("  topics  : " + topics);
        System.out.println("  intents : " + intents);
        if (!memorySummary.isEmpty()) {
            System.out.println("  summary : " + memorySummary);
        }
        System.out.println("  history : " + history.size() + " turns in window");
    }

    private static String buildContext() {
        List<String> parts = new ArrayList<>();

        if (!memorySummary.isEmpty()) {
            parts.add("Earlier in this conversation: " + memorySummary);
        }

        List<Message> recent = history.size() <= WINDOW
                ? history
                : history.subList(history.size() - WINDOW, history.size());
        if (!recent.isEmpty()) {
            StringBuilder dialogue = new StringBuilder();
            for (Message m : recent) {
                String who = "user".equals(m.role()) ? "User" : "Assistant";
                if (dialogue.length() > 0) dialogue.append('\n');
                dialogue.append("  ").append(who).append(": ").append(m.text());
            }
            parts.add("Recent conversation:\n" + dialogue);
        }

        if (!topics.isEmpty()) {
            parts.add("Topics discussed so far: " + String.join(", ", topics));
        }

        if (!intents.isEmpty()) {
            List<String> recentIntents = intents.size() <= 4
                    ? intents
                    : intents.subList(intents.size() - 4, intents.size());
            parts.add("User intent pattern: " + String.join(" → ", recentIntents));
        }

        return String.join("\n\n", parts);
    }

    private static void maybeCompress() {
        if (history.size() <= COMPRESS_AFTER) {
            return;
        }

        List<Message> older = history.subList(0, history.size() - WINDOW);
        StringBuilder olderText = new StringBuilder();
        for (Message m : older) {
            if (olderText.length() > 0) olderText.append(' ');
            olderText.append(m.text());
        }

        var result = Summarize.summarize(config, new SummarizeInput(olderText.toString(), 40)).join();
        memorySummary = result.getSummary();

        List<Message> keep = new ArrayList<>(history.subList(history.size() - WINDOW, history.size()));
        history.clear();
        history.addAll(keep);
        System.out.println("  [memory compressed → \"" + memorySummary + "\"]\n");
    }
}
