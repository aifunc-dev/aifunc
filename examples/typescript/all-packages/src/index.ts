import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';
import { translate, TranslateInput } from './aifunc/translate';
import { analyzeSentiment, AnalyzeSentimentInput } from './aifunc/analyze-sentiment';
import { detectLanguage, DetectLanguageInput } from './aifunc/detect-language';
import { rewrite, RewriteInput } from './aifunc/rewrite';
import { extractKeywords, ExtractKeywordsInput } from './aifunc/extract-keywords';
import { classify, ClassifyInput } from './aifunc/classify';
import { recognizeIntent, RecognizeIntentInput } from './aifunc/recognize-intent';
import { extractEntities, ExtractEntitiesInput } from './aifunc/extract-entities';
import { extractJson, ExtractJsonInput } from './aifunc/extract-json';
import { generateSlug, GenerateSlugInput } from './aifunc/generate-slug';
import { generateReply, GenerateReplyInput } from './aifunc/generate-reply';
import { generatePost, GeneratePostInput } from './aifunc/generate-post';
import { generateEmail, GenerateEmailInput } from './aifunc/generate-email';
import { generateTitle, GenerateTitleInput } from './aifunc/generate-title';
import { answerQuestion, AnswerQuestionInput } from './aifunc/answer-question';
import { scoreQuality, ScoreQualityInput } from './aifunc/score-quality';

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

function section(title: string) {
  console.log(`\n${'='.repeat(60)}`);
  console.log(`  ${title}`);
  console.log(`${'='.repeat(60)}`);
}

async function main() {
  const ARTICLE =
    'In 1915, Albert Einstein published the General Theory of Relativity, ' +
    'fundamentally transforming our understanding of physics. The theory posits ' +
    'that gravity is not an invisible force, but rather a curvature of spacetime ' +
    'caused by the presence of mass and energy. This groundbreaking framework ' +
    'revolutionized modern science and introduced the famous equation E=mc².';

  section('1. DETECT LANGUAGE');
  const langSamples = [
    'The quick brown fox jumps over the lazy dog.',
    'Der schnelle braune Fuchs springt über den faulen Hund.',
    'Le renard brun rapide saute par-dessus le chien paresseux.',
    'El veloz zorro marrón salta sobre el perro perezoso.',
  ];
  for (const text of langSamples) {
    try {
      const r = await detectLanguage(config, { text } as DetectLanguageInput);
      console.log(`  [${r.language}] ${r.languageName} (conf: ${(r.confidence * 100).toFixed(0)}%)  "${text.slice(0, 40)}"`);
    } catch (e: any) { console.log(`  Error: ${e.message}`); }
  }

  section('2. GENERATE SLUG');
  try {
    const r = await generateSlug(config, { title: '10 Practical Tips for Writing Faster Python Code', language: 'en' } as GenerateSlugInput);
    console.log(`Title : 10 Practical Tips for Writing Faster Python Code`);
    console.log(`Slug  : ${r.slug}`);
    console.log(`Meta  : ${r.metaDescription}`);
    console.log(`Tags  : ${r.tags}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('3. SUMMARIZE');
  try {
    const r = await summarize(config, { text: ARTICLE, maxLength: 30 } as SummarizeInput);
    console.log(`Summary   : ${r.summary}`);
    console.log(`Word count: ${r.wordCount}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('4. TRANSLATE');
  try {
    const r = await translate(config, { text: 'The meeting has been moved to Friday at 3 PM.', targetLang: 'es' } as TranslateInput);
    console.log(`Original : The meeting has been moved to Friday at 3 PM.`);
    console.log(`Spanish  : ${r.translation}`);
    console.log(`Detected : ${r.sourceLang}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('5. REWRITE');
  try {
    const original = 'hey, just wanna let u know the deploy went fine, no issues at all';
    const r = await rewrite(config, { text: original, style: 'formal' } as RewriteInput);
    console.log(`Casual : ${original}`);
    console.log(`Formal : ${r.rewritten}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('6. GENERATE TITLE');
  try {
    const content = 'This guide covers how to use Docker and GitHub Actions to automate testing and deployment of a Node.js application to a cloud server.';
    const r = await generateTitle(config, { content, style: 'seo', count: 4 } as GenerateTitleInput);
    console.log(`Content: ${content}`);
    console.log('Titles:');
    for (let i = 0; i < r.titles.length; i++) console.log(`  ${i + 1}. ${r.titles[i]}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('7. EXTRACT KEYWORDS');
  try {
    const r = await extractKeywords(config, { text: ARTICLE, maxKeywords: 5 } as ExtractKeywordsInput);
    console.log('Keywords from article:');
    for (const kw of r.keywords) console.log(`  ${kw.word.padEnd(30)} relevance: ${kw.relevance}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('8. ANALYZE SENTIMENT');
  const sentSamples = [
    'The product arrived on time and works perfectly. Very happy!',
    'Terrible experience. The package was damaged and support ignored my emails.',
    'Item received. Does what it says.',
  ];
  for (const text of sentSamples) {
    try {
      const r = await analyzeSentiment(config, { text, labels: ['positive', 'negative', 'neutral'] } as AnalyzeSentimentInput);
      console.log(`  [${r.label.padEnd(8)} ${(r.confidence * 100).toFixed(0)}%] ${text.slice(0, 55)}`);
    } catch (e: any) { console.log(`  Error: ${e.message}`); }
  }

  section('9. CLASSIFY');
  const tickets = [
    "My order hasn't shipped after five days. Please help.",
    'The API returns a 500 error when the payload exceeds 1 MB.',
    'It would be great to have a dark mode option.',
    'I was charged twice for the same subscription this month.',
  ];
  const categories = ['shipping', 'technical', 'feature request', 'billing', 'other'];
  for (const ticket of tickets) {
    try {
      const r = await classify(config, { text: ticket, categories } as ClassifyInput);
      const top = r.classifications[0];
      console.log(`  [${top.category.padEnd(16)} ${(top.confidence * 100).toFixed(0)}%]  ${ticket.slice(0, 55)}`);
    } catch (e: any) { console.log(`  Error: ${e.message}`); }
  }

  section('10. RECOGNIZE INTENT');
  const intentMsgs = [
    'Where is my order? I placed it three days ago.',
    'I want a refund for the broken item.',
    'Can you tell me your business hours?',
    "I'd like to upgrade my subscription to the pro plan.",
  ];
  const intents = ['query_order', 'request_refund', 'general_inquiry', 'manage_subscription'];
  for (const msg of intentMsgs) {
    try {
      const r = await recognizeIntent(config, { text: msg, intents, context: 'You are a customer support routing system.' } as RecognizeIntentInput);
      console.log(`  [${r.intent.padEnd(20)} ${(r.confidence * 100).toFixed(0)}%]  "${msg.slice(0, 50)}"`);
    } catch (e: any) { console.log(`  Error: ${e.message}`); }
  }

  section('11. EXTRACT ENTITIES');
  try {
    const text = 'On March 10, 2024, NASA astronaut Sarah Mitchell landed at Kennedy Space Center in Florida.';
    const r = await extractEntities(config, { text, entityTypes: ['person', 'organization', 'location', 'date'] } as ExtractEntitiesInput);
    console.log(`Text: ${text}`);
    console.log('Entities:');
    for (const ent of r.entities) console.log(`  [${ent.type.padEnd(12)}] "${ent.text}"`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('12. EXTRACT JSON');
  try {
    const jobPost = 'We are looking for a Senior Backend Engineer in Berlin. Requirements: 5+ years experience, Go or Rust, Kubernetes.';
    const r = await extractJson(config, {
      text: jobPost,
      fields: [
        { name: 'title', description: 'Job title', type: 'string' },
        { name: 'location', description: 'City or country', type: 'string' },
        { name: 'skills', description: 'Required technical skills', type: 'array' },
        { name: 'experience_years', description: 'Minimum years of experience', type: 'number' },
      ],
    } as ExtractJsonInput);
    console.log(`Text     : ${jobPost}`);
    console.log(`Extracted: ${JSON.stringify(r.extracted)}`);
    console.log(`Missing  : ${JSON.stringify(r.missing)}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('13. ANSWER QUESTION');
  const qaCtx = 'AIFunc is a function-based AI toolkit. The CLI generates type-safe wrappers for Python, TypeScript, or Go.';
  const qaPairs: [string, string | undefined][] = [
    ['Which languages does AIFunc support?', qaCtx],
    ['What is a monad in functional programming?', undefined],
  ];
  for (const [q, ctx] of qaPairs) {
    try {
      const input: any = { question: q, maxLength: 60 };
      if (ctx) input.context = ctx;
      const r = await answerQuestion(config, input as AnswerQuestionInput);
      const source = r.grounded ? 'from context' : 'general knowledge';
      console.log(`  Q: ${q}`);
      console.log(`  A: ${r.answer}  [${source}, conf: ${(r.confidence * 100).toFixed(0)}%]\n`);
    } catch (e: any) { console.log(`  Error: ${e.message}\n`); }
  }

  section('14. GENERATE REPLY');
  try {
    const message = "I placed an order three days ago but haven't received a shipping confirmation yet.";
    const r = await generateReply(config, { message, tone: 'empathetic', context: 'You are a customer support agent.' } as GenerateReplyInput);
    console.log(`Customer : ${message}`);
    console.log(`Reply    : ${r.reply}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('15. GENERATE POST');
  try {
    const r = await generatePost(config, { topic: 'How async Python cut our API response time by 60%', platform: 'linkedin', tone: 'professional', includeHashtags: true } as GeneratePostInput);
    console.log(`Post     : ${r.post}`);
    console.log(`Hashtags : ${r.hashtags.map(t => '#' + t)}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('16. GENERATE EMAIL');
  try {
    const r = await generateEmail(config, {
      intent: 'Apologize to a customer for a billing error',
      tone: 'formal',
      senderName: 'Billing Support Team',
      recipientName: 'Alex',
      keyPoints: ['Incorrect charge of $29.99 on June 1st', 'Charge has been fully refunded', '20% discount applied to next invoice'],
      language: 'English',
    } as GenerateEmailInput);
    console.log(`Subject: ${r.subject}`);
    console.log(`Body:\n${r.body}`);
  } catch (e: any) { console.log(`Error: ${e.message}`); }

  section('17. SCORE QUALITY');
  const qualitySamples: [string, string, string][] = [
    ['Our product is good. It has many features. Users like it.', 'customers', 'marketing'],
    ['To set up CI: 1) Install Docker. 2) Create deploy.yml. 3) Push to main.', 'developers', 'explanation'],
  ];
  for (const [text, audience, purpose] of qualitySamples) {
    try {
      const r = await scoreQuality(config, { text, targetAudience: audience, purpose, maxSuggestions: 3, strictness: 3 } as ScoreQualityInput);
      console.log(`Text       : ${text.slice(0, 55)}...`);
      console.log(`Score      : ${r.overallScore}/100  [${r.level}]`);
      console.log(`Summary    : ${r.summary}`);
      console.log(`Suggestions:`);
      for (const s of r.suggestions) console.log(`  - ${s}`);
      console.log();
    } catch (e: any) { console.log(`Error: ${e.message}\n`); }
  }

  if (config.mock) {
  console.log(
    'Notice: You are using mock mode for offline testing. ' +
    'Configure a real model for the full experience.'
  );
}
}

main().catch(console.error);
