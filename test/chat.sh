OPENAI_API_KEY=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

curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "user",
        "content": "介绍一下mysql的mvcc"
      }
    ]
  }'

