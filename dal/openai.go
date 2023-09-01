package dal

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"io"
)

const (
	ApiKey = "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"

	SystemMessage = "I want you to act as a spoken English teacher and improver. I want you to keep your reply neat, limiting the reply to 100 words. I want you to strictly correct my grammar mistakes, typos, and factual errors."
)

var (
	client *openai.Client
)

func InitClient() {
	client = openai.NewClient(ApiKey)
}

func Transcribe(ctx context.Context, audioReader io.Reader) (text string, err error) {
	rsp, err := client.CreateTranscription(ctx, openai.AudioRequest{
		Model:  openai.Whisper1,
		Reader: audioReader,
		Format: openai.AudioResponseFormatJSON,
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
		MaxTokens: 200,
	})
	if err != nil {
		logrus.Errorf("[ChatCompletion]: %v", err)
		return
	}
	reply = rsp.Choices[0].Message.Content
	return
}
