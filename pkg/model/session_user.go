package model

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	SessionMember     = 1 // 普通成员 可以查询session信息，session会话历史消息
	SessionAdmin      = 2 // 管理员 模式下可以修改session基本信息，禁言单个用户, 添加/删除普通成员
	SessionSuperAdmin = 3 // 超级管理员 可以全员禁言, 添加/删除管理员
	SessionOwner      = 4 // 拥有者 可以添加超级管理员, 删除管理员，删除session
)

type (
	SessionUser struct {
		SessionId  int64 `gorm:"session_id" json:"session_id"`
		UserId     int64 `gorm:"user_id" json:"user_id"`
		Type       int   `gorm:"type" json:"type"`
		Role       int   `gorm:"role" json:"role"`
		Mute       int   `gorm:"mute" json:"mute"`
		Status     int   `gorm:"status" json:"status"`
		CreateTime int64 `gorm:"create_time" json:"create_time"`
		UpdateTime int64 `gorm:"update_time" json:"update_time"`
		Deleted    int8  `gorm:"deleted" json:"deleted"`
	}

	SessionUserModel interface {
		FindAllSessionUsers(sId int64) ([]*SessionUser, error)
		FindSessionUsers(sId int64, uIds []int64) ([]*SessionUser, error)
		FindSessionUser(sId, uId int64) (*SessionUser, error)
		FindSessionUserCount(sId int64) (int, error)
		FindUIdsInSessionWithoutStatus(sId int64, status int, uIds []int64) []int64
		FindUIdsInSessionContainStatus(sId int64, status int, uIds []int64) []int64
		AddUser(session *Session, entityIds []int64, uIds []int64, role []int, maxCount int) (err error)
		DelUser(session *Session, uIds []int64) (err error)
		UpdateUser(sId int64, uIds []int64, role, status *int, mute *string, tx *gorm.DB) (err error)
	}

	defaultSessionUserModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionUserModel) FindAllSessionUsers(sId int64) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?", tableName)
	err := d.db.Raw(sqlStr, sId).Scan(sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUsers(sId int64, uIds []int64) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ? and user_id in ? and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sId, uIds).Scan(sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUser(sId, uId int64) (*SessionUser, error) {
	sessionUser := &SessionUser{}
	tableName := d.genSessionUserTableName(sId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?  and user_id = ? and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sId, uId).Scan(sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUserCount(sessionId int64) (int, error) {
	count := 0
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ?  and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId).Scan(&count).Error
	return count, err
}

func (d defaultSessionUserModel) FindUIdsInSessionWithoutStatus(sessionId int64, status int, uIds []int64) []int64 {
	sessionUsers := make([]*SessionUser, 0)
	uIdsCondition := ""
	if len(uIds) > 0 {
		uIdsCondition = " and user_id in ? "
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? = 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(uIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, uIds, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	userIds := make([]int64, 0)
	for _, su := range sessionUsers {
		userIds = append(userIds, su.UserId)
	}
	return userIds
}

func (d defaultSessionUserModel) FindUIdsInSessionContainStatus(sessionId int64, status int, uIds []int64) []int64 {
	sessionUsers := make([]*SessionUser, 0)
	uIdsCondition := ""
	if len(uIds) > 0 {
		uIdsCondition = " and user_id in ? "
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? > 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(uIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, uIds, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	userIds := make([]int64, 0)
	for _, su := range sessionUsers {
		userIds = append(userIds, su.UserId)
	}
	return userIds
}

func (d defaultSessionUserModel) AddUser(session *Session, entityIds []int64, uIds []int64, role []int, maxCount int) (err error) {
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
	sqlStr := fmt.Sprintf("select count(0) from %s where session_id = ? and user_in not in ? and deleted = 0", tableName)
	if err = tx.Raw(sqlStr, session.Id, uIds).Scan(&count).Error; err != nil {
		return err
	}

	if count > maxCount-len(uIds) {
		return errorx.ErrGroupMemberCountBeyond
	}

	t := time.Now().UnixMilli()
	sql1 := "insert into " + d.genSessionUserTableName(session.Id) +
		" (session_id, user_id, role, type, create_time, update_time) values (?, ?, ?, ?, ?, ?) " +
		"on duplicate key update role = ?, deleted = ?, update_time = ? "

	for index, id := range uIds {
		if err = tx.Exec(sql1, session.Id, id, role[index], session.Type, t, t, role[index], 0, t).Error; err != nil {
			return err
		}
		sql2 := "insert into " + d.genUserSessionTableName(id) +
			" (session_id, user_id, role, type, entity_id, create_time, update_time) values (?, ?, ?, ?, ?, ?, ?) " +
			"on duplicate key update top = ?, role = ?, deleted = ?, update_time = ? "
		if err = tx.Exec(sql2, session.Id, id, role[index], session.Type, entityIds[index], t, t, 0, role[index], 0, t).Error; err != nil {
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

func (d defaultSessionUserModel) UpdateUser(sId int64, uIds []int64, role, status *int, mute *string, tx *gorm.DB) (err error) {
	if role == nil && status == nil && mute == nil {
		return nil
	}
	t := time.Now().UnixMilli()
	updates := make(map[string]interface{})
	if role != nil {
		updates["role"] = *role
	}
	if status != nil {
		updates["status"] = *status
	}
	if mute != nil {
		updates["mute"] = *mute
	}
	updates["update_time"] = t
	db := tx
	if db == nil {
		db = d.db
	}
	return tx.Table(d.genSessionUserTableName(sId)).Updates(updates).Where("session_id = ? and user_id in ?", sId, uIds).Error
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
