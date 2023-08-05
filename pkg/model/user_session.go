package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	// MutedBitInUserSessionStatus 被禁言标志位
	MutedBitInUserSessionStatus = 1
	// RejectBitInUserSessionStatus 拒收标志位
	RejectBitInUserSessionStatus = 2
	// SilenceBitInUserSessionStatus 静音标志位
	SilenceBitInUserSessionStatus = 4
)

type (
	UserSession struct {
		SessionId  int64 `gorm:"session_id" json:"session_id"`
		UserId     int64 `gorm:"user_id" json:"user_id"`
		Type       int   `gorm:"type" json:"type"`
		EntityId   int64 `gorm:"entity_id" json:"entity_id"`
		Top        int64 `gorm:"top" json:"top"`
		Status     int   `gorm:"status" json:"status"`
		CreateTime int64 `gorm:"create_time" json:"create_time"`
		UpdateTime int64 `gorm:"update_time" json:"update_time"`
		Deleted    int8  `gorm:"deleted" json:"deleted"`
	}

	UserSessionModel interface {
		FindUserSessionByEntityId(uId, entityId int64, sessionType int) (*UserSession, error)
		UpdateUserSession(userId, sessionId int64, top *int64, status *int, tx *gorm.DB) error
		FindEntityIdsInUserSession(userId int64, sessionId int64) []int64
		GetUserSessions(uid, mTime int64, offset, count int) ([]*UserSession, error)
		GetUserSession(uid, sid int64) (*UserSession, error)
	}

	defaultUserSessionModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserSessionModel) FindUserSessionByEntityId(uId, entityId int64, sessionType int) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.genUserSessionTableName(uId) + " where user_id = ? and entity_id = ? and type = ?"
	err := d.db.Raw(sqlStr, uId, entityId, sessionType).Scan(&userSession).Error
	return userSession, err
}

func (d defaultUserSessionModel) UpdateUserSession(userId, sessionId int64, top *int64, status *int, tx *gorm.DB) error {
	if top == nil && status == nil {
		return nil
	}
	conditions := []int64{userId, sessionId}
	updateMap := make(map[string]interface{})
	if top != nil {
		updateMap["top"] = *top
	}
	if status != nil {
		updateMap["status"] = *status
	}
	updateMap["update_time"] = time.Now().UnixMilli()
	if tx != nil {
		return tx.Table(d.genUserSessionTableName(userId)).Where("user_id = ? and session_id = ?", conditions).Updates(updateMap).Error
	} else {
		return d.db.Table(d.genUserSessionTableName(userId)).Where("user_id = ? and session_id = ?", conditions).Updates(updateMap).Error
	}
}

func (d defaultUserSessionModel) FindEntityIdsInUserSession(userId int64, sessionId int64) []int64 {
	entityIds := make([]int64, 0)
	sqlStr := fmt.Sprintf("select entity_id from %s where user_id = ? and session_id = ? and deleted = 0", d.genUserSessionTableName(userId))
	_ = d.db.Raw(sqlStr, userId, sessionId).Scan(&entityIds).Error
	return entityIds
}

func (d defaultUserSessionModel) GetUserSessions(uid, mTime int64, offset, count int) ([]*UserSession, error) {
	userSessions := make([]*UserSession, 0)
	sqlStr := "select * from " + d.genUserSessionTableName(uid) + " where user_id = ? and update_time > ? limit ? offset ?"
	err := d.db.Raw(sqlStr, uid, mTime, count, offset).Scan(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (d defaultUserSessionModel) GetUserSession(uid, sid int64) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.genUserSessionTableName(uid) + " where user_id = ? and session_id = ?"
	err := d.db.Raw(sqlStr, uid, sid).Scan(userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (d defaultUserSessionModel) genUserSessionTableName(userId int64) string {
	return "user_session_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func NewUserSessionModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserSessionModel {
	return defaultUserSessionModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
