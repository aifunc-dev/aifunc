// =============================================================================
// Test instructions:
//   1. npm install
//   2. npm run build
//   3. npm run start
//
// This example demonstrates a customer support pipeline that chains multiple
// aifunc calls: sentiment analysis → intent recognition → structured extraction.
// It requires a real LLM to produce meaningful results. Update the config object
// with your API endpoint, model name, and API key to run the full pipeline.
// =============================================================================

import { analyzeSentiment, AIFuncConfig, AnalyzeSentimentInput } from './aifunc/analyze-sentiment';
import { recognizeIntent, RecognizeIntentInput } from './aifunc/recognize-intent';
import { extractJson, ExtractJsonInput } from './aifunc/extract-json';

// const config: AIFuncConfig = {
//   baseURL: 'https://your-api-endpoint/v1',
//   model: 'your-model-name',
//   apiKey: 'your-api-key',
// };

// To use a real model, replace the line below with the commented config above.
const config: AIFuncConfig = { mock: true };

// Tip: Each call accepts its own config — use cheaper models for simple tasks to save cost.
//
// const cheap: AIFuncConfig = { baseURL: '...', model: '...', apiKey: '...' };
// const strong: AIFuncConfig = { baseURL: '...', model: '...', apiKey: '...' };
//
// analyzeSentiment(cheap, ...)   // classification is simple, cheap model is fine
// extractJson(strong, ...)       // extraction needs accuracy, use a stronger model

if (config.mock) {
  console.log(
    'This example requires a real LLM to produce meaningful results.\n' +
    'Mock mode cannot simulate multi-step reasoning (sentiment → intent → extraction).\n' +
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

const MESSAGES = [
  "What the hell?! I ordered this a WEEK ago and it still hasn't shipped! I want my money back NOW!",
  "Hi, I'd like to check on my order #ORD-20240601-123. It's been three days with no shipping update.",
  "Your stupid app crashed again and I lost all my data! Fix it or I'm leaving!",
  "I was charged twice this month. Transaction IDs: TXN-88201 and TXN-88202. Please help.",
  "It would be cool if you added a dark mode. The bright screen hurts my eyes at night.",
  "How do I export my purchase history to CSV? I can't find the option.",
  "I am SO FURIOUS! Your delivery guy threw my package over the fence and it's destroyed! I want a manager NOW!",
  "Any ongoing promotions for loyal customers? I've been a member for 2 years.",
];

const INTENTS = ['query_order', 'request_refund', 'technical_support', 'billing_issue', 'feature_request', 'general_inquiry'];

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

async function main() {
  const message = MESSAGES[Math.floor(Math.random() * MESSAGES.length)];
  console.log(`Customer: ${message}\n`);

  // Step 1: Sentiment analysis
  const sentiment = await retry(
    () => analyzeSentiment(config, {
      text: message,
      labels: ['angry', 'frustrated', 'neutral', 'happy', 'other'],
    } as AnalyzeSentimentInput),
    "analyzeSentiment",
  );
  console.log(`Sentiment: ${sentiment.label} (${(sentiment.confidence * 100).toFixed(0)}%)`);

  if (sentiment.label === 'angry' && sentiment.confidence > 0.7) {
    console.log('\n=> call_human_agent(message, priority="HIGH")');
    return;
  }

  // Step 2: Intent recognition
  const intentResult = await retry(
    () => recognizeIntent(config, { text: message, intents: INTENTS } as RecognizeIntentInput),
    "recognizeIntent",
  );
  const intent = intentResult.intent;
  console.log(`Intent: ${intent} (${(intentResult.confidence * 100).toFixed(0)}%)`);

  // Step 3: Route by intent
  switch (intent) {
    case 'query_order': {
      const info = await retry(() => extractJson(config, { text: message, fields: [
        { name: 'order_id', description: 'Order number', type: 'string' },
        { name: 'issue', description: 'What the customer wants to know', type: 'string' },
      ] } as ExtractJsonInput), "extractJson");
      console.log(`\n=> query_order_system(order_id="${info.extracted['order_id']}", issue="${info.extracted['issue']}")`);
      break;
    }
    case 'request_refund': {
      const info = await retry(() => extractJson(config, { text: message, fields: [
        { name: 'order_id', description: 'Order number', type: 'string' },
        { name: 'reason', description: 'Reason for refund', type: 'string' },
      ] } as ExtractJsonInput), "extractJson");
      console.log(`\n=> submit_refund(order_id="${info.extracted['order_id']}", reason="${info.extracted['reason']}")`);
      break;
    }
    case 'technical_support': {
      const info = await retry(() => extractJson(config, { text: message, fields: [
        { name: 'issue', description: 'Technical problem', type: 'string' },
        { name: 'platform', description: 'Device or platform', type: 'string' },
      ] } as ExtractJsonInput), "extractJson");
      console.log(`\n=> create_tech_ticket(issue="${info.extracted['issue']}", platform="${info.extracted['platform']}")`);
      break;
    }
    case 'billing_issue': {
      const info = await retry(() => extractJson(config, { text: message, fields: [
        { name: 'transaction_id', description: 'Transaction ID', type: 'string' },
        { name: 'problem', description: 'Billing problem', type: 'string' },
      ] } as ExtractJsonInput), "extractJson");
      console.log(`\n=> flag_billing_dispute(transaction_id="${info.extracted['transaction_id']}", problem="${info.extracted['problem']}")`);
      break;
    }
    case 'feature_request': {
      const info = await retry(() => extractJson(config, { text: message, fields: [
        { name: 'feature', description: 'Requested feature', type: 'string' },
      ] } as ExtractJsonInput), "extractJson");
      console.log(`\n=> log_feature_request(feature="${info.extracted['feature']}")`);
      break;
    }
    default: {
      const info = await retry(() => extractJson(config, { text: message, fields: [
        { name: 'question', description: "Customer's question", type: 'string' },
      ] } as ExtractJsonInput), "extractJson");
      console.log(`\n=> send_to_faq_system(question="${info.extracted['question']}")`);
    }
  }
}

main().catch(console.error);
