import sys
import io
import asyncio
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")

from aifunc.chat_stream import chat_stream, AIFuncConfig, ChatStreamInput
from aifunc.answer_stream import answer_stream, AnswerStreamInput
from aifunc.explain_stream import explain_stream, ExplainStreamInput
from aifunc.article_stream import article_stream, ArticleStreamInput
from aifunc.write_stream import write_stream, WriteStreamInput
from aifunc.translate_stream import translate_stream, TranslateStreamInput
from aifunc.review_stream import review_stream, ReviewStreamInput

# config = AIFuncConfig(
#     base_url="https://your-api-endpoint/v1",
#     model="your-model-name",
#     api_key="your-api-key",
#     max_retries=3,
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


async def stream_print(gen):
    """Consume a token stream and write each chunk to stdout immediately."""
    async for token in gen:
        sys.stdout.write(token)
        sys.stdout.flush()
    sys.stdout.write("\n")
    sys.stdout.flush()


async def main():
    ARTICLE = (
        "In 1915, Albert Einstein published the General Theory of Relativity, "
        "fundamentally transforming our understanding of physics. The theory posits "
        "that gravity is not an invisible force, but rather a curvature of spacetime "
        "caused by the presence of mass and energy. This groundbreaking framework "
        "revolutionized modern science and introduced the famous equation E=mc²."
    )

    CODE_SNIPPET = '''def fetch_user(user_id):
    conn = get_connection()
    result = conn.execute(f"SELECT * FROM users WHERE id = {user_id}")
    return result.fetchone()
'''

    # ─── Conversational & Q&A ─────────────────────────────────────────

    section("1. CHAT STREAM")
    print("User: Explain async/await in Python in 3 sentences.\n")
    print("Assistant: ", end="", flush=True)
    await stream_print(await chat_stream(config, ChatStreamInput(
        messages=[
            {"role": "user", "content": "Explain async/await in Python in 3 sentences."},
        ],
    )))

    section("2. ANSWER STREAM (with context / RAG)")
    context = (
        "AIFunc is a function-based AI toolkit. Developers declare the packages they need "
        "in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. "
        "Each package supports a mock mode for testing without consuming API credits. "
        "Streaming packages return tokens incrementally via async generators."
    )
    question = "How does AIFunc support offline testing, and what do streaming packages return?"
    print(f"Q: {question}\n")
    print("A: ", end="", flush=True)
    await stream_print(await answer_stream(config, AnswerStreamInput(
        question=question,
        context=context,
        depth="concise",
        audience="technical",
    )))

    section("3. EXPLAIN STREAM")
    print("Topic: the GIL (Global Interpreter Lock) in CPython\n")
    await stream_print(await explain_stream(config, ExplainStreamInput(
        topic="the GIL (Global Interpreter Lock) in CPython",
        audience="intermediate",
        depth="standard",
    )))

    # ─── Long-form writing ────────────────────────────────────────────

    section("4. ARTICLE STREAM")
    title = "Why Typed AI Functions Beat Ad-Hoc Prompt Scripts"
    outline = (
        "- The cost of untyped prompt glue code\n"
        "- How function-shaped AI APIs improve testability\n"
        "- Streaming vs batch for product UX\n"
        "- Practical adoption tips"
    )
    print(f"Title  : {title}")
    print(f"Outline: {outline}\n")
    await stream_print(await article_stream(config, ArticleStreamInput(
        title=title,
        outline=outline,
        style="informational",
        audience="developers",
        word_count=250,
    )))

    section("5. WRITE STREAM")
    prompt = (
        "Write a short internal proposal recommending that our team adopt AIFunc "
        "for customer-support reply generation."
    )
    structure = (
        "1. Problem\n"
        "2. Proposed approach\n"
        "3. Expected benefits\n"
        "4. Next steps"
    )
    print(f"Prompt   : {prompt}")
    print(f"Structure: {structure}\n")
    await stream_print(await write_stream(config, WriteStreamInput(
        prompt=prompt,
        format="proposal",
        structure=structure,
        tone="professional",
        audience="engineers",
        word_count=300,
    )))

    # ─── Translation & review ─────────────────────────────────────────

    section("6. TRANSLATE STREAM")
    print(f"Original (EN):\n{ARTICLE}\n")
    print("Translation (zh-CN):\n")
    await stream_print(await translate_stream(config, TranslateStreamInput(
        text=ARTICLE,
        target_lang="zh-CN",
        style="natural",
        domain="technical",
    )))

    section("7. REVIEW STREAM")
    print(f"Code under review:\n{CODE_SNIPPET}")
    print("Findings:\n")
    await stream_print(await review_stream(config, ReviewStreamInput(
        content=CODE_SNIPPET,
        type="code",
        language="Python",
        focus="correctness, security",
        context="Simple data-access helper in a web API.",
        severity="all",
        output_language="English",
    )))

    if config.mock:
        print(
            "Notice: You are using mock mode for offline testing. "
            "Configure a real model for the full experience."
        )


if __name__ == "__main__":
    asyncio.run(main())
