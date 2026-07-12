import aifunc.AIFuncConfig;
import aifunc.analyze_sentiment.AnalyzeSentiment;
import aifunc.analyze_sentiment.AnalyzeSentimentTypes.AnalyzeSentimentInput;
import aifunc.extract_json.ExtractJson;
import aifunc.extract_json.ExtractJsonTypes.ExtractJsonInput;
import aifunc.recognize_intent.RecognizeIntent;
import aifunc.recognize_intent.RecognizeIntentTypes.RecognizeIntentInput;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ThreadLocalRandom;

public class Main {
    // AIFuncConfig config = AIFuncConfig.builder()
    //         .baseUrl("https://your-api-endpoint/v1")
    //         .model("your-model-name")
    //         .apiKey("your-api-key")
    //         .maxRetries(3)
    //         .build();

    // To use a real model, replace the line below with the commented config above.
    static final AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

    // Tip: Each call accepts its own config — use cheaper models for simple tasks to save cost.
    //
    // AIFuncConfig cheap = AIFuncConfig.builder().baseUrl("...").model("...").apiKey("...").build();
    // AIFuncConfig strong = AIFuncConfig.builder().baseUrl("...").model("...").apiKey("...").build();
    //
    // AnalyzeSentiment.analyzeSentiment(cheap, ...);   // classification is simple, cheap model is fine
    // ExtractJson.extractJson(strong, ...);            // extraction needs accuracy, use a stronger model

    static final List<String> MESSAGES = List.of(
            "What the hell?! I ordered this a WEEK ago and it still hasn't shipped! I want my money back NOW!",
            "Hi, I'd like to check on my order #ORD-20240601-123. It's been three days with no shipping update.",
            "Your stupid app crashed again and I lost all my data! Fix it or I'm leaving!",
            "I was charged twice this month. Transaction IDs: TXN-88201 and TXN-88202. Please help.",
            "It would be cool if you added a dark mode. The bright screen hurts my eyes at night.",
            "How do I export my purchase history to CSV? I can't find the option.",
            "I am SO FURIOUS! Your delivery guy threw my package over the fence and it's destroyed! I want a manager NOW!",
            "Any ongoing promotions for loyal customers? I've been a member for 2 years."
    );

    static final List<String> INTENTS = List.of(
            "query_order",
            "request_refund",
            "technical_support",
            "billing_issue",
            "feature_request",
            "general_inquiry"
    );

    public static void main(String[] args) {
        if (config.isMock()) {
            System.out.println(
                    "This example requires a real LLM to produce meaningful results.\n"
                            + "Mock mode cannot simulate multi-step reasoning (sentiment → intent → extraction).\n"
                            + "\n"
                            + "To run this example, replace the line:\n"
                            + "\n"
                            + "  static final AIFuncConfig config = AIFuncConfig.builder().mock(true).build();\n"
                            + "\n"
                            + "with the commented real-credentials config above.\n"
            );
            return;
        }

        String message = MESSAGES.get(ThreadLocalRandom.current().nextInt(MESSAGES.size()));
        System.out.println("Customer: " + message + "\n");

        // Step 1: Sentiment analysis
        var sentiment = AnalyzeSentiment.analyzeSentiment(config, new AnalyzeSentimentInput(
                message,
                List.of("angry", "frustrated", "neutral", "happy", "other"),
                null
        )).join();
        System.out.printf("Sentiment: %s (%.0f%%)%n",
                sentiment.getLabel(), sentiment.getConfidence() * 100);

        if ("angry".equals(sentiment.getLabel()) && sentiment.getConfidence() > 0.7) {
            System.out.println("\n=> call_human_agent(message, priority=\"HIGH\")");
            return;
        }

        // Step 2: Intent recognition
        var intentResult = RecognizeIntent.recognizeIntent(config,
                new RecognizeIntentInput(INTENTS, message, null)).join();
        String intent = intentResult.getIntent();
        System.out.printf("Intent: %s (%.0f%%)%n", intent, intentResult.getConfidence() * 100);

        // Step 3: Route by intent
        switch (intent) {
            case "query_order" -> {
                var info = ExtractJson.extractJson(config, new ExtractJsonInput(List.of(
                        field("order_id", "Order number", "string"),
                        field("issue", "What the customer wants to know", "string")
                ), message)).join();
                System.out.printf("%n=> query_order_system(order_id=\"%s\", issue=\"%s\")%n",
                        info.getExtracted().get("order_id"), info.getExtracted().get("issue"));
            }
            case "request_refund" -> {
                var info = ExtractJson.extractJson(config, new ExtractJsonInput(List.of(
                        field("order_id", "Order number", "string"),
                        field("reason", "Reason for refund", "string")
                ), message)).join();
                System.out.printf("%n=> submit_refund(order_id=\"%s\", reason=\"%s\")%n",
                        info.getExtracted().get("order_id"), info.getExtracted().get("reason"));
            }
            case "technical_support" -> {
                var info = ExtractJson.extractJson(config, new ExtractJsonInput(List.of(
                        field("issue", "Technical problem", "string"),
                        field("platform", "Device or platform", "string")
                ), message)).join();
                System.out.printf("%n=> create_tech_ticket(issue=\"%s\", platform=\"%s\")%n",
                        info.getExtracted().get("issue"), info.getExtracted().get("platform"));
            }
            case "billing_issue" -> {
                var info = ExtractJson.extractJson(config, new ExtractJsonInput(List.of(
                        field("transaction_id", "Transaction ID", "string"),
                        field("problem", "Billing problem", "string")
                ), message)).join();
                System.out.printf("%n=> flag_billing_dispute(transaction_id=\"%s\", problem=\"%s\")%n",
                        info.getExtracted().get("transaction_id"), info.getExtracted().get("problem"));
            }
            case "feature_request" -> {
                var info = ExtractJson.extractJson(config, new ExtractJsonInput(List.of(
                        field("feature", "Requested feature", "string")
                ), message)).join();
                System.out.printf("%n=> log_feature_request(feature=\"%s\")%n",
                        info.getExtracted().get("feature"));
            }
            default -> {
                var info = ExtractJson.extractJson(config, new ExtractJsonInput(List.of(
                        field("question", "Customer's question", "string")
                ), message)).join();
                System.out.printf("%n=> send_to_faq_system(question=\"%s\")%n",
                        info.getExtracted().get("question"));
            }
        }
    }

    private static Map<String, Object> field(String name, String description, String type) {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("name", name);
        m.put("description", description);
        m.put("type", type);
        return m;
    }
}
