const mockData = {
  "cases": [
    {
      "description": "User asking about order status.",
      "id": "order-query",
      "output": {
        "confidence": 0.91,
        "intent": "query_order",
        "rankings": [
          {
            "confidence": 0.91,
            "intent": "query_order"
          },
          {
            "confidence": 0.05,
            "intent": "request_refund"
          },
          {
            "confidence": 0.04,
            "intent": "general_inquiry"
          }
        ]
      }
    }
  ],
  "delay": {
    "maxMs": 100,
    "minMs": 30
  },
  "random": {
    "enabled": false,
    "seed": "intent-recognition"
  },
  "version": "1.0.0"
};

export default mockData;
