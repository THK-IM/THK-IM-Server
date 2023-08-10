package model

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type (
	SessionUser struct {
		SessionId  int64 `gorm:"session_id" json:"session_id"`
		UserId     int64 `gorm:"user_id" json:"user_id"`
		Type       int   `gorm:"type" json:"type"`
		Status     int   `gorm:"status" json:"status"`
		CreateTime int64 `gorm:"create_time" json:"create_time"`
		UpdateTime int64 `gorm:"update_time" json:"update_time"`
		Deleted    int8  `gorm:"deleted" json:"deleted"`
	}

	SessionUserModel interface {
		FindSessionUserCount(sessionId int64) (int, error)
		FindUIdsInSessionWithoutStatus(sessionId int64, status int, uIds []int64) []int64
		FindUIdsInSessionContainStatus(sessionId int64, status int, uIds []int64) []int64
		AddUser(session *Session, entityId int64, uIds []int64, maxCount int) (err error)
		DelUser(session *Session, uIds []int64) (err error)
		UpdateUser(sId, uId int64, status int, tx *gorm.DB) (err error)
	}

	defaultSessionUserModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionUserModel) FindSessionUserCount(sessionId int64) (int, error) {
	count := 0
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ?  and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId).Scan(&count).Error
	return count, err
}

func (d defaultSessionUserModel) FindUIdsInSessionWithoutStatus(sessionId int64, status int, uIds []int64) []int64 {
	userIds := make([]int64, 0)
	uIdsCondition := ""
	if len(uIds) > 0 {
		uIdsCondition = " and user_id in ?"
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? = 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(uIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, uIds, status).Scan(&userIds)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&userIds)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	return userIds
}

func (d defaultSessionUserModel) FindUIdsInSessionContainStatus(sessionId int64, status int, uIds []int64) []int64 {
	userIds := make([]int64, 0)
	uIdsCondition := ""
	if len(uIds) > 0 {
		uIdsCondition = " and user_id in ?"
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? > 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(uIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, uIds, status).Scan(&userIds)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&userIds)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	return userIds
}

func (d defaultSessionUserModel) AddUser(session *Session, entityId int64, uIds []int64, maxCount int) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()
	count := 0
	tableName := d.genSessionUserTableName(session.Id)
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? and deleted = 0", tableName)
	if err = tx.Raw(sqlStr, session.Id).Scan(&count).Error; err != nil {
		return err
	}
	if count > maxCount-len(uIds) {
		return errorx.ErrGroupMemberCount
	}
	t := time.Now().UnixMilli()
	sql1 := "insert into " + d.genSessionUserTableName(session.Id) +
		" (session_id, user_id, type, create_time, update_time) values (?, ?, ?, ?, ?) " +
		"on duplicate key update status =?, deleted = ?, update_time = ? "
	for _, id := range uIds {
		if err = tx.Exec(sql1, session.Id, id, session.Type, t, t, 0, 0, t).Error; err != nil {
			return err
		}

		sql2 := "insert into " + d.genUserSessionTableName(id) +
			" (session_id, user_id, type, entity_id, create_time, update_time) values (?, ?, ?, ?, ?, ?) " +
			"on duplicate key update top = ?, status =?, deleted = ?, update_time = ? "
		if err = tx.Exec(sql2, session.Id, id, session.Type, entityId, t, t, 0, 0, 0, t).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d defaultSessionUserModel) DelUser(session *Session, uIds []int64) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()
	t := time.Now().UnixMilli()
	sql1 := "update " + d.genSessionUserTableName(session.Id) +
		" set deleted = ?, update_time = ? where session_id = ? and user_id = ?"
	for _, id := range uIds {
		if err = tx.Exec(sql1, 1, t, session.Id, id).Error; err != nil {
			return err
		}

		sql2 := "update " + d.genUserSessionTableName(id) +
			" set deleted = ?, update_time = ? where session_id = ? and user_id = ?"
		if err = tx.Exec(sql2, 1, t, session.Id, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d defaultSessionUserModel) UpdateUser(sId, uId int64, status int, tx *gorm.DB) (err error) {
	t := time.Now().UnixMilli()
	sql := "update " + d.genSessionUserTableName(sId) +
		" set status = ?, update_time = ? where session_id = ? and user_id = ?"
	if tx == nil {
		return d.db.Exec(sql, status, t, sId, uId).Error
	} else {
		return tx.Exec(sql, status, t, sId, uId).Error
	}
}

func (d defaultSessionUserModel) genUserSessionTableName(userId int64) string {
	return "user_session_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func (d defaultSessionUserModel) genSessionUserTableName(sessionId int64) string {
	return "session_user_" + fmt.Sprintf("%02d", sessionId%(d.shards))
}

func NewSessionUserModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionUserModel {
	return defaultSessionUserModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
