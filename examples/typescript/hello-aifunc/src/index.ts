import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

// const config: AIFuncConfig = {
//   baseURL: 'https://your-api-endpoint/v1',
//   model: 'your-model-name',
//   apiKey: 'your-api-key',
// };

// To use a real model, replace the line below with the commented config above.
const config: AIFuncConfig = { mock: true };

if (config.mock) {
  console.log(
    'Notice: You are using mock mode for offline testing. ' +
    'Configure a real model for the full experience. Continuing with mock responses...'
  );
}

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a patch of sky smaller than a grain of sand held ' +
  "at arm's length. The images show galaxies as they appeared over 13 billion years ago, " +
  'providing a glimpse into the early universe shortly after the Big Bang.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Original  : ${text}`);
  console.log(`Summary   : ${result.summary}`);
  console.log(`Word count: ${result.wordCount}`);
}

main().catch(console.error);
