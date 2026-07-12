import aifunc.AIFuncConfig;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;

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
            System.out.println("Notice: You are using mock mode for offline testing. " +
                "Configure a real model for the full experience. Continuing with mock responses...");
        }

        String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
            "revealing thousands of galaxies in a patch of sky smaller than a grain of sand held " +
            "at arm's length. The images show galaxies as they appeared over 13 billion years ago, " +
            "providing a glimpse into the early universe shortly after the Big Bang.";

        Summarize.summarize(config, new SummarizeInput(text, 30))
                .thenAccept(result -> {
                    System.out.println("Original  : " + text);
                    System.out.println("Summary   : " + result.getSummary());
                    System.out.println("Word count: " + result.getWordCount());
                })
                .join();
    }
}
