package service

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cmz2012/AITalk/biz/model/chat"
	"github.com/cmz2012/AITalk/dal"
	"github.com/cmz2012/AITalk/dal/model"
	"github.com/hertz-contrib/websocket"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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

		// 插入db
		msg := &model.Message{
			SessionID: cc.req.SessionID,
			UserID:    cc.req.UserID,
			Data:      text,
		}
		err = dal.InsertMsg(nil, msg)
		if err != nil {
			logrus.Errorf("[ChatHandler]: %v", err)
			break
		}

		// write transcribe
		conn.WriteJSON(msg)

		reply, err := dal.ChatCompletion(ctx, text)
		if err != nil {
			logrus.Errorf("[ChatHandler]: ChatCompletion %v", err)
			break
		}
		logrus.Infof("[ChatHandler]: gpt reply : %v", reply)

		// write reply to db
		msg.UserID = 0 // bot
		msg.ID = 0
		msg.Data = reply
		err = dal.InsertMsg(nil, msg)
		if err != nil {
			logrus.Errorf("[ChatHandler]: %v", err)
			break
		}

		conn.WriteJSON(msg)

		out, err := dal.Text2Speech(ctx, reply)
		if err != nil {
			logrus.Errorf("[ChatHandler]: Text2Speech %v", err)
			break
		}
		replyAudio, err := ioutil.ReadFile(out)
		if err != nil {
			logrus.Errorf("[ChatHandler]: ReadFile %v", err)
			break
		}

		// write reply audio
		err = conn.WriteMessage(websocket.BinaryMessage, replyAudio)
		if err != nil {
			logrus.Errorf("[ChatHandler]: write websocket %v", err)
			break
		}
	}
}
