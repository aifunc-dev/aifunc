import sys
import io
import asyncio
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")

from aifunc.recognize_intent import recognize_intent, AIFuncConfig, RecognizeIntentInput
from aifunc.extract_keywords import extract_keywords, ExtractKeywordsInput
from aifunc.summarize import summarize, SummarizeInput
from aifunc.generate_reply import generate_reply, GenerateReplyInput

# config = AIFuncConfig(
#     base_url="https://your-api-endpoint/v1",
#     model="your-model-name",
#     api_key="your-api-key",
#     max_retries=3,
# )

# To run this example, replace the mock config below with real credentials:
config = AIFuncConfig(mock=True)

if config.mock:
    print(
        "This example requires a real LLM to produce meaningful results.\n"
        "Mock mode cannot simulate intent-aware replies grounded in conversation history.\n"
        "\n"
        "To run this example, replace the line:\n"
        "\n"
        "  config = AIFuncConfig(mock=True)\n"
        "\n"
        "with:\n"
        "\n"
        "  config = AIFuncConfig(\n"
        "      base_url=\"https://your-api-endpoint/v1\",\n"
        "      model=\"your-model-name\",\n"
        "      api_key=\"your-api-key\",\n"
        "  )\n"
    )
    sys.exit(0)

history: list[dict] = []
topics: list[str] = []
intents: list[str] = []

# Conversation window: keep the last N turns in context before compressing
WINDOW = 4
COMPRESS_AFTER = 6  # compress history into a summary once it exceeds this length

INTENTS = [
    "ask_recommendation",
    "ask_logistics",
    "ask_budget",
    "share_preference",
    "confirm",
    "other",
]

# Simulated user turns about planning a trip to Europe
messages = [
    "I'm planning a three-week trip across Europe in September. Where should I start?",
    "I enjoy hiking and local markets. I'd rather skip the big tourist traps.",
    "What's the best way to get from Paris to Barcelona — high-speed train or budget flight?",
    "How much should I budget per day for food and transport in Western Europe?",
    "I've heard the Dolomites are stunning in early autumn. Is it worth a detour from Venice?",
    "Alright, I think I'll do Paris → Barcelona → Rome → Venice → Dolomites. Does that route make sense?",
]

# Running compressed summary of older history (set once COMPRESS_AFTER is hit)
memory_summary = ""


async def build_context() -> str:
    """Build a context string for generate-reply from current memory state."""
    parts = []

    if memory_summary:
        parts.append(f"Earlier in this conversation: {memory_summary}")

    # Recent turns within the sliding window
    recent = history[-WINDOW:]
    if recent:
        dialogue = "\n".join(
            f"  {'User' if m['role'] == 'user' else 'Assistant'}: {m['text']}"
            for m in recent
        )
        parts.append(f"Recent conversation:\n{dialogue}")

    if topics:
        parts.append(f"Topics discussed so far: {', '.join(topics)}")

    if intents:
        parts.append(f"User intent pattern: {' → '.join(intents[-4:])}")

    return "\n\n".join(parts)


async def maybe_compress():
    """When history grows long, compress older turns into a summary and trim."""
    global history, memory_summary

    if len(history) <= COMPRESS_AFTER:
        return

    older = history[:-WINDOW]
    older_text = " ".join(m["text"] for m in older)

    result = await summarize(config, SummarizeInput(text=older_text, max_length=40))
    memory_summary = result.summary

    # Keep only the recent window in the live history array
    history = history[-WINDOW:]
    print(f"  [memory compressed → \"{memory_summary}\"]\n")


async def main():
    global history, topics, intents

    for i, user_msg in enumerate(messages, 1):
        print(f"[Turn {i}] User: {user_msg}")

        # 1. Classify the user's intent
        intent_result = await recognize_intent(config, RecognizeIntentInput(text=user_msg, intents=INTENTS))
        intents.append(intent_result.intent)

        # 2. Extract keywords and accumulate into the topics array
        kw_result = await extract_keywords(config, ExtractKeywordsInput(text=user_msg, max_keywords=3))
        for kw in kw_result.keywords:
            word = kw["word"]
            if word not in topics:
                topics.append(word)

        # 3. Append user turn to history array
        history.append({"role": "user", "text": user_msg})

        # 4. Compress old history before replying if it has grown too long
        await maybe_compress()

        # 5. Build context from memory arrays and generate a reply
        ctx = await build_context()
        reply_result = await generate_reply(config, GenerateReplyInput(message=user_msg, tone="friendly", context=ctx))

        # 6. Append assistant reply to history array
        history.append({"role": "assistant", "text": reply_result.reply})

        print(f"         Intent  : {intent_result.intent} ({intent_result.confidence:.0%})")
        print(f"         Topics  : {topics}")
        print(f"         Reply   : {reply_result.reply}")
        print()

    # Final state of memory arrays
    print("=" * 60)
    print("Final memory state")
    print(f"  topics  : {topics}")
    print(f"  intents : {intents}")
    if memory_summary:
        print(f"  summary : {memory_summary}")
    print(f"  history : {len(history)} turns in window")


asyncio.run(main())
