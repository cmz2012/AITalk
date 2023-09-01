package dal

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

const (
	model = "tts_models/en/ljspeech/tacotron2-DCA"
)

func Text2Speech(ctx context.Context, text string) (r io.ReadCloser, err error) {
	u, err := uuid.NewUUID()
	if err != nil {
		logrus.Errorf("[Text2Speech]: generate uuid %v", err)
		return
	}
	dir := os.Getenv("CURDIR") + "/tmp/"
	out := u.String() + ".wav"
	cmd := exec.CommandContext(ctx, "tts", "--text", "\""+text+"\"", "--model_name", model, "--out_path", dir+out)
	logrus.Infof("[Text2Speech]: cmd = %v", cmd.String())
	b, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("[Text2Speech]: tts run %v, %v", err, string(b))
		return
	}
	r, err = os.Open(dir + out)
	if err != nil {
		logrus.Errorf("[Text2Speech]: open file %v", err)
		return
	}
	return
}
