package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type (
	SessionMessage struct {
		MsgId      int64  `gorm:"msg_id" json:"msg_id"`
		ClientId   int64  `gorm:"client_id" json:"client_id"`
		SessionId  int64  `gorm:"session_id" json:"session_id"`
		FromUserId int64  `gorm:"from_user_id" json:"from_user_id"`
		AtUsers    string `gorm:"at_users" json:"at_users"`
		MsgType    int    `gorm:"msg_type" json:"msg_type"`
		MsgContent string `gorm:"msg_content" json:"msg_content"`
		ReplyMsgId int64  `gorm:"reply_msg_id" json:"reply_msg_id"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
		UpdateTime int64  `gorm:"update_time" json:"update_time"`
		Deleted    int8   `gorm:"deleted" json:"deleted"`
	}

	SessionMessageModel interface {
		DelMessages(sessionId int64, messageIds []int64, from, to int64) error
		SendMessage(clientId int64, fromUserId int64, sessionId int64, msgId int64, msgContent string,
			msgType int, atUserIds string, replayMsgId int64) (*SessionMessage, error)
		FindMessageByClientId(sessionId, clientId, fromUId int64) (*SessionMessage, error)
		GetSessionMessages(sid int64, ctime int, offset, count int, msgIds []int64) ([]*SessionMessage, error)
	}

	defaultSessionMessageModel struct {
		shards        int64
		db            *gorm.DB
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionMessageModel) DelMessages(sessionId int64, messageIds []int64, from, to int64) error {
	if len(messageIds) > 0 {
		sqlStr := fmt.Sprintf("update %s set deleted = 1 where session_id = ? and msg_id in ? and create_time >= ? and create_time <= ? ", d.genSessionMessageTableName(sessionId))
		err := d.db.Exec(sqlStr, sessionId, messageIds, from, to).Error
		return err
	} else {
		sqlStr := fmt.Sprintf("update %s set deleted = 1 where session_id = ? and create_time >= ? and create_time <= ?", d.genSessionMessageTableName(sessionId))
		err := d.db.Exec(sqlStr, sessionId, from, to).Error
		return err
	}
}

func (d defaultSessionMessageModel) SendMessage(clientId int64, fromUserId int64, sessionId int64, msgId int64, msgContent string, msgType int, atUserIds string, replayMsgId int64) (*SessionMessage, error) {
	currTime := time.Now().UnixMilli()
	sessionMessage := &SessionMessage{
		MsgId:      msgId,
		ClientId:   clientId,
		SessionId:  sessionId,
		FromUserId: fromUserId,
		AtUsers:    atUserIds,
		MsgType:    msgType,
		MsgContent: msgContent,
		ReplyMsgId: replayMsgId,
		CreateTime: currTime,
		UpdateTime: currTime,
		Deleted:    0,
	}
	tx := d.db.Table(d.genSessionMessageTableName(sessionId)).Create(sessionMessage)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sessionMessage, nil
}

func (d defaultSessionMessageModel) FindMessageByClientId(sessionId, clientId, fromUId int64) (*SessionMessage, error) {
	result := &SessionMessage{}
	strSql := "select * from " + d.genSessionMessageTableName(sessionId) + " where session_id = ? and client_id = ? and from_user_id = ?"
	err := d.db.Raw(strSql, sessionId, clientId, fromUId).Scan(result).Error
	return result, err
}

func (d defaultSessionMessageModel) GetSessionMessages(sid int64, ctime int, offset, count int, msgIds []int64) ([]*SessionMessage, error) {
	result := make([]*SessionMessage, 0)
	if len(msgIds) == 0 {
		strSql := "select * from " + d.genSessionMessageTableName(sid) + " where session_id = ? and deleted = 0 and from_user_id != 0 and create_time < ? order by create_time desc limit ?,?"
		err := d.db.Raw(strSql, sid, ctime, offset, count).Scan(&result).Error
		return result, err
	} else {
		strSql := "select * from " + d.genSessionMessageTableName(sid) + " where session_id = ? and msg_id in ? and create_time < ? order by create_time desc limit ?,?"
		err := d.db.Raw(strSql, sid, msgIds, ctime, offset, count).Scan(&result).Error
		return result, err
	}
}

func (d defaultSessionMessageModel) genSessionMessageTableName(sessionId int64) string {
	return "session_message_" + fmt.Sprintf("%02d", sessionId%(d.shards))
}

func NewSessionMessageModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionMessageModel {
	return defaultSessionMessageModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
