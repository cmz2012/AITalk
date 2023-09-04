package dal

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	InitClient()
	type args struct {
		ctx context.Context
		msg string
	}
	tests := []struct {
		name      string
		args      args
		wantReply string
		wantErr   bool
	}{
		// : Add test cases.
		{
			args: args{
				ctx: context.Background(),
				msg: "Who's Yao Ming",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReply, err := ChatCompletion(tt.args.ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatCompletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotReply)
		})
	}
}

func TestTranscribe(t *testing.T) {
	InitClient()
	f, err := os.Open("../test/Sports.wav")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	type args struct {
		ctx         context.Context
		audioReader io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		//: Add test cases.
		{
			args: args{
				ctx:         context.Background(),
				audioReader: f,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, err := Transcribe(tt.args.ctx, tt.args.audioReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transcribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotText)
		})
	}
}
