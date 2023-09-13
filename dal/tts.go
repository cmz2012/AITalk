package dal

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

const (
	model = "tts_models/en/ljspeech/tacotron2-DCA"
)

func Text2Speech(ctx context.Context, text string) (out string, err error) {
	u, err := uuid.NewUUID()
	if err != nil {
		logrus.Errorf("[Text2Speech]: generate uuid %v", err)
		return
	}
	dir := os.Getenv("CURDIR") + "/tmp/"
	out = dir + u.String() + ".wav"
	cmd := exec.Command("tts", "--text", "\""+text+"\"", "--model_name", model, "--out_path", out)
	logrus.Infof("[Text2Speech]: cmd = %v", cmd.String())
	b, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("[Text2Speech]: tts run %v, %v", err, string(b))
		return
	}
	return
}
