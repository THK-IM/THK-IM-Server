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
	SessionObject struct {
		Id         int64  `gorm:"id" json:"id"`
		SId        int64  `gorm:"s_id" json:"s_id"`
		FromUserId int64  `gorm:"from_user_id" json:"from_user_id"`
		ClientId   int64  `gorm:"client_id" json:"client_id"`
		Engine     string `gorm:"engine" json:"engine"`
		Key        string `gorm:"key" json:"key"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
	}

	SessionObjectModel interface {
		AddSession(sId int64, fromUIds, clientMsgIds []int64, newFromUId, newClientMsgId, newSId int64) error
		Insert(sId, fromUId, clientId int64, engine, key string) (int64, error)
		FindObject(id, sId int64) (*SessionObject, error)
	}

	defaultSessionObjectModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionObjectModel) AddSession(sId int64, fromUIds, clientMsgIds []int64, newFromUId, newClientMsgId, newSId int64) error {
	db := d.db
	tableName := d.genObjectTableName(sId)
	sql := fmt.Sprintf("select * from %s where s_id = ? and from_user_id in ? and client_id in ?", tableName)
	objects := make([]*SessionObject, 0)
	err := db.Raw(sql, sId, fromUIds, clientMsgIds).Scan(&objects).Error
	if err != nil {
		return err
	}
	if len(objects) == 0 {
		return errorx.ErrParamsError
	}
	now := time.Now().UnixMilli()
	for i, object := range objects {
		object.Id = objects[i].Id
		object.SId = newSId
		object.FromUserId = newFromUId
		object.ClientId = newClientMsgId
		object.Key = objects[i].Key
		object.CreateTime = now
	}
	newTableName := d.genObjectTableName(newSId)
	err = db.Table(newTableName).CreateInBatches(objects, len(objects)).Error
	return err
}

func (d defaultSessionObjectModel) Insert(sId, fromUId, clientId int64, engine, key string) (int64, error) {
	id := d.snowflakeNode.Generate().Int64()
	o := &SessionObject{
		Id:         id,
		SId:        sId,
		FromUserId: fromUId,
		ClientId:   clientId,
		Engine:     engine,
		Key:        key,
		CreateTime: time.Now().UnixMilli(),
	}
	tableName := d.genObjectTableName(sId)
	return id, d.db.Table(tableName).Create(o).Error
}

func (d defaultSessionObjectModel) FindObject(id, sId int64) (*SessionObject, error) {
	tableName := d.genObjectTableName(sId)
	sql := fmt.Sprintf("select * from %s where id = ? and s_id = ? limit 0, 1", tableName)
	object := &SessionObject{}
	err := d.db.Raw(sql, id, sId).Scan(object).Error
	return object, err
}

func (d defaultSessionObjectModel) genObjectTableName(sId int64) string {
	return "session_object_" + fmt.Sprintf("%02d", sId%(d.shards))
}

func NewSessionObjectModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionObjectModel {
	return defaultSessionObjectModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
