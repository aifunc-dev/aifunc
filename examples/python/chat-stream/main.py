import asyncio
import sys
from aifunc.chat_stream import chat_stream, AIFuncConfig, ChatStreamInput

# config = AIFuncConfig(
#     base_url="https://your-api-endpoint/v1",
#     model="your-model-name",
#     api_key="your-api-key",
#     max_retries=3,
# )

# To use a real model, replace the line below with the commented config above.
config = AIFuncConfig(mock=True)

# Minimal call — message only
input_short = ChatStreamInput(
    message="What is the difference between a process and a thread? Answer in 3 sentences.",
)

# With optional context (caller-owned history / background)
input_with_context = ChatStreamInput(
    message="Should I prefer threads or processes for CPU-bound work on multi-core machines?",
    context=(
        "Conversation history:\n"
        "User: What is the difference between a process and a thread?\n"
        "Assistant: Processes have separate memory; threads share an address space."
    ),
)

input_long = ChatStreamInput(
    message="Explain the entire history of the internet from ARPANET to today, in detail.",
)


async def main() -> None:
    # Short reply — run to completion
    print("--- short reply (run to completion) ---")
    async for token in await chat_stream(config, input_short):
        sys.stdout.write(token)
        sys.stdout.flush()
    sys.stdout.write("\n\n")

    # Follow-up with context
    print("--- reply with context ---")
    async for token in await chat_stream(config, input_with_context):
        sys.stdout.write(token)
        sys.stdout.flush()
    sys.stdout.write("\n\n")

    # Long reply — cancel after 500 characters
    print("--- long reply (cancel after 500 chars) ---")
    chars = 0
    async for token in await chat_stream(config, input_long):
        sys.stdout.write(token)
        sys.stdout.flush()
        chars += len(token)
        if chars >= 500:
            break
    sys.stdout.write("\n[cancelled]\n")


if __name__ == "__main__":
    asyncio.run(main())
