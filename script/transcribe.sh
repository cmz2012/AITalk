OPENAI_API_KEY=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

curl https://api.openai.com/v1/audio/transcriptions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: multipart/form-data" \
  -F file="@/Users/bytedance/go/src/github.com/cmz2012/AITalk/script/Sports.wav" \
  -F model="whisper-1"
