package dti

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ggoop/mdf/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

type ParamValue struct {
	Name  string
	Value interface{}
}
type DtiHandProc struct {
	Group string
	Ctx   *gin.Context
	Path  string
	db    *sql.DB
}

func (c *DtiHandProc) getDSN() (driverName, dataSourceName string, err error) {
	dbConfig := utils.DbConfig{}
	if c.Group == "dti" {
		utils.DefaultConfig.UnmarshalValue("dti", &dbConfig)
	} else if c.Group == "localhost" || c.Group == "db" {
		utils.DefaultConfig.UnmarshalValue("db", &dbConfig)
	} else {
		utils.DefaultConfig.UnmarshalValue("dti."+c.Group, &dbConfig)
	}
	if dbConfig.Driver == "sqlserver" {
		dataSourceName, err = c.getDSN_sqlserver(dbConfig)
		return dbConfig.Driver, dataSourceName, err
	} else if dbConfig.Driver == "oci8" {
		dataSourceName, err = c.getDSN_oci8(dbConfig)
	} else if dbConfig.Driver == "mysql" {
		dataSourceName, err = c.getDSN_mysql(dbConfig)
	}
	return dbConfig.Driver, dataSourceName, err
}

func (c *DtiHandProc) getDSN_mysql(dbConfig utils.DbConfig) (dataSourceName string, err error) {
	config := mysql.Config{
		User:   dbConfig.Username,
		Passwd: dbConfig.Password, Net: "tcp", Addr: dbConfig.Host,
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}
	config.DBName = dbConfig.Database
	if dbConfig.Port != "" {
		config.Addr = fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port)
	}
	return config.FormatDSN(), nil
}

func (c *DtiHandProc) getDSN_oci8(dbConfig utils.DbConfig) (dataSourceName string, err error) {
	//oci8 =[username/[password]@]host[:port][/service_name][?param1=value1&...&paramN=valueN]
	//"oci8", "j1_bibox/oracle123@10.1.196.200:1521/zjdevdb"
	s := strings.ReplaceAll(url.UserPassword(dbConfig.Username, dbConfig.Password).String(), ":", "/")
	u := &url.URL{
		Host: dbConfig.Host,
		Path: dbConfig.Database,
	}
	if utils.DefaultConfig.Db.Port != "" {
		u.Host = fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port)
	}
	return s + "@" + u.String(), nil
}
func (c *DtiHandProc) getDSN_sqlserver(dc utils.DbConfig) (driverName string, err error) {
	//sqlserver://sa:mypass@localhost:1234?database=master&connection+timeout=30
	//sqlserver://username:password@localhost:1433?database=dbname
	query := url.Values{}
	query.Add("encrypt", "disable")
	if dc.Database != "" {
		query.Add("database", dc.Database)
	}
	query.Add("connection timeout", "0")
	u := &url.URL{
		Scheme:   dc.Driver,
		User:     url.UserPassword(dc.Username, dc.Password),
		Host:     dc.Host,
		RawQuery: query.Encode(),
	}
	if utils.DefaultConfig.Db.Port != "" {
		u.Host = fmt.Sprintf("%s:%s", dc.Host, dc.Port)
	}
	return u.String(), nil
}
func (c *DtiHandProc) Do() {
	driverName, ds, err := c.getDSN()
	if err != nil {
		c.toError(err)
		return
	}
	glog.Errorf("请求连接为:%s", driverName, ds)
	if driverName == "" || ds == "" {
		c.toError("连接为空")
		return
	}
	bodyParams := make(map[string]interface{})
	paramKey := ""
	otherParams := make(map[string]interface{})
	if err := c.Ctx.Bind(&otherParams); err != nil {
		glog.Error(err)
	}
	if len(otherParams) > 0 {
		for k, v := range otherParams {
			paramKey = strings.ToLower(k)
			if _, ok := bodyParams[paramKey]; !ok {
				bodyParams[paramKey] = v
			}
		}
	}
	glog.Errorf("解析到请求 path:%s,参数为:%v", c.Path, bodyParams)
	db, err := sql.Open(driverName, ds)
	defer db.Close()
	c.db = db
	if err != nil {
		glog.Error(err)
		c.toError(err)
		return
	}
	paramsIn := make([]ParamValue, 0)
	rtn := utils.Map{"data": nil, "path": c.Path, "name": c.Ctx.GetHeader("REMOTE_CODE")}
	if script, ok := bodyParams["script"]; ok && script != nil && script.(string) != "" {
		datas, err := c.execQuery(script.(string), utils.ToInterfaceSlice(bodyParams["params"])...)
		if err != nil {
			c.toError(err)
			return
		}
		rtn["data"] = datas
	} else {
		if spParams, err := c.getSpParams(driverName, c.Path); err != nil {
			glog.Error(err)
			c.toError(err)
			return
		} else {
			glog.Errorf("存储过程:%s,参数为:%v", c.Path, spParams)
			for _, pv := range spParams {
				if pv["name"] == nil {
					break
				}
				paramKey = pv["name"].(string)
				if pv, ok := bodyParams[strings.ToLower(paramKey)]; ok && pv != nil {
					paramsIn = append(paramsIn, ParamValue{Name: paramKey, Value: pv})
				}
			}
		}
		// 执行SQL语句
		fm_time := time.Now()
		glog.Errorf("存储过程:%s,传入参数:%v,开始执行！", c.Path, paramsIn)
		maps, err := c.execProc(driverName, c.Path, paramsIn)
		glog.Errorf("执行:%s 结束,%v条,time:%v Seconds", c.Path, len(maps), time.Now().Sub(fm_time).Seconds())
		if err != nil {
			glog.Error(err)
			c.toError(err)
			return
		}
		rtn["data"] = maps
	}
	c.Ctx.JSON(http.StatusOK, rtn)
}
func (c *DtiHandProc) toError(data interface{}, code ...int) {
	obj := utils.Map{}
	if ev, ok := data.(utils.GError); ok {
		obj["msg"] = ev.Error()
	} else if ev, ok := data.(error); ok {
		obj["msg"] = ev.Error()
	} else {
		obj["msg"] = data
	}
	if code != nil && len(code) > 0 {
		obj["code"] = code[0]
	}
	c.Ctx.JSON(http.StatusBadRequest, obj)
}
func (c *DtiHandProc) getSpParams(driverName, name string) ([]map[string]interface{}, error) {
	cmd := fmt.Sprintf(`select o.object_id,substring(p.name,2,100) as name,lt.name as type,lt.max_length as length 
	from  sys.objects  o 
	left join sys.parameters  p on p.object_id=o.object_id
	left join sys.types  lt on p.system_type_id=lt.user_type_id
	where o.type='P' and o.name='%s'
	Order By p.parameter_id`, name)
	spParams, err := c.execQuery(cmd)
	if err != nil {
		return nil, err
	}
	if spParams == nil || len(spParams) == 0 {
		return nil, errors.New(fmt.Sprintf("找不到实现:%v", name))
	}
	return spParams, nil
}
func (c *DtiHandProc) execProc(driverName, name string, params []ParamValue) ([]map[string]interface{}, error) {
	cmd := fmt.Sprintf("exec %s", name)
	paramValues := make([]interface{}, 0)
	for i, k := range params {
		if i == 0 {
			cmd = fmt.Sprintf("%v @%v=@p%v", cmd, k.Name, i+1)
		} else {
			cmd = fmt.Sprintf("%v,@%v=@p%v", cmd, k.Name, i+1)
		}
		paramValues = append(paramValues, k.Value)
	}
	stmt, err := c.db.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result := paramValues[:]
	rows, err := stmt.Query(result...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return c.rowsToMap(rows)
}
func (c *DtiHandProc) execQuery(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return c.rowsToMap(rows)
}
func (c *DtiHandProc) rowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	var maps = make([]map[string]interface{}, 0)
	colNames, _ := rows.Columns()
	var cols = make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		cols[i] = new(interface{})
	}
	for rows.Next() {
		err := rows.Scan(cols...)
		if err != nil {
			return nil, err
		}
		var rowMap = make(map[string]interface{})
		for i := 0; i < len(colNames); i++ {
			rowMap[colNames[i]] = c.convertRow(*(cols[i].(*interface{})))
		}
		maps = append(maps, rowMap)
	}
	return maps, nil
}
func (c *DtiHandProc) convertRow(row interface{}) interface{} {
	switch row.(type) {
	case int:
		return utils.ToInt(row)
	case string:
		return utils.ToString(row)
	case []byte:
		return utils.ToString(row)
	case bool:
		return utils.ToBool(row)
	}
	return row
}
