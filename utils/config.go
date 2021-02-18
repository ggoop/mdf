package utils

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"strings"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/gmap"
)

type AppConfig struct {
	Token   string `mapstructure:"token" json:"token"`
	Code    string `mapstructure:"code" json:"code"`
	Active  string `mapstructure:"active" json:"active"`
	Name    string `mapstructure:"name" json:"name"`
	Port    string `mapstructure:"port" json:"port"`
	Locale  string `mapstructure:"locale" json:"locale"`
	Debug   bool   `mapstructure:"debug" json:"debug"`
	Storage string `mapstructure:"storage" json:"storage"`
	//注册中心
	Registry string `mapstructure:"registry" json:"registry"`
	//服务地址，带端口号
	Address       string `mapstructure:"address" json:"address"`
	PublicAddress string `mapstructure:"public_address" json:"public_address"`
}
type DbConfig struct {
	Driver    string `mapstructure:"driver" json:"driver"`
	Host      string `mapstructure:"host" json:"host"`
	Port      string `mapstructure:"port" json:"port"`
	Database  string `mapstructure:"database" json:"database"`
	Username  string `mapstructure:"username" json:"username"`
	Password  string `mapstructure:"password" json:"password"`
	Charset   string `mapstructure:"charset" json:"charset"`
	Collation string `mapstructure:"collation" json:"collation"`
}
type LogConfig struct {
	Level string `mapstructure:"level" json:"level"`
	Path  string `mapstructure:"path" json:"path"`
	Stack bool   `mapstructure:"stack" json:"stack"`
}
type AuthConfig struct {
	//权限中心地址
	Address string `mapstructure:"address" json:"address"`
	//权限中心编码
	Code string `mapstructure:"code" json:"code"`
}

const AppConfigName = "app"

var DefaultConfig *Config
var _ENVMaps = gmap.New()

type jsonConfig struct {
	App  AppConfig              `mapstructure:"app" json:"app"`
	Db   DbConfig               `mapstructure:"db" json:"db"`
	Log  LogConfig              `mapstructure:"log" json:"log"`
	Auth AuthConfig             `mapstructure:"auth" json:"auth"`
	Data map[string]interface{} `json:"data"`
}

type Config struct {
	App  AppConfig  `mapstructure:"app" json:"app"`
	Db   DbConfig   `mapstructure:"db" json:"db"`
	Log  LogConfig  `mapstructure:"log" json:"log"`
	Auth AuthConfig `mapstructure:"auth" json:"auth"`
	data map[string]interface{}
}

func (c Config) MarshalJSON() ([]byte, error) {
	jsonMap := jsonConfig{}
	jsonMap.App = c.App
	jsonMap.Db = c.Db
	jsonMap.Log = c.Log
	jsonMap.Auth = c.Auth
	jsonMap.Data = c.data
	return json.Marshal(jsonMap)
}

func (c *Config) UnmarshalJSON(b []byte) error {
	jsonMap := jsonConfig{}
	json.Unmarshal(b, &jsonMap)
	c.App = jsonMap.App
	c.Db = jsonMap.Db
	c.Log = jsonMap.Log
	c.Auth = jsonMap.Auth
	c.data = jsonMap.Data
	return nil
}
func (s *Config) GetValue(name string, envNames ...string) string {
	return s.getViper(envNames...).GetString(name)
}
func (s *Config) UnmarshalValue(name string, rawVal interface{}, envNames ...string) error {
	return s.getViper(envNames...).UnmarshalKey(name, rawVal)
}
func (s *Config) Unmarshal(rawVal interface{}, envNames ...string) error {
	return s.getViper(envNames...).Unmarshal(rawVal)
}
func (s *Config) GetBool(name string, envNames ...string) bool {
	return s.getViper(envNames...).GetBool(name)
}
func (s *Config) GetObject(name string, envNames ...string) interface{} {
	return s.getViper(envNames...).Get(name)
}
func (s *Config) SetValue(name string, value interface{}, envNames ...string) *Config {
	s.getViper(envNames...).Set(name, value)
	return s
}
func (s *Config) getViper(envNames ...string) *viper.Viper {
	if len(envNames) > 0 {
		return getConfigViper(envNames[0])
	} else {
		return getConfigViper(AppConfigName)
	}
}
func getConfigViper(name string) *viper.Viper {
	if name == "" {
		name = AppConfigName
	}
	name = strings.ToLower(name)
	if v, ok := _ENVMaps.Get(name); ok && v != nil {
		if vv, ok := v.(*viper.Viper); ok && vv != nil {
			return vv
		}
	}

	envType := "yaml"
	envPath := "env"

	envFile := JoinCurrentPath(path.Join(envPath, name+"."+envType))
	if !PathExists(envFile) {
		paths := []string{path.Join(envPath, name+".prod."+envType), path.Join(envPath, name+".dev."+envType)}
		for _, p := range paths {
			//不存在时，自动由dev创建
			if !PathExists(envFile) {
				devFile := JoinCurrentPath(p)
				if PathExists(devFile) {
					if s, err := os.Open(devFile); err == nil {
						defer s.Close()
						if newEnv, err := os.Create(envFile); err == nil {
							defer newEnv.Close()
							io.Copy(newEnv, s)
							break
						}
					}
				}
			}
		}
	}
	v := viper.New()
	v.SetConfigType(envType)
	v.SetConfigName(name)
	v.AddConfigPath(JoinCurrentPath(envPath))
	v.SetConfigFile(envFile)

	if err := v.ReadInConfig(); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", name, err)
	}
	_ENVMaps.Set(name, v)
	return v
}
func newInitConfig() {
	DefaultConfig = &Config{}
	vp := getConfigViper(AppConfigName)
	if err := vp.Unmarshal(&DefaultConfig); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if DefaultConfig.App.Port == "" {
		DefaultConfig.App.Port = "8080"
	}
	if DefaultConfig.App.Locale == "" {
		DefaultConfig.App.Locale = "zh-cn"
	}
	if DefaultConfig.App.Token == "" {
		DefaultConfig.App.Token = "01e8f6a984101b20b24e4d172ec741be"
	}
	if DefaultConfig.App.Storage == "" {
		DefaultConfig.App.Storage = "./storage"
	}
	if DefaultConfig.Db.Driver == "" {
		DefaultConfig.Db.Driver = "mysql"
	}
	if DefaultConfig.Db.Host == "" {
		DefaultConfig.Db.Host = "localhost"
	}
	if DefaultConfig.Db.Port == "" {
		if DefaultConfig.Db.Driver == "mysql" {
			DefaultConfig.Db.Port = "3306"
		}
		if DefaultConfig.Db.Driver == "mssql" {
			DefaultConfig.Db.Port = "1433"
		}
	}
	if DefaultConfig.Db.Charset == "" {
		DefaultConfig.Db.Charset = "utf8mb4"
	}
	if DefaultConfig.Db.Collation == "" {
		DefaultConfig.Db.Collation = "utf8mb4_general_ci"
	}
	if DefaultConfig.Auth.Code == "" {
		DefaultConfig.Auth.Code = DefaultConfig.App.Code
	}
	kvs := make(map[string]interface{})
	if err := vp.Unmarshal(&kvs); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if len(kvs) > 0 {
		for k, v := range kvs {
			DefaultConfig.SetValue(k, v)
		}
	}
}
