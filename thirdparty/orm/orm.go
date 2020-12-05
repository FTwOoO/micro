package orm

import (
	"github.com/FTwOoO/micro/thirdparty/jaeger/gormtracing"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var g_db *gorm.DB

func InitOrm(dbName string, dbSource string, ormLog bool, enableTracing bool) error {
	if dbName == "mysql" {
		source := mysql.Open(dbSource)
		db, err := gorm.Open(source, &gorm.Config{})
		if err != nil {
			return err
		}

		if ormLog {
			db = db.Debug()
		}
		db.AllowGlobalUpdate = true

		if enableTracing {
			err := db.Use(&gormtracing.OpentracingPlugin{})
			if err != nil {
				return err
			}
		}
		g_db = db
	}
	return nil
}
