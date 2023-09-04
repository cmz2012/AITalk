package dal

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

const (
	SystemMessage = "You are a spoken English teacher and improver. Your reply should be limited to 100 words."
)

var (
	client *openai.Client
)

func InitClient() {
	ApiKey := os.Getenv("OPENAIAPIKEY")
	if ApiKey == "" {
		panic("$OPENAIAPIKEY not setup")
	}
	client = openai.NewClient(ApiKey)
}

func Transcribe(ctx context.Context, audioReader io.Reader) (text string, err error) {
	rsp, err := client.CreateTranscription(ctx, openai.AudioRequest{
		Model:    openai.Whisper1,
		Reader:   audioReader,
		FilePath: "example.wav",
		Format:   openai.AudioResponseFormatJSON,
	})
	if err != nil {
		logrus.Errorf("[CreateTranscription]: %v", err)
		return
	}
	text = rsp.Text
	return
}

func ChatCompletion(ctx context.Context, msg string) (reply string, err error) {
	rsp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: SystemMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			},
		},
		MaxTokens:   200,
		Temperature: 1.0,
	})
	if err != nil {
		logrus.Errorf("[ChatCompletion]: %v", err)
		return
	}
	reply = rsp.Choices[0].Message.Content
	return
}
