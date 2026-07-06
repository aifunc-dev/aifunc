import sys
import io
import asyncio
import random
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")

from aifunc.analyze_sentiment import analyze_sentiment, AIFuncConfig, AnalyzeSentimentInput
from aifunc.recognize_intent import recognize_intent, RecognizeIntentInput
from aifunc.extract_json import extract_json, ExtractJsonInput

# config = AIFuncConfig(
#     base_url="https://your-api-endpoint/v1",
#     model="your-model-name",
#     api_key="your-api-key",
# )

# To run this example, replace the mock config below with real credentials:
config = AIFuncConfig(mock=True)

# Tip: Each call accepts its own config — use cheaper models for simple tasks to save cost.
#
# cheap = AIFuncConfig(base_url="...", model="...", api_key="...")
# strong = AIFuncConfig(base_url="...", model="...", api_key="...")
#
# await analyze_sentiment(cheap, ...)   # classification is simple, cheap model is fine
# await extract_json(strong, ...)       # extraction needs accuracy, use a stronger model

if config.mock:
    print(
        "This example requires a real LLM to produce meaningful results.\n"
        "Mock mode cannot simulate multi-step reasoning (sentiment → intent → extraction).\n"
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

MESSAGES = [
    "What the hell?! I ordered this a WEEK ago and it still hasn't shipped! I want my money back NOW!",
    "Hi, I'd like to check on my order #ORD-20240601-123. It's been three days with no shipping update.",
    "Your stupid app crashed again and I lost all my data! Fix it or I'm leaving!",
    "I was charged twice this month. Transaction IDs: TXN-88201 and TXN-88202. Please help.",
    "It would be cool if you added a dark mode. The bright screen hurts my eyes at night.",
    "How do I export my purchase history to CSV? I can't find the option.",
    "I am SO FURIOUS! Your delivery guy threw my package over the fence and it's destroyed! I want a manager NOW!",
    "Any ongoing promotions for loyal customers? I've been a member for 2 years.",
]


async def retry(fn, label, retries=3):
    """Call an async function with retries. Gives up after `retries` failures."""
    for attempt in range(retries):
        try:
            return await fn()
        except Exception as e:
            if attempt == retries - 1:
                raise
            print(f"  [retry {label} ({attempt+1}/{retries}): {e}]")
            await asyncio.sleep(1)


async def main():
    message = random.choice(MESSAGES)
    print(f"Customer: {message}\n")

    # Step 1: Sentiment analysis
    sentiment = await retry(
        lambda: analyze_sentiment(config, AnalyzeSentimentInput(
            text=message,
            labels=["angry", "frustrated", "neutral", "happy", "other"],
        )),
        "analyze_sentiment",
    )
    print(f"Sentiment: {sentiment.label} ({sentiment.confidence:.0%})")

    if sentiment.label == "angry" and sentiment.confidence > 0.7:
        print(f"\n=> call_human_agent(message, priority=\"HIGH\")")
        return

    # Step 2: Intent recognition
    INTENTS = ["query_order", "request_refund", "technical_support", "billing_issue", "feature_request", "general_inquiry"]
    intent_result = await retry(
        lambda: recognize_intent(config, RecognizeIntentInput(text=message, intents=INTENTS)),
        "recognize_intent",
    )
    intent = intent_result.intent
    print(f"Intent: {intent} ({intent_result.confidence:.0%})")

    # Step 3: Route by intent
    match intent:
        case "query_order":
            info = await retry(lambda: extract_json(config, ExtractJsonInput(text=message, fields=[
                {"name": "order_id", "description": "Order number", "type": "string"},
                {"name": "issue", "description": "What the customer wants to know", "type": "string"},
            ])), "extract_json")
            print(f"\n=> query_order_system(order_id=\"{info.extracted.get('order_id')}\", issue=\"{info.extracted.get('issue')}\")")

        case "request_refund":
            info = await retry(lambda: extract_json(config, ExtractJsonInput(text=message, fields=[
                {"name": "order_id", "description": "Order number", "type": "string"},
                {"name": "reason", "description": "Reason for refund", "type": "string"},
            ])), "extract_json")
            print(f"\n=> submit_refund(order_id=\"{info.extracted.get('order_id')}\", reason=\"{info.extracted.get('reason')}\")")

        case "technical_support":
            info = await retry(lambda: extract_json(config, ExtractJsonInput(text=message, fields=[
                {"name": "issue", "description": "Technical problem", "type": "string"},
                {"name": "platform", "description": "Device or platform", "type": "string"},
            ])), "extract_json")
            print(f"\n=> create_tech_ticket(issue=\"{info.extracted.get('issue')}\", platform=\"{info.extracted.get('platform')}\")")

        case "billing_issue":
            info = await retry(lambda: extract_json(config, ExtractJsonInput(text=message, fields=[
                {"name": "transaction_id", "description": "Transaction ID", "type": "string"},
                {"name": "problem", "description": "Billing problem", "type": "string"},
            ])), "extract_json")
            print(f"\n=> flag_billing_dispute(transaction_id=\"{info.extracted.get('transaction_id')}\", problem=\"{info.extracted.get('problem')}\")")

        case "feature_request":
            info = await retry(lambda: extract_json(config, ExtractJsonInput(text=message, fields=[
                {"name": "feature", "description": "Requested feature", "type": "string"},
            ])), "extract_json")
            print(f"\n=> log_feature_request(feature=\"{info.extracted.get('feature')}\")")

        case _:
            info = await retry(lambda: extract_json(config, ExtractJsonInput(text=message, fields=[
                {"name": "question", "description": "Customer's question", "type": "string"},
            ])), "extract_json")
            print(f"\n=> send_to_faq_system(question=\"{info.extracted.get('question')}\")")


asyncio.run(main())
