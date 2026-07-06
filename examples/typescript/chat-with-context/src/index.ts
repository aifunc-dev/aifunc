import { recognizeIntent, AIFuncConfig, RecognizeIntentOutput } from './aifunc/recognize-intent';
import { extractKeywords } from './aifunc/extract-keywords';
import { summarize } from './aifunc/summarize';
import { generateReply } from './aifunc/generate-reply';

// const config: AIFuncConfig = {
//   baseURL: 'https://your-api-endpoint/v1',
//   model: 'your-model-name',
//   apiKey: 'your-api-key',
// };

// To use a real model, replace the line below with the commented config above.
const config: AIFuncConfig = { mock: true };

if (config.mock) {
  console.log(
    'This example requires a real LLM to produce meaningful results.\n' +
    'Mock mode cannot simulate intent-aware replies grounded in conversation history.\n' +
    '\n' +
    'To run this example, replace the line:\n' +
    '\n' +
    '  const config: AIFuncConfig = { mock: true };\n' +
    '\n' +
    'with:\n' +
    '\n' +
    '  const config: AIFuncConfig = {\n' +
    '    baseURL: "https://your-api-endpoint/v1",\n' +
    '    model: "your-model-name",\n' +
    '    apiKey: "your-api-key",\n' +
    '  };\n'
  );
  process.exit(0);
}

// ---------------------------------------------------------------------------
// Memory is just plain arrays — no special memory object needed.
// history   : { role, text } tuples accumulated across turns
// topics    : deduplicated keywords accumulated across all turns
// intents   : intent label per user turn, in order
// ---------------------------------------------------------------------------
type Message = { role: 'user' | 'assistant'; text: string };

const history: Message[] = [];
const topics: string[] = [];
const intents: string[] = [];

const WINDOW = 4;
const COMPRESS_AFTER = 6;

const INTENTS = [
  'ask_recommendation',
  'ask_logistics',
  'ask_budget',
  'share_preference',
  'confirm',
  'other',
];

const messages = [
  "I'm planning a three-week trip across Europe in September. Where should I start?",
  "I enjoy hiking and local markets. I'd rather skip the big tourist traps.",
  "What's the best way to get from Paris to Barcelona — high-speed train or budget flight?",
  'How much should I budget per day for food and transport in Western Europe?',
  "I've heard the Dolomites are stunning in early autumn. Is it worth a detour from Venice?",
  "Alright, I think I'll do Paris → Barcelona → Rome → Venice → Dolomites. Does that route make sense?",
];

let memorySummary = '';

async function buildContext(): Promise<string> {
  const parts: string[] = [];

  if (memorySummary) {
    parts.push(`Earlier in this conversation: ${memorySummary}`);
  }

  const recent = history.slice(-WINDOW);
  if (recent.length > 0) {
    const dialogue = recent
      .map(m => `  ${m.role === 'user' ? 'User' : 'Assistant'}: ${m.text}`)
      .join('\n');
    parts.push(`Recent conversation:\n${dialogue}`);
  }

  if (topics.length > 0) {
    parts.push(`Topics discussed so far: ${topics.join(', ')}`);
  }

  if (intents.length > 0) {
    parts.push(`User intent pattern: ${intents.slice(-4).join(' → ')}`);
  }

  return parts.join('\n\n');
}

async function maybeCompress(): Promise<void> {
  if (history.length <= COMPRESS_AFTER) return;

  const older = history.slice(0, -WINDOW);
  const olderText = older.map(m => m.text).join(' ');

  const result = await summarize(config, { text: olderText, maxLength: 40 });
  memorySummary = result.summary;

  history.splice(0, history.length - WINDOW);
  console.log(`  [memory compressed → "${memorySummary}"]\n`);
}

async function retry<T>(fn: () => Promise<T>, label: string, retries = 3): Promise<T> {
  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      return await fn();
    } catch (e) {
      if (attempt === retries - 1) throw e;
      console.log(`  [retry ${label} (${attempt + 1}/${retries}): ${e}]`);
      await new Promise(r => setTimeout(r, 1000));
    }
  }
  throw new Error('unreachable');
}

async function main(): Promise<void> {
  for (let i = 0; i < messages.length; i++) {
    const userMsg = messages[i];
    console.log(`[Turn ${i + 1}] User: ${userMsg}`);

    // 1. Classify the user's intent
    const intentResult = await retry(
      () => recognizeIntent(config, { text: userMsg, intents: INTENTS }),
      "recognizeIntent",
    );
    intents.push(intentResult.intent);

    // 2. Extract keywords and accumulate into the topics array
    const kwResult = await retry(
      () => extractKeywords(config, { text: userMsg, maxKeywords: 3 }),
      "extractKeywords",
    );
    for (const kw of kwResult.keywords) {
      if (!topics.includes(kw.word)) topics.push(kw.word);
    }

    // 3. Append user turn to history array
    history.push({ role: 'user', text: userMsg });

    // 4. Compress old history before replying if it has grown too long
    await retry(() => maybeCompress(), "maybeCompress");

    // 5. Build context from memory arrays and generate a reply
    const ctx = await buildContext();
    const replyResult = await retry(
      () => generateReply(config, { message: userMsg, tone: 'friendly', context: ctx }),
      "generateReply",
    );

    // 6. Append assistant reply to history array
    history.push({ role: 'assistant', text: replyResult.reply });

    console.log(`         Intent  : ${intentResult.intent} (${Math.round(intentResult.confidence * 100)}%)`);
    console.log(`         Topics  : ${JSON.stringify(topics)}`);
    console.log(`         Reply   : ${replyResult.reply}`);
    console.log();
  }

  console.log('='.repeat(60));
  console.log('Final memory state');
  console.log(`  topics  : ${JSON.stringify(topics)}`);
  console.log(`  intents : ${JSON.stringify(intents)}`);
  if (memorySummary) console.log(`  summary : ${memorySummary}`);
  console.log(`  history : ${history.length} turns in window`);
}

main().catch(console.error);
