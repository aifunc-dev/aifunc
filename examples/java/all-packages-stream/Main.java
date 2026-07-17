import aifunc.AIFuncConfig;
import aifunc.answer_stream.AnswerStream;
import aifunc.answer_stream.AnswerStreamTypes.AnswerStreamInput;
import aifunc.article_stream.ArticleStream;
import aifunc.article_stream.ArticleStreamTypes.ArticleStreamInput;
import aifunc.chat_stream.ChatStream;
import aifunc.chat_stream.ChatStreamTypes.ChatStreamInput;
import aifunc.explain_stream.ExplainStream;
import aifunc.explain_stream.ExplainStreamTypes.ExplainStreamInput;
import aifunc.review_stream.ReviewStream;
import aifunc.review_stream.ReviewStreamTypes.ReviewStreamInput;
import aifunc.translate_stream.TranslateStream;
import aifunc.translate_stream.TranslateStreamTypes.TranslateStreamInput;
import aifunc.write_stream.WriteStream;
import aifunc.write_stream.WriteStreamTypes.WriteStreamInput;

import java.util.Iterator;
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

    static void section(String title) {
        System.out.println();
        System.out.println("=".repeat(60));
        System.out.println("  " + title);
        System.out.println("=".repeat(60));
    }

    /** Consume a token stream and write each chunk to stdout immediately. */
    static void streamPrint(Iterator<String> tokens) {
        try {
            while (tokens.hasNext()) {
                System.out.print(tokens.next());
            }
        } finally {
            if (tokens instanceof AutoCloseable closeable) {
                try {
                    closeable.close();
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            }
        }
        System.out.println();
    }

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

        String codeSnippet = """
                def fetch_user(user_id):
                    conn = get_connection()
                    result = conn.execute(f"SELECT * FROM users WHERE id = {user_id}")
                    return result.fetchone()
                """;

        // ─── Conversational & Q&A ─────────────────────────────────────────

        section("1. CHAT STREAM");
        System.out.println("User: Explain streams and try-with-resources in Java in 3 sentences.");
        System.out.println();
        System.out.print("Assistant: ");
        streamPrint(ChatStream.chatStream(config, new ChatStreamInput(
                List.of(Map.of("role", "user",
                        "content", "Explain streams and try-with-resources in Java in 3 sentences.")),
                null, null
        )));

        section("2. ANSWER STREAM (with context / RAG)");
        String context = "AIFunc is a function-based AI toolkit. Developers declare the packages they need "
                + "in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. "
                + "Each package supports a mock mode for testing without consuming API credits. "
                + "Streaming packages return tokens incrementally via a TokenStream iterator.";
        String question = "How does AIFunc support offline testing, and what do streaming packages return?";
        System.out.println("Q: " + question);
        System.out.println();
        System.out.print("A: ");
        streamPrint(AnswerStream.answerStream(config, new AnswerStreamInput(
                question, "technical", context, "concise", null
        )));

        section("3. EXPLAIN STREAM");
        System.out.println("Topic: the JVM garbage collector");
        System.out.println();
        streamPrint(ExplainStream.explainStream(config, new ExplainStreamInput(
                "the JVM garbage collector", "intermediate", null, "standard", null
        )));

        // ─── Long-form writing ────────────────────────────────────────────

        section("4. ARTICLE STREAM");
        String title = "Why Typed AI Functions Beat Ad-Hoc Prompt Scripts";
        String outline = "- The cost of untyped prompt glue code\n"
                + "- How function-shaped AI APIs improve testability\n"
                + "- Streaming vs batch for product UX\n"
                + "- Practical adoption tips";
        System.out.println("Title  : " + title);
        System.out.println("Outline: " + outline);
        System.out.println();
        streamPrint(ArticleStream.articleStream(config, new ArticleStreamInput(
                title, "developers", null, outline, "informational", 250
        )));

        section("5. WRITE STREAM");
        String prompt = "Write a short internal proposal recommending that our team adopt AIFunc "
                + "for customer-support reply generation.";
        String structure = "1. Problem\n"
                + "2. Proposed approach\n"
                + "3. Expected benefits\n"
                + "4. Next steps";
        System.out.println("Prompt   : " + prompt);
        System.out.println("Structure: " + structure);
        System.out.println();
        streamPrint(WriteStream.writeStream(config, new WriteStreamInput(
                prompt, "engineers", "proposal", null, structure, "professional", 300
        )));

        // ─── Translation & review ─────────────────────────────────────────

        section("6. TRANSLATE STREAM");
        System.out.println("Original (EN):\n" + article);
        System.out.println();
        System.out.println("Translation (zh-CN):");
        System.out.println();
        streamPrint(TranslateStream.translateStream(config, new TranslateStreamInput(
                "zh-CN", article, "technical", null, "natural"
        )));

        section("7. REVIEW STREAM");
        System.out.println("Code under review:\n" + codeSnippet);
        System.out.println("Findings:");
        System.out.println();
        streamPrint(ReviewStream.reviewStream(config, new ReviewStreamInput(
                codeSnippet,
                "Simple data-access helper in a web API.",
                "correctness, security",
                "Python",
                "English",
                "all",
                "code"
        )));

        if (config.isMock()) {
            System.out.println("Notice: You are using mock mode for offline testing. "
                    + "Configure a real model for the full experience.");
        }
    }
}
