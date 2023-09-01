package model

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	// MutedAllBitInUserSessionStatus 全员被禁言标志位
	MutedAllBitInUserSessionStatus = 1 << 0
	// MutedSingleBitInUserSessionStatus 用户被禁言标志位
	MutedSingleBitInUserSessionStatus = 1 << 1
	// RejectBitInUserSessionStatus 拒收标志位
	RejectBitInUserSessionStatus = 1 << 0
	// SilenceBitInUserSessionStatus 静音标志位
	SilenceBitInUserSessionStatus = 1 << 1
)

type (
	UserSession struct {
		SessionId  int64  `gorm:"session_id" json:"session_id"`
		UserId     int64  `gorm:"user_id" json:"user_id"`
		Type       int    `gorm:"type" json:"type"`
		EntityId   int64  `gorm:"entity_id" json:"entity_id"`
		Name       string `gorm:"name" json:"name"`
		Remark     string `gorm:"remark" json:"remark"`
		Top        int64  `gorm:"top" json:"top"`
		Role       int    `gorm:"role" json:"role"`
		Mute       int    `gorm:"mute" json:"mute"`
		Status     int    `gorm:"status" json:"status"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
		UpdateTime int64  `gorm:"update_time" json:"update_time"`
		Deleted    int8   `gorm:"deleted" json:"deleted"`
	}

	UserSessionModel interface {
		RecoverUserSession(uId, sId, time int64) error
		FindUserSessionByEntityId(uId, entityId int64, sessionType int, containDeleted bool) (*UserSession, error)
		UpdateUserSession(uId []int64, sId int64, sessionName, sessionRemark, mute *string, top *int64, status, role *int, tx *gorm.DB) error
		FindEntityIdsInUserSession(uId int64, sId int64) []int64
		GetUserSessions(uId, mTime int64, offset, count int) ([]*UserSession, error)
		GetUserSession(uId, sId int64) (*UserSession, error)
	}

	defaultUserSessionModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserSessionModel) RecoverUserSession(userId, sessionId, time int64) error {
	conditions := []int64{userId, sessionId}
	updateMap := make(map[string]interface{})
	updateMap["top"] = 0
	updateMap["status"] = 0
	updateMap["create_time"] = time
	updateMap["update_time"] = time
	updateMap["deleted"] = 0
	return d.db.Table(d.genUserSessionTableName(userId)).Where("user_id = ? and session_id = ?", conditions).Updates(updateMap).Error

}

func (d defaultUserSessionModel) FindUserSessionByEntityId(uId, entityId int64, sessionType int, containDeleted bool) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.genUserSessionTableName(uId) + " where user_id = ? and entity_id = ? and type = ?"
	if !containDeleted {
		sqlStr += " and deleted = 0"
	}
	err := d.db.Raw(sqlStr, uId, entityId, sessionType).Scan(&userSession).Error
	return userSession, err
}

func (d defaultUserSessionModel) UpdateUserSession(uIds []int64, sId int64, sessionName, sessionRemark, mute *string, top *int64, status, role *int, tx *gorm.DB) error {
	// 分表uid数组
	sharedUIds := make(map[int64][]int64, 0)
	for _, uId := range uIds {
		share := uId % d.shards
		if sharedUIds[share] == nil {
			sharedUIds[share] = make([]int64, 0)
		}
		sharedUIds[share] = append(sharedUIds[share], uId)
	}
	for k, v := range sharedUIds {
		if sessionName == nil && sessionRemark == nil && top == nil && status == nil && mute == nil && role == nil {
			continue
		}
		sqlBuffer := bytes.Buffer{}
		sqlBuffer.WriteString(fmt.Sprintf("update %s set ", d.genUserSessionTableName(k)))
		if sessionName != nil {
			sqlBuffer.WriteString(fmt.Sprintf("name = '%s', ", *sessionName))
		}
		if sessionRemark != nil {
			sqlBuffer.WriteString(fmt.Sprintf("remark = '%s', ", *sessionRemark))
		}
		if top != nil {
			sqlBuffer.WriteString(fmt.Sprintf("top = %d, ", *top))
		}
		if status != nil {
			sqlBuffer.WriteString(fmt.Sprintf("status = %d, ", *status))
		}
		if mute != nil {
			sqlBuffer.WriteString(fmt.Sprintf("mute = %s, ", *mute))
		}
		if role != nil {
			sqlBuffer.WriteString(fmt.Sprintf("role = %d, ", *role))
		}
		sqlBuffer.WriteString(fmt.Sprintf("update_time = %d ", time.Now().UnixMilli()))
		sqlBuffer.WriteString("where session_id = ? and user_id in ? ")
		db := tx
		if db == nil {
			db = d.db
		}
		err := tx.Exec(sqlBuffer.String(), sId, v).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (d defaultUserSessionModel) FindEntityIdsInUserSession(uId int64, sId int64) []int64 {
	entityIds := make([]int64, 0)
	sqlStr := fmt.Sprintf("select entity_id from %s where user_id = ? and session_id = ? and deleted = 0", d.genUserSessionTableName(uId))
	_ = d.db.Raw(sqlStr, uId, sId).Scan(&entityIds).Error
	return entityIds
}

func (d defaultUserSessionModel) GetUserSessions(uId, mTime int64, offset, count int) ([]*UserSession, error) {
	userSessions := make([]*UserSession, 0)
	sqlStr := "select * from " + d.genUserSessionTableName(uId) + " where user_id = ? and update_time > ? limit ? offset ?"
	err := d.db.Raw(sqlStr, uId, mTime, count, offset).Scan(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (d defaultUserSessionModel) GetUserSession(uId, sId int64) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.genUserSessionTableName(uId) + " where user_id = ? and session_id = ?"
	err := d.db.Raw(sqlStr, uId, sId).Scan(userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (d defaultUserSessionModel) genUserSessionTableName(uId int64) string {
	return "user_session_" + fmt.Sprintf("%02d", uId%(d.shards))
}

func NewUserSessionModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserSessionModel {
	return defaultUserSessionModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
