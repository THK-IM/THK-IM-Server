package loader

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

func LoadDataBase(source *conf.DataSource) *gorm.DB {
	if source == nil {
		return nil
	}
	endpoint := source.Endpoint
	if strings.HasPrefix(source.Endpoint, "{{") && strings.HasSuffix(source.Endpoint, "}}") {
		endpointEnvKey := strings.Replace(source.Endpoint, "{{", "", -1)
		endpointEnvKey = strings.Replace(endpointEnvKey, "}}", "", -1)
		endpoint = os.Getenv(endpointEnvKey)
	}
	db, errOpen := gorm.Open(mysql.Open(fmt.Sprintf("%s%s", endpoint, source.Uri)), &gorm.Config{})
	if errOpen != nil {
		panic(errOpen)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxIdleConns(source.MaxIdleConn)
	sqlDb.SetMaxOpenConns(source.MaxOpenConn)
	sqlDb.SetConnMaxIdleTime(time.Duration(source.ConnMaxIdleTime) * time.Second)
	sqlDb.SetConnMaxLifetime(time.Duration(source.ConnMaxLifeTime) * time.Second)

	return db
}
