package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type (
	Object struct {
		Id         int64  `gorm:"id" json:"id"`
		SId        int64  `gorm:"s_id" json:"s_id"`
		Engine     string `gorm:"engine" json:"engine"`
		Key        string `gorm:"key" json:"key"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
	}

	ObjectModel interface {
		Insert(sId int64, engine, key string) (int64, error)
		FindOne(id int64) (*Object, error)
	}

	defaultObjectModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultObjectModel) Insert(sId int64, engine, key string) (int64, error) {
	id := d.snowflakeNode.Generate().Int64()
	o := &Object{
		Id:         id,
		SId:        sId,
		Engine:     engine,
		Key:        key,
		CreateTime: time.Now().UnixMilli(),
	}
	tableName := d.genObjectTableName(id)
	return id, d.db.Table(tableName).Create(o).Error
}

func (d defaultObjectModel) FindOne(id int64) (*Object, error) {
	tableName := d.genObjectTableName(id)
	sql := fmt.Sprintf("select * from %s where id = ?", tableName)
	o := &Object{}
	err := d.db.Raw(sql, id).Scan(o).Error
	return o, err
}

func (d defaultObjectModel) genObjectTableName(id int64) string {
	return "object_" + fmt.Sprintf("%02d", id%(d.shards))
}

func NewObjectModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) ObjectModel {
	return defaultObjectModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
