package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sacf-rpcx/common/constant"
	"sacf-rpcx/common/jaeger/gormtracing"
)

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
		constant.Global_DB = db
	}
	return nil
}
