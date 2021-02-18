package db

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/db/gorm"
)

type Repo struct {
	*gorm.DB
}

var dbIns *Repo
var _cacheDefaultMu sync.Once

func Default() *Repo {
	if dbIns == nil {
		_cacheDefaultMu.Do(func() {
			dbIns = NewMysqlRepo()
		})
	}
	return dbIns
}
func SetDefault(d *Repo) {
	dbIns = d
}
func Open() *Repo {
	db, err := gorm.Open(utils.DefaultConfig.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}

	db.LogMode(utils.DefaultConfig.App.Debug)
	repo := &Repo{db}
	return repo
}
func (s *Repo) Close() error {
	return s.DB.Close()
}
func (s *Repo) Begin() *Repo {
	return &Repo{s.DB.Begin()}
}
func (s *Repo) New() *Repo {
	return &Repo{s.DB.New()}
}
func NewMysqlRepo() *Repo {
	db, err := gorm.Open(utils.DefaultConfig.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
		panic(err)
	}
	db.LogMode(utils.DefaultConfig.App.Debug)
	repo := &Repo{db}
	return repo
}
func getDsnString(inDb bool) string {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//mssql:   =>  sqlserver://username:password@localhost:1433?database=dbname
	//mysql => user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	str := ""
	// 创建连接
	if utils.DefaultConfig.Db.Driver == utils.ORM_DRIVER_MSSQL {
		var buf bytes.Buffer
		buf.WriteString("sqlserver://")
		buf.WriteString(utils.DefaultConfig.Db.Username)
		if utils.DefaultConfig.Db.Password != "" {
			buf.WriteByte(':')
			buf.WriteString(utils.DefaultConfig.Db.Password)
		}
		buf.WriteByte('@')
		if utils.DefaultConfig.Db.Host != "" {
			buf.WriteString(utils.DefaultConfig.Db.Host)
			if utils.DefaultConfig.Db.Port != "" {
				buf.WriteByte(':')
				buf.WriteString(utils.DefaultConfig.Db.Port)
			} else {
				buf.WriteString(":1433")
			}
		}
		if utils.DefaultConfig.Db.Database != "" && inDb {
			buf.WriteString("?database=")
			buf.WriteString(utils.DefaultConfig.Db.Database)
		} else {
			buf.WriteString("?database=master")
		}
		str = buf.String()
		return str
	}
	{
		config := mysql.Config{
			User:   utils.DefaultConfig.Db.Username,
			Passwd: utils.DefaultConfig.Db.Password, Net: "tcp", Addr: utils.DefaultConfig.Db.Host,
			AllowNativePasswords: true,
			ParseTime:            true,
			Loc:                  time.Local,
		}
		if inDb {
			config.DBName = utils.DefaultConfig.Db.Database
		}
		if utils.DefaultConfig.Db.Port != "" {
			config.Addr = fmt.Sprintf("%s:%s", utils.DefaultConfig.Db.Host, utils.DefaultConfig.Db.Port)
		}
		str = config.FormatDSN()
	}
	return str
}
func DestroyDB(name string) error {
	db, err := gorm.Open(utils.DefaultConfig.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	defer db.Close()
	return db.Exec(fmt.Sprintf("Drop Database if exists %s;", name)).Error
}
func CreateDB(name string) {
	db, err := gorm.Open(utils.DefaultConfig.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	script := ""
	if utils.DefaultConfig.Db.Driver == utils.ORM_DRIVER_MSSQL {
		script = fmt.Sprintf("if not exists (select * from sysdatabases where name='%s') begin create database %s end;", name, name)
	} else {
		script = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s;", name, utils.DefaultConfig.Db.Charset, utils.DefaultConfig.Db.Collation)
	}
	err = db.Exec(script).Error
	if err != nil {
		glog.Errorf("create DATABASE err: %v", err)
	}

	defer db.Close()
}
