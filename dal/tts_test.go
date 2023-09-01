package dal

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestText2Speech(t *testing.T) {
	os.Setenv("CURDIR", "..")
	type args struct {
		ctx  context.Context
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantR   io.ReadCloser
		wantErr bool
	}{
		// : Add test cases.
		{
			name: "",
			args: args{
				ctx:  context.Background(),
				text: "What sampling temperature to use",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := Text2Speech(tt.args.ctx, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Text2Speech() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotR.Close()
		})
	}
}
