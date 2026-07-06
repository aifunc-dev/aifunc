import sys
import io
import asyncio
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")

from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput
from aifunc.translate import translate, TranslateInput
from aifunc.analyze_sentiment import analyze_sentiment, AnalyzeSentimentInput
from aifunc.detect_language import detect_language, DetectLanguageInput
from aifunc.rewrite import rewrite, RewriteInput
from aifunc.extract_keywords import extract_keywords, ExtractKeywordsInput
from aifunc.classify import classify, ClassifyInput
from aifunc.recognize_intent import recognize_intent, RecognizeIntentInput
from aifunc.extract_entities import extract_entities, ExtractEntitiesInput
from aifunc.extract_json import extract_json, ExtractJsonInput
from aifunc.generate_slug import generate_slug, GenerateSlugInput
from aifunc.generate_reply import generate_reply, GenerateReplyInput
from aifunc.generate_post import generate_post, GeneratePostInput
from aifunc.generate_email import generate_email, GenerateEmailInput
from aifunc.generate_title import generate_title, GenerateTitleInput
from aifunc.answer_question import answer_question, AnswerQuestionInput
from aifunc.score_quality import score_quality, ScoreQualityInput

# config = AIFuncConfig(
#     base_url="https://your-api-endpoint/v1",
#     model="your-model-name",
#     api_key="your-api-key",
# )

# To use a real model, replace the line below with the commented config above.
config = AIFuncConfig(mock=True)

if config.mock:
    print(
        "Notice: You are using mock mode for offline testing. "
        "Configure a real model for the full experience. Continuing with mock responses..."
    )

def section(title):
    print(f"\n{'=' * 60}")
    print(f"  {title}")
    print(f"{'=' * 60}")

async def main():
    ARTICLE = (
        "In 1915, Albert Einstein published the General Theory of Relativity, "
        "fundamentally transforming our understanding of physics. The theory posits "
        "that gravity is not an invisible force, but rather a curvature of spacetime "
        "caused by the presence of mass and energy. This groundbreaking framework "
        "revolutionized modern science and introduced the famous equation E=mc²."
    )

    # ─── Easy: single input → single output ───────────────────────────

    section("1. DETECT LANGUAGE")
    samples = [
        "The quick brown fox jumps over the lazy dog.",
        "Der schnelle braune Fuchs springt über den faulen Hund.",
        "Le renard brun rapide saute par-dessus le chien paresseux.",
        "El veloz zorro marrón salta sobre el perro perezoso.",
    ]
    for text in samples:
        try:
            result = await detect_language(config, DetectLanguageInput(text=text))
            print(f"  [{result.language}] {result.language_name} (conf: {result.confidence:.0%})  \"{text[:40]}\"")
        except Exception as e:
            print(f"  Error: {e}")

    section("2. GENERATE SLUG")
    try:
        result = await generate_slug(config, GenerateSlugInput(
            title="10 Practical Tips for Writing Faster Python Code",
            language="en",
        ))
        print(f"Title : 10 Practical Tips for Writing Faster Python Code")
        print(f"Slug  : {result.slug}")
        print(f"Meta  : {result.meta_description}")
        print(f"Tags  : {result.tags}")
    except Exception as e:
        print(f"Error: {e}")

    section("3. SUMMARIZE")
    try:
        result = await summarize(config, SummarizeInput(text=ARTICLE, max_length=30))
        print(f"Summary   : {result.summary}")
        print(f"Word count: {result.word_count}")
    except Exception as e:
        print(f"Error: {e}")

    section("4. TRANSLATE")
    try:
        result = await translate(config, TranslateInput(
            text="The meeting has been moved to Friday at 3 PM.",
            target_lang="es",
        ))
        print(f"Original : The meeting has been moved to Friday at 3 PM.")
        print(f"Spanish  : {result.translation}")
        print(f"Detected : {result.source_lang}")
    except Exception as e:
        print(f"Error: {e}")

    section("5. REWRITE")
    try:
        original = "hey, just wanna let u know the deploy went fine, no issues at all"
        formal = await rewrite(config, RewriteInput(text=original, style="formal"))
        print(f"Casual : {original}")
        print(f"Formal : {formal.rewritten}")
    except Exception as e:
        print(f"Error: {e}")

    # ─── Medium: structured output or multiple parameters ─────────────

    section("6. GENERATE TITLE")
    try:
        content = (
            "This guide covers how to use Docker and GitHub Actions to automate "
            "testing and deployment of a Node.js application to a cloud server."
        )
        result = await generate_title(config, GenerateTitleInput(content=content, style="seo", count=4))
        print(f"Content: {content}")
        print("Titles:")
        for i, title in enumerate(result.titles, 1):
            print(f"  {i}. {title}")
    except Exception as e:
        print(f"Error: {e}")

    section("7. EXTRACT KEYWORDS")
    try:
        result = await extract_keywords(config, ExtractKeywordsInput(text=ARTICLE, max_keywords=5))
        print("Keywords from article:")
        for kw in result.keywords:
            print(f"  {kw['word']:30s} relevance: {kw['relevance']}")
    except Exception as e:
        print(f"Error: {e}")

    section("8. ANALYZE SENTIMENT")
    samples = [
        "The product arrived on time and works perfectly. Very happy!",
        "Terrible experience. The package was damaged and support ignored my emails.",
        "Item received. Does what it says.",
    ]
    for text in samples:
        try:
            result = await analyze_sentiment(config, AnalyzeSentimentInput(
                text=text,
                labels=["positive", "negative", "neutral"],
            ))
            print(f"  [{result.label:8s} {result.confidence:.0%}] {text[:55]}")
        except Exception as e:
            print(f"  Error: {e}")

    section("9. CLASSIFY")
    tickets = [
        "My order hasn't shipped after five days. Please help.",
        "The API returns a 500 error when the payload exceeds 1 MB.",
        "It would be great to have a dark mode option.",
        "I was charged twice for the same subscription this month.",
    ]
    categories = ["shipping", "technical", "feature request", "billing", "other"]
    for ticket in tickets:
        try:
            result = await classify(config, ClassifyInput(text=ticket, categories=categories))
            top = result.classifications[0]
            print(f"  [{top['category']:16s} {top['confidence']:.0%}]  {ticket[:55]}")
        except Exception as e:
            print(f"  Error: {e}")

    section("10. RECOGNIZE INTENT")
    messages = [
        "Where is my order? I placed it three days ago.",
        "I want a refund for the broken item.",
        "Can you tell me your business hours?",
        "I'd like to upgrade my subscription to the pro plan.",
    ]
    intents = ["query_order", "request_refund", "general_inquiry", "manage_subscription"]
    for msg in messages:
        try:
            result = await recognize_intent(config, RecognizeIntentInput(
                text=msg,
                intents=intents,
                context="You are a customer support routing system for an e-commerce platform.",
            ))
            print(f"  [{result.intent:20s} {result.confidence:.0%}]  \"{msg[:50]}\"")
        except Exception as e:
            print(f"  Error: {e}")

    # ─── Advanced: complex extraction and generation ──────────────────

    section("11. EXTRACT ENTITIES")
    try:
        text = "On March 10, 2024, NASA astronaut Sarah Mitchell landed at Kennedy Space Center in Florida after a six-month mission."
        result = await extract_entities(config, ExtractEntitiesInput(
            text=text,
            entity_types=["person", "organization", "location", "date"],
        ))
        print(f"Text: {text}")
        print("Entities:")
        for e in result.entities:
            print(f"  [{e['type']:12s}] \"{e['text']}\"")
    except Exception as e:
        print(f"Error: {e}")

    section("12. EXTRACT JSON")
    try:
        job_post = (
            "We are looking for a Senior Backend Engineer in Berlin. "
            "Requirements: 5+ years of experience, proficiency in Go or Rust, "
            "experience with Kubernetes. Salary range: €80,000–€110,000."
        )
        result = await extract_json(config, ExtractJsonInput(
            text=job_post,
            fields=[
                {"name": "title", "description": "Job title", "type": "string"},
                {"name": "location", "description": "City or country", "type": "string"},
                {"name": "skills", "description": "Required technical skills", "type": "array"},
                {"name": "experience_years", "description": "Minimum years of experience", "type": "number"},
                {"name": "salary_range", "description": "Salary range", "type": "string"},
            ],
        ))
        print(f"Text     : {job_post}")
        print(f"Extracted: {result.extracted}")
        print(f"Missing  : {result.missing}")
    except Exception as e:
        print(f"Error: {e}")

    section("13. ANSWER QUESTION")
    context = (
        "AIFunc is a function-based AI toolkit. Developers declare the packages they need "
        "in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. "
        "Each package supports a mock mode for testing without consuming API credits."
    )
    questions = [
        ("Which languages does AIFunc support?", context),
        ("What is mock mode used for?", context),
        ("What is a monad in functional programming?", None),
    ]
    for q, ctx in questions:
        try:
            result = await answer_question(config, AnswerQuestionInput(question=q, context=ctx, max_length=60))
            source = "from context" if result.grounded else "general knowledge"
            print(f"  Q: {q}")
            print(f"  A: {result.answer}  [{source}, conf: {result.confidence:.0%}]")
            print()
        except Exception as e:
            print(f"  Error: {e}")
            print()

    section("14. GENERATE REPLY")
    try:
        message = "I placed an order three days ago but haven't received a shipping confirmation yet."
        result = await generate_reply(config, GenerateReplyInput(
            message=message,
            tone="empathetic",
            context="You are a customer support agent for an online store.",
        ))
        print(f"Customer : {message}")
        print(f"Reply    : {result.reply}")
    except Exception as e:
        print(f"Error: {e}")

    section("15. GENERATE POST")
    try:
        result = await generate_post(config, GeneratePostInput(
            topic="How switching to async Python cut our API response time by 60%",
            platform="linkedin",
            tone="professional",
            include_hashtags=True,
        ))
        print(f"Post     : {result.post}")
        print(f"Hashtags : {['#' + t for t in result.hashtags]}")
    except Exception as e:
        print(f"Error: {e}")

    section("16. GENERATE EMAIL")
    try:
        result = await generate_email(config, GenerateEmailInput(
            intent="Apologize to a customer for a billing error and explain the resolution",
            tone="formal",
            sender_name="Billing Support Team",
            recipient_name="Alex",
            key_points=[
                "An incorrect charge of $29.99 was applied on June 1st",
                "The charge has been fully refunded and will appear within 3–5 business days",
                "We have applied a 20% discount to the next invoice as compensation",
            ],
            language="English",
        ))
        print(f"Subject: {result.subject}")
        print(f"Body:\n{result.body}")
    except Exception as e:
        print(f"Error: {e}")

    section("17. SCORE QUALITY")
    samples = [
        (
            "Our product is good. It has many features. Users like it.",
            "customers",
            "marketing",
        ),
        (
            "To set up the CI pipeline: 1) Install Docker and the GitHub CLI. "
            "2) Create .github/workflows/deploy.yml. 3) Push to main to trigger deployment.",
            "developers",
            "explanation",
        ),
    ]
    for text, audience, purpose in samples:
        try:
            result = await score_quality(config, ScoreQualityInput(
                text=text,
                target_audience=audience,
                purpose=purpose,
                max_suggestions=3,
                strictness=3,
            ))
            print(f"Text       : {text[:55]}...")
            print(f"Score      : {result.overall_score}/100  [{result.level}]")
            print(f"Summary    : {result.summary}")
            print(f"Suggestions:")
            for s in result.suggestions:
                print(f"  - {s}")
            print()
        except Exception as e:
            print(f"Error: {e}")
            print()
    if config.mock:
        print(
            "Notice: You are using mock mode for offline testing. "
            "Configure a real model for the full experience."
        )


if __name__ == "__main__":
    asyncio.run(main())
