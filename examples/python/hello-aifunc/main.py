import sys
import io
import asyncio
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")

from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput

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

text = (
    "The James Webb Space Telescope captured its first full-color images in July 2022, "
    "revealing thousands of galaxies in a patch of sky smaller than a grain of sand held "
    "at arm's length. The images show galaxies as they appeared over 13 billion years ago, "
    "providing a glimpse into the early universe shortly after the Big Bang."
)


async def main():
    result = await summarize(config, SummarizeInput(text=text, max_length=30))
    print(f"Original  : {text}")
    print(f"Summary   : {result.summary}")
    print(f"Word count: {result.word_count}")


asyncio.run(main())
