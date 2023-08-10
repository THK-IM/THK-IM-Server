package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	SingleSessionType     = 1
	GroupSessionType      = 2
	SuperGroupSessionType = 3
)

const (
	// MutedBitInSessionStatus 被禁言标志位
	MutedBitInSessionStatus = 1
)

type (
	Session struct {
		Id         int64  `gorm:"id" gorm:"primaryKey" json:"id"`
		Name       string `gorm:"name" json:"name"`
		Remark     string `gorm:"remark" json:"remark"`
		Type       int    `gorm:"type" json:"type"`
		Status     int    `gorm:"status" json:"status"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
		UpdateTime int64  `gorm:"update_time" json:"update_time"`
		Deleted    int8   `gorm:"deleted" json:"deleted"`
	}

	SessionModel interface {
		UpdateSession(id int64, status *int, name, remark *string) error
		FindSession(id int64, tx *gorm.DB) (*Session, error)
		CreateEmptySession(sessionType int, tx *gorm.DB) (*Session, error)
	}

	defaultSessionModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionModel) UpdateSession(id int64, status *int, name, remark *string) error {
	if status == nil && name == nil && remark == nil {
		return nil
	}
	updateMap := make(map[string]interface{})
	if status != nil {
		updateMap["status"] = *status
	}
	if name != nil {
		updateMap["name"] = *name
	}
	if remark != nil {
		updateMap["remark"] = *remark
	}
	updateMap["update_time"] = time.Now().UnixMilli()
	return d.db.Table(d.genSessionTableName(id)).Where("id = ?", id).Updates(updateMap).Error
}

func (d defaultSessionModel) FindSession(id int64, tx *gorm.DB) (*Session, error) {
	sqlStr := "select * from " + d.genSessionTableName(id) + " where id = ?"
	session := &Session{}
	var err error
	if tx != nil {
		tx = tx.Raw(sqlStr, id).Scan(session)
		err = tx.Error
	} else {
		err = d.db.Raw(sqlStr, id).Scan(session).Error
	}
	if err != nil {
		return nil, err
	}
	return session, err
}

func (d defaultSessionModel) CreateEmptySession(sessionType int, tx *gorm.DB) (*Session, error) {
	sessionId := int64(d.snowflakeNode.Generate())
	currTime := time.Now().UnixMilli()
	session := Session{
		Id:         sessionId,
		Type:       sessionType,
		CreateTime: currTime,
		UpdateTime: currTime,
	}
	var err error
	if tx != nil {
		tx = tx.Table(d.genSessionTableName(sessionId)).Create(&session)
		err = tx.Error
	} else {
		err = d.db.Table(d.genSessionTableName(sessionId)).Create(&session).Error
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (d defaultSessionModel) genSessionTableName(sessionId int64) string {
	return "session_" + fmt.Sprintf("%02d", sessionId%(d.shards))
}

func NewSessionModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionModel {
	return defaultSessionModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
