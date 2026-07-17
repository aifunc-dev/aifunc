import { chatStream, AIFuncConfig, ChatStreamInput } from './aifunc/chat-stream';
import { answerStream, AnswerStreamInput } from './aifunc/answer-stream';
import { explainStream, ExplainStreamInput } from './aifunc/explain-stream';
import { articleStream, ArticleStreamInput } from './aifunc/article-stream';
import { writeStream, WriteStreamInput } from './aifunc/write-stream';
import { translateStream, TranslateStreamInput } from './aifunc/translate-stream';
import { reviewStream, ReviewStreamInput } from './aifunc/review-stream';

// const config: AIFuncConfig = {
//   baseURL: 'https://your-api-endpoint/v1',
//   model: 'your-model-name',
//   apiKey: 'your-api-key',
//   maxRetries: 3,
// };

// To use a real model, replace the line below with the commented config above.
const config: AIFuncConfig = { mock: true };

if (config.mock) {
  console.log(
    'Notice: You are using mock mode for offline testing. ' +
    'Configure a real model for the full experience. Continuing with mock responses...'
  );
}

function section(title: string) {
  console.log(`\n${'='.repeat(60)}`);
  console.log(`  ${title}`);
  console.log(`${'='.repeat(60)}`);
}

async function streamPrint(tokens: AsyncIterable<string>) {
  for await (const token of tokens) {
    process.stdout.write(token);
  }
  process.stdout.write('\n');
}

async function main() {
  const ARTICLE =
    'In 1915, Albert Einstein published the General Theory of Relativity, ' +
    'fundamentally transforming our understanding of physics. The theory posits ' +
    'that gravity is not an invisible force, but rather a curvature of spacetime ' +
    'caused by the presence of mass and energy. This groundbreaking framework ' +
    'revolutionized modern science and introduced the famous equation E=mc².';

  const CODE_SNIPPET = `def fetch_user(user_id):
    conn = get_connection()
    result = conn.execute(f"SELECT * FROM users WHERE id = {user_id}")
    return result.fetchone()
`;

  // ─── Conversational & Q&A ─────────────────────────────────────────

  section('1. CHAT STREAM');
  console.log('User: Explain async/await in TypeScript in 3 sentences.\n');
  process.stdout.write('Assistant: ');
  await streamPrint(chatStream(config, {
    messages: [
      { role: 'user', content: 'Explain async/await in TypeScript in 3 sentences.' },
    ],
  } as ChatStreamInput));

  section('2. ANSWER STREAM (with context / RAG)');
  const context =
    'AIFunc is a function-based AI toolkit. Developers declare the packages they need ' +
    'in aifunc.json. The CLI generates type-safe wrappers for Python, TypeScript, or Go. ' +
    'Each package supports a mock mode for testing without consuming API credits. ' +
    'Streaming packages return tokens incrementally via AsyncIterable.';
  const question = 'How does AIFunc support offline testing, and what do streaming packages return?';
  console.log(`Q: ${question}\n`);
  process.stdout.write('A: ');
  await streamPrint(answerStream(config, {
    question,
    context,
    depth: 'concise',
    audience: 'technical',
  } as AnswerStreamInput));

  section('3. EXPLAIN STREAM');
  console.log('Topic: the event loop in Node.js\n');
  await streamPrint(explainStream(config, {
    topic: 'the event loop in Node.js',
    audience: 'intermediate',
    depth: 'standard',
  } as ExplainStreamInput));

  // ─── Long-form writing ────────────────────────────────────────────

  section('4. ARTICLE STREAM');
  const title = 'Why Typed AI Functions Beat Ad-Hoc Prompt Scripts';
  const outline =
    '- The cost of untyped prompt glue code\n' +
    '- How function-shaped AI APIs improve testability\n' +
    '- Streaming vs batch for product UX\n' +
    '- Practical adoption tips';
  console.log(`Title  : ${title}`);
  console.log(`Outline: ${outline}\n`);
  await streamPrint(articleStream(config, {
    title,
    outline,
    style: 'informational',
    audience: 'developers',
    wordCount: 250,
  } as ArticleStreamInput));

  section('5. WRITE STREAM');
  const prompt =
    'Write a short internal proposal recommending that our team adopt AIFunc ' +
    'for customer-support reply generation.';
  const structure =
    '1. Problem\n' +
    '2. Proposed approach\n' +
    '3. Expected benefits\n' +
    '4. Next steps';
  console.log(`Prompt   : ${prompt}`);
  console.log(`Structure: ${structure}\n`);
  await streamPrint(writeStream(config, {
    prompt,
    format: 'proposal',
    structure,
    tone: 'professional',
    audience: 'engineers',
    wordCount: 300,
  } as WriteStreamInput));

  // ─── Translation & review ─────────────────────────────────────────

  section('6. TRANSLATE STREAM');
  console.log(`Original (EN):\n${ARTICLE}\n`);
  console.log('Translation (zh-CN):\n');
  await streamPrint(translateStream(config, {
    text: ARTICLE,
    targetLang: 'zh-CN',
    style: 'natural',
    domain: 'technical',
  } as TranslateStreamInput));

  section('7. REVIEW STREAM');
  console.log(`Code under review:\n${CODE_SNIPPET}`);
  console.log('Findings:\n');
  await streamPrint(reviewStream(config, {
    content: CODE_SNIPPET,
    type: 'code',
    language: 'Python',
    focus: 'correctness, security',
    context: 'Simple data-access helper in a web API.',
    severity: 'all',
    outputLanguage: 'English',
  } as ReviewStreamInput));

  if (config.mock) {
    console.log(
      'Notice: You are using mock mode for offline testing. ' +
      'Configure a real model for the full experience.'
    );
  }
}

main().catch(console.error);
