package db

import (
	"admin/setting"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"time"
)

type MysqlPool struct {

}

var instance *MysqlPool
var once sync.Once
var err error
var db  *gorm.DB

func GetInstance() *MysqlPool {
	if instance != nil {
		return instance
	}
	once.Do(func() {
		instance = &MysqlPool{}
	})

	return instance
}

func (mysqlPool *MysqlPool) InitMysqlPool() *MysqlPool {
	var dbStr string
	if setting.DBType == "mysql" {
		dbStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&&parseTime=True&loc=Local",
			setting.Username,
			setting.Password,
			setting.Host,
			setting.Port,
			setting.DbName)
	} else {
		dbStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=require sslcert=%s sslkey=%s sslrootcert=%s",
			setting.Username,
			setting.Password,
			setting.DbName,
			setting.Host,
			setting.Port,
			setting.SslCert,
			setting.SslKey,
			setting.SslRootCert)
	}

	db, err = gorm.Open(setting.DBType, dbStr)
	if err != nil {
		log.Fatal(err)
		return mysqlPool
	}

	db.DB().SetConnMaxLifetime(100*time.Second)  //最大连接周期，超过时间的连接就close
	db.DB().SetMaxOpenConns(20) //设置最大连接数
	db.DB().SetMaxIdleConns(4) //设置闲置连接数
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return setting.TablePrefix + defaultTableName
	}

	if err := db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	return mysqlPool
}

func (mysqlPool *MysqlPool) GetMysqlDb() *gorm.DB {
	return db
}

func (mysqlPool *MysqlPool) Close() error {
	 err := db.Close()
	 if err != nil {
		 log.Fatal(err)
	 }
	 return err
}


