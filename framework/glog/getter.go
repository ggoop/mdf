package glog

import (
	"database/sql/driver"
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	oracleRegexp             = regexp.MustCompile(`\:\d+`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
	Stack bool   `mapstructure:"stack"`
}

func readConfig() *LogConfig {
	config := &LogConfig{}
	viper.SetConfigType("yaml")

	viper.SetConfigName("app")
	viper.AddConfigPath(joinCurrentPath("env"))
	if err := viper.ReadInConfig(); err != nil {
		//Errorf("Fatal error when reading %s config file:%s", "app", err)
	}
	if err := viper.UnmarshalKey("log", config); err != nil {
		//Errorf("Fatal error when reading %s config file:%s", "app", err)
	}
	if config.Path == "" {
		config.Path = "./storage/logs"
	}
	return config
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
func NowFunc() time.Time {
	return time.Now()
}

func (l *logger) sqlLog(values ...interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
		)

		if level == "sql" {
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						switch value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						default:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						}
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if oracleRegexp.MatchString(values[3].(string)) {
				formattedValuesLength := len(formattedValues)
				for index, value := range oracleRegexp.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			} else if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}
			l.Debug(sql,
				String("type", "sql"),
				String("rows", strconv.FormatInt(values[5].(int64), 10)),
				String("duration", fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)))

		}
	}

	return
}
