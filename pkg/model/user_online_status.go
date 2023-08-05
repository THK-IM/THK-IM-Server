package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type (
	UserOnlineStatus struct {
		UserId     int64 `gorm:"user_id"`
		IsOnline   int8  `gorm:"is_online"`
		OnlineTime int64 `gorm:"online_time"`
	}

	UserOnlineStatusModel interface {
		GetUsersOnlineStatus(userIds []int64) ([]*UserOnlineStatus, error)
		UpdateUserOnlineStatus(userId int64, isOnline int8) error
		GetOnlineUsers(userIds []int64, offlineTime int64) ([]int64, error)
		GetOfflineUsers(userIds []int64, offlineTime int64) ([]int64, error)
	}

	defaultUserOnlineStatusModel struct {
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
		shards        int64
	}
)

func (d defaultUserOnlineStatusModel) GetUsersOnlineStatus(userIds []int64) ([]*UserOnlineStatus, error) {
	usersOnlineStatus := make([]*UserOnlineStatus, 0)
	sql := "select * from user_online_status where user_id in ?"
	err := d.db.Raw(sql, userIds).Scan(&usersOnlineStatus).Error
	return usersOnlineStatus, err
}

func (d defaultUserOnlineStatusModel) UpdateUserOnlineStatus(userId int64, isOnline int8) (err error) {
	now := time.Now().UnixMilli()
	sqlStr := "insert into user_online_status " +
		" (user_id, online_time, is_online) values (?, ?, ?)" +
		" on duplicate key update online_time = ?, is_online = ?"
	return d.db.Exec(sqlStr, userId, now, isOnline, now, isOnline).Error
}

func (d defaultUserOnlineStatusModel) GetOnlineUsers(userIds []int64, offlineTime int64) ([]int64, error) {
	sqlStr := "select user_id from user_online_status where user_id in ? and last_login_time >= ?"
	offlineUserIds := make([]int64, 0)
	tx := d.db.Raw(sqlStr, userIds, offlineTime).Scan(&offlineUserIds)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return offlineUserIds, nil
}

func (d defaultUserOnlineStatusModel) GetOfflineUsers(userIds []int64, offlineTime int64) ([]int64, error) {
	sqlStr := "select user_id from user_online_status where user_id in ? and last_login_time < ?"
	offlineUserIds := make([]int64, 0)
	tx := d.db.Raw(sqlStr, userIds, offlineTime).Scan(&offlineUserIds)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return offlineUserIds, nil
}

func (d defaultUserOnlineStatusModel) genUserOnlineStatusTable(userId int64) string {
	return "user_message_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func NewUserOnlineStatusModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserOnlineStatusModel {
	return defaultUserOnlineStatusModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
