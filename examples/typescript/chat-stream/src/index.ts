import { chatStream, AIFuncConfig, ChatStreamInput } from './aifunc/chat-stream';

// const config: AIFuncConfig = {
//   baseURL: 'https://your-api-endpoint/v1',
//   model: 'your-model-name',
//   apiKey: 'your-api-key',
//   maxRetries: 3,
// };

// To use a real model, replace the line below with the commented config above.
const config: AIFuncConfig = { mock: true };

const inputShort: ChatStreamInput = {
  message: 'What is the difference between a process and a thread? Answer in 3 sentences.',
};

const inputWithContext: ChatStreamInput = {
  message: 'Should I prefer threads or processes for CPU-bound work on multi-core machines?',
  context:
    'Conversation history:\n' +
    'User: What is the difference between a process and a thread?\n' +
    'Assistant: Processes have separate memory; threads share an address space.',
};

const inputLong: ChatStreamInput = {
  message: 'Explain the entire history of the internet from ARPANET to today, in detail.',
};

async function main() {
  // Short reply — run to completion
  console.log('--- short reply (run to completion) ---');
  for await (const token of chatStream(config, inputShort)) {
    process.stdout.write(token);
  }
  process.stdout.write('\n\n');

  // Follow-up with context
  console.log('--- reply with context ---');
  for await (const token of chatStream(config, inputWithContext)) {
    process.stdout.write(token);
  }
  process.stdout.write('\n\n');

  // Long reply — cancel after 500 characters
  console.log('--- long reply (cancel after 500 chars) ---');
  let chars = 0;
  for await (const token of chatStream(config, inputLong)) {
    process.stdout.write(token);
    if ((chars += token.length) >= 500) {
      break;
    }
  }
  process.stdout.write('\n[cancelled]\n');
}

main().catch(console.error);
