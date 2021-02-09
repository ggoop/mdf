package repositories

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/ggoop/mdf/configs"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/utils"

	"github.com/ggoop/mdf/gorm"
)

type MysqlRepo struct {
	*gorm.DB
}

var _default *MysqlRepo

func Default() *MysqlRepo {
	if _default == nil {
		_default = NewMysqlRepo()
	}
	return _default
}

func SetDefault(d *MysqlRepo) {
	_default = d
}
func Open() *MysqlRepo {
	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}

	db.LogMode(configs.Default.App.Debug)
	repo := createMysqlRepo(db)
	return repo
}
func (s *MysqlRepo) Close() error {
	return s.DB.Close()
}
func (s *MysqlRepo) Begin() *MysqlRepo {
	return &MysqlRepo{s.DB.Begin()}
}
func (s *MysqlRepo) New() *MysqlRepo {
	return &MysqlRepo{s.DB.New()}
}
func NewMysqlRepo() *MysqlRepo {
	// 生成数据库
	CreateDB(configs.Default.Db.Database)

	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
		panic(err)
	}

	db.LogMode(configs.Default.App.Debug)
	repo := createMysqlRepo(db)

	if _default == nil {
		_default = repo
	}
	return repo
}
func getDsnString(inDb bool) string {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//mssql:   =>  sqlserver://username:password@localhost:1433?database=dbname
	//mysql => user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	str := ""
	// 创建连接
	if configs.Default.Db.Driver == utils.ORM_DRIVER_MSSQL {
		var buf bytes.Buffer
		buf.WriteString("sqlserver://")
		buf.WriteString(configs.Default.Db.Username)
		if configs.Default.Db.Password != "" {
			buf.WriteByte(':')
			buf.WriteString(configs.Default.Db.Password)
		}
		buf.WriteByte('@')
		if configs.Default.Db.Host != "" {
			buf.WriteString(configs.Default.Db.Host)
			if configs.Default.Db.Port != "" {
				buf.WriteByte(':')
				buf.WriteString(configs.Default.Db.Port)
			} else {
				buf.WriteString(":1433")
			}
		}
		if configs.Default.Db.Database != "" && inDb {
			buf.WriteString("?database=")
			buf.WriteString(configs.Default.Db.Database)
		} else {
			buf.WriteString("?database=master")
		}
		str = buf.String()
		return str
	}
	{
		config := mysql.Config{
			User:   configs.Default.Db.Username,
			Passwd: configs.Default.Db.Password, Net: "tcp", Addr: configs.Default.Db.Host,
			AllowNativePasswords: true,
			ParseTime:            true,
			Loc:                  time.Local,
		}
		if inDb {
			config.DBName = configs.Default.Db.Database
		}
		if configs.Default.Db.Port != "" {
			config.Addr = fmt.Sprintf("%s:%s", configs.Default.Db.Host, configs.Default.Db.Port)
		}
		str = config.FormatDSN()
	}
	return str
}
func createMysqlRepo(db *gorm.DB) *MysqlRepo {
	repo := &MysqlRepo{db}
	return repo
}
func DestroyDB(name string) error {
	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	defer db.Close()
	return db.Exec(fmt.Sprintf("Drop Database if exists %s;", name)).Error
}
func CreateDB(name string) {
	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	script := ""
	if configs.Default.Db.Driver == utils.ORM_DRIVER_MSSQL {
		script = fmt.Sprintf("if not exists (select * from sysdatabases where name='%s') begin create database %s end;", name, name)
	} else {
		script = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s;", name, configs.Default.Db.Charset, configs.Default.Db.Collation)
	}
	err = db.Exec(script).Error
	if err != nil {
		glog.Errorf("create DATABASE err: %v", err)
	}

	defer db.Close()
}
