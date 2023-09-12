package service

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cmz2012/AITalk/biz/model/chat"
	"github.com/cmz2012/AITalk/dal"
	"github.com/cmz2012/AITalk/dal/model"
	"github.com/cmz2012/AITalk/utils"
	"github.com/hertz-contrib/websocket"
	"github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

var upgrader = websocket.HertzUpgrader{CheckOrigin: func(ctx *app.RequestContext) bool { return true }}

func ChatUpgrade(ctx context.Context, c *app.RequestContext, req *chat.CreateChatReq) {
	cc := ChatContext{
		c:   c,
		req: req,
	}
	err := upgrader.Upgrade(c, cc.ChatHandler)
	if err != nil {
		logrus.Errorf("[ChatUpgrade]: %v", err)
		return
	}
}

type ChatContext struct {
	c   *app.RequestContext
	req *chat.CreateChatReq
}

func (cc ChatContext) WriteFileAndDB(text string, audio []byte, user int64) (msg *model.Message, err error) {
	// 写文件
	name, _ := utils.GenStrUUID()
	out := os.Getenv("CURDIR") + "/tmp/" + name + ".wav"
	err = ioutil.WriteFile(out, audio, fs.ModePerm)
	if err != nil {
		return
	}

	// 插入db
	msg = &model.Message{
		SessionID: cc.req.SessionID,
		UserID:    user,
		Data:      text,
		AudioKey:  name + ".wav",
	}
	err = dal.InsertMsg(nil, msg)
	if err != nil {
		return
	}
	return
}

func (cc ChatContext) ChatHandler(conn *websocket.Conn) {
	// 循环读取wav bytes
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Infof("[ChatHandler]: %v", err)
			break
		}
		logrus.Printf("[ChatHandler]: wav file length : %v", len(message))

		// 处理chat
		ctx := context.Background()
		text, err := dal.Transcribe(ctx, bytes.NewReader(message))
		if err != nil {
			logrus.Errorf("[ChatHandler]: Transcribe %v", err)
			break
		}
		logrus.Infof("[ChatHandler]: wav -> text : %v", text)

		msg, err := cc.WriteFileAndDB(text, message, cc.req.UserID)
		if err != nil {
			logrus.Errorf("[ChatHandler]: Transcribe %v", err)
			break
		}

		// write transcribe msg
		conn.WriteJSON(msg)

		reply, err := dal.ChatCompletion(ctx, text)
		if err != nil {
			logrus.Errorf("[ChatHandler]: ChatCompletion %v", err)
			break
		}
		logrus.Infof("[ChatHandler]: gpt reply : %v", reply)

		out, err := dal.Text2Speech(ctx, reply)
		if err != nil {
			logrus.Errorf("[ChatHandler]: Text2Speech %v", err)
			break
		}

		// write reply to db
		dirs := strings.Split(out, "/")

		msg.UserID = 0 // bot
		msg.ID = 0
		msg.Data = reply
		msg.AudioKey = dirs[len(dirs)-1]

		err = dal.InsertMsg(nil, msg)
		if err != nil {
			logrus.Errorf("[ChatHandler]: %v", err)
			break
		}

		// write reply msg
		conn.WriteJSON(msg)
	}
}
