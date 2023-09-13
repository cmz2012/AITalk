package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cmz2012/AITalk/dal"
	"github.com/cmz2012/AITalk/dal/model"
	_utils "github.com/cmz2012/AITalk/utils"
	"github.com/spf13/cast"
	"io/fs"
	"io/ioutil"
	"os"
)

func Transcribe(ctx context.Context, c *app.RequestContext) {
	sessionID := cast.ToInt64(string(c.FormValue("session_id")))
	userID := cast.ToInt64(string(c.FormValue("user_id")))

	if userID <= 0 {
		c.String(consts.StatusBadRequest, "userID must > 0")
		return
	}
	if sessionID <= 0 {
		uu, _ := _utils.GenIntUUID()
		sessionID = int64(uu)
	}
	audio, err := c.FormFile("audio")
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	f, err := audio.Open()
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	defer f.Close()
	text, err := dal.Transcribe(ctx, f)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	msg := &model.Message{
		SessionID: sessionID,
		UserID:    userID,
		Data:      text,
	}

	// 写文件
	name, _ := _utils.GenStrUUID()
	out := os.Getenv("CURDIR") + "/tmp/" + name + ".wav"
	f.Seek(0, 0)
	bs, _ := ioutil.ReadAll(f)
	err = ioutil.WriteFile(out, bs, fs.ModePerm)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	msg.AudioKey = name + ".wav"
	err = dal.InsertMsg(ctx, msg)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.JSON(consts.StatusOK, msg)
}

func Audio(ctx context.Context, c *app.RequestContext) {
	fileName := c.Param("file_name")
	out := os.Getenv("CURDIR") + "/tmp/" + fileName
	data, err := ioutil.ReadFile(out)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.Data(consts.StatusOK, "audio/wav", data)
}
