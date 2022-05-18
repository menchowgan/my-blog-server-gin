package db

import (
	"fmt"
	"gmc-blog-server/config"
	"gmc-blog-server/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB Repo

type Repo interface {
	i()
	// 获取读数据库连接
	GetDbR() *gorm.DB
	// 获取写数据库连接
	GetDbW() *gorm.DB
	// 关闭读数据库连接
	DbRClose() error
	// 关闭写数据库连接
	DbWClose() error
}

func InitDB() error {
	db, err := New()
	if err == nil {
		DB = db
		return nil
	}
	return err
}

type dbRepo struct {
	DbR *gorm.DB
	DbW *gorm.DB
}

func (d *dbRepo) i() {}

func (d *dbRepo) GetDbR() *gorm.DB {
	return d.DbR
}

func (d *dbRepo) GetDbW() *gorm.DB {
	return d.DbW
}

func (d *dbRepo) DbRClose() error {
	sqlDB, err := d.DbR.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *dbRepo) DbWClose() error {
	sqlDB, err := d.DbW.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// New 创建连接
func New() (Repo, error) {
	dbr, err := dbConnect(config.USER, config.PASSWORD, config.ADDR, config.DATABASE_NAME)
	if err != nil {
		return nil, err
	}

	dbw, err := dbConnect(config.USER, config.PASSWORD, config.ADDR, config.DATABASE_NAME)
	if err != nil {
		return nil, err
	}

	return &dbRepo{
		DbR: dbr,
		DbW: dbw,
	}, nil
}

// dbConnect 数据库连接
func dbConnect(user, pass, addr, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		user,
		pass,
		addr,
		dbName,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info), // 日志配置
	})

	if err != nil {
		return nil, err
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(60)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(60)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(10))

	return db, nil
}

func InitTables() {
	dw := DB.GetDbW()
	dr := DB.GetDbR()

	dw.Transaction(func(tx *gorm.DB) error {
		var err error
		has := dr.Migrator().HasTable(&model.User{})

		if !has {
			err = dw.AutoMigrate(&model.User{})
			if err != nil {
				return err
			}
		}

		has = dr.Migrator().HasTable(&model.Articles{})

		if !has {
			err = dw.AutoMigrate(&model.Articles{})
			if err != nil {
				return err
			}
		}

		has = dr.Migrator().HasTable(&model.Photos{})

		if !has {
			err = dw.AutoMigrate(&model.Photos{})
			if err != nil {
				return err
			}
		}

		has = dr.Migrator().HasTable(&model.Music{})

		if !has {
			err = dw.AutoMigrate(&model.Music{})
			if err != nil {
				return err
			}
		}

		has = dr.Migrator().HasTable(&model.Video{})

		if !has {
			err = dw.AutoMigrate(&model.Video{})
			if err != nil {
				return err
			}
		}

		return nil
	})

}
