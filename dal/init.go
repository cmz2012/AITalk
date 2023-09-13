package dal

import (
	"context"
	_model "github.com/cmz2012/AITalk/dal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	db *gorm.DB
)

func InitDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic(err)
	}
}

// GetSessionByUser 查询用户所有对话session
func GetSessionByUser(ctx context.Context, userID int64) (sessions []int64, err error) {
	sessions = make([]int64, 0)
	err = db.Model(&_model.Message{}).Distinct("session_id").Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("[GetSessionByUser]: %v", err)
		return
	}
	return
}

// GetMessageBySession 查询会话内的消息
func GetMessageBySession(ctx context.Context, sessionID int64) (msg []*_model.Message, err error) {
	msg = make([]*_model.Message, 0)
	err = db.Model(&_model.Message{}).Where("session_id = ?", sessionID).Find(&msg).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("[GetMessageBySession]: %v", err)
		return
	}
	return
}

func InsertMsg(ctx context.Context, msg *_model.Message) (err error) {
	msg.CreateTime = time.Now()
	msg.UpdateTime = msg.CreateTime
	err = db.Model(&_model.Message{}).Create(msg).Error
	return
}
