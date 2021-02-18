package files

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/ggoop/mdf/utils"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ggoop/mdf/framework/glog"
)

type FileData struct {
	FileName string
	Dir      string
	FullPath string
}

func (s FileData) String() string {
	return s.FileName
}

type ImportData struct {
	EntityCode string
	EntityName string
	SheetName  string
	Columns    map[string]string
	Data       []map[string]string
}

func (s ImportData) String() string {
	return fmt.Sprintf("%s %s", s.SheetName, s.EntityCode)
}

type ExcelColumn struct {
	Name      string
	Title     string
	Hidden    bool
	excelName string
}
type ExcelCell struct {
	Name  string
	Value interface{}
}
type ToExcel struct {
	FileName string
	Dir      string
	Columns  []ExcelColumn
	Datas    [][]ExcelCell
}

func (s ToExcel) String() string {
	return s.FileName
}

type ExcelSv struct {
}

/**
* 创建服务实例
 */
func NewExcelSv() *ExcelSv {
	return &ExcelSv{}
}
func GetCellValue(key string, row map[string]string) string {
	if v, ok := row[key]; ok {
		return v
	}
	if v, ok := row[utils.SnakeString(key)]; ok {
		return v
	}
	return ""
}
func (s *ExcelSv) GetExcelDatasByReader(r io.Reader, sheetNames ...string) ([]ImportData, error) {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	return s.toExcelDatas(xlsx, sheetNames...)
}
func (s *ExcelSv) GetExcelDatas(filePath string, sheetNames ...string) ([]ImportData, error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	return s.toExcelDatas(xlsx, sheetNames...)
}
func (s *ExcelSv) toExcelDatas(xlsx *excelize.File, sheetNames ...string) ([]ImportData, error) {
	rtnDatas := make([]ImportData, 0)
	if len(sheetNames) == 0 || sheetNames[0] == "*" {
		for _, sheetName := range xlsx.GetSheetMap() {
			if data, err := s.getSheetData(xlsx, sheetName); err != nil {
				return nil, err
			} else if len(data) > 0 {
				for i, _ := range data {
					rtnDatas = append(rtnDatas, data[i])
				}
			}
		}
	} else {
		for _, sheetName := range sheetNames {
			if data, err := s.getSheetData(xlsx, sheetName); err != nil {
				return nil, err
			} else if len(data) > 0 {
				for i, _ := range data {
					rtnDatas = append(rtnDatas, data[i])
				}
			}
		}
	}
	return rtnDatas, nil
}
func (s *ExcelSv) GetExcelDataByReader(r io.Reader) (data ImportData, err error) {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return data, err
	}
	return s.toExcelData(xlsx)
}
func (s *ExcelSv) GetExcelData(filePath string) (data ImportData, err error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return data, err
	}
	return s.toExcelData(xlsx)
}
func (s *ExcelSv) toExcelData(xlsx *excelize.File) (data ImportData, err error) {
	datas, err := s.getSheetData(xlsx, xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	if err != nil {
		return data, err
	}
	if len(datas) > 0 && len(datas[0].Data) > 0 {
		return datas[0], nil
	}
	return data, err
}
func (s *ExcelSv) getSheetData(xlsx *excelize.File, sheetName string) ([]ImportData, error) {
	entityDatas := make([]ImportData, 0)
	idMap := make(map[string]string)
	//第一行为字段标识
	//条二行为行名称
	// 获取 Sheet1 上所有单元格,模板，需要预制一个列标识，_state,为空后，后边的行将都不会导入
	if sheetName == "" {
		return entityDatas, nil
	}
	allRows := xlsx.GetRows(sheetName)
	if len(allRows) <= 2 {
		return entityDatas, nil
	}
	rowCount := len(allRows)
	hasEntity := false
	//判断是否标准导入模板，标准导入模型有实体标识，列标识，标题标识
	for _, row := range allRows {
		firstValue := row[0]
		if strings.HasPrefix(firstValue, "[") && strings.HasSuffix(firstValue, "]") {
			hasEntity = true
			break
		}
	}
	isData := true
	//取列数
	colsMap := make(map[int]string)
	if hasEntity {
		currentPart := ImportData{SheetName: sheetName, Columns: make(map[string]string), Data: make([]map[string]string, 0)}
		for i := 0; i < rowCount; i++ {
			row := allRows[i]
			values := make(map[string]string, 0)
			firstValue := row[0]

			//如果是空行
			if firstValue == "" {
				if currentPart.EntityCode != "" && len(currentPart.Data) > 0 {
					entityDatas = append(entityDatas, currentPart)
				}
				currentPart = ImportData{SheetName: sheetName, Columns: make(map[string]string), Data: make([]map[string]string, 0)}
				continue
			}
			//如果是表标记，则取出表名
			if firstValue != "" && strings.HasPrefix(firstValue, "[") && strings.HasSuffix(firstValue, "]") {
				currentPart.EntityCode = strings.ReplaceAll(strings.ReplaceAll(firstValue, "[", ""), "]", "")
				if len(row) > 1 {
					currentPart.EntityName = row[1]
				}
				//如果没有配置数据行或者标题，栏目，则退出
				if i+2 >= rowCount {
					break
				}
				colsMap = make(map[int]string)
				cols := allRows[i+1]   //取出列名,表名下一行为列名
				titles := allRows[i+2] //取出标题

				for c, name := range cols {
					if name != "" {
						colsMap[c] = name
						currentPart.Columns[name] = titles[c]
					}
				}
				//跳过标题和栏目行
				i = i + 2
				continue
			}
			//如果是数据行
			if firstValue != "" && currentPart.EntityCode != "" {
				isData = true
				for c, value := range row {
					if cName, ok := colsMap[c]; ok && cName != "" {
						if cName == utils.STATE_FIELD && value == utils.STATE_IGNORED {
							isData = false
							break
						}
						values[cName] = s.getCellValue(value, idMap)
					}
				}
				if isData {
					currentPart.Data = append(currentPart.Data, values)
				}
			}
		}
		if currentPart.EntityCode != "" && len(currentPart.Data) > 0 {
			entityDatas = append(entityDatas, currentPart)
		}
	} else {
		currentPart := ImportData{SheetName: sheetName, Columns: make(map[string]string), Data: make([]map[string]string, 0)}
		cols := allRows[0]
		titles := allRows[1]
		for c, name := range cols {
			if name != "" {
				colsMap[c] = name
				currentPart.Columns[name] = titles[c]
			}
		}
		if len(colsMap) == 0 {
			return entityDatas, nil
		}
		datas := make([]map[string]string, 0)
		isData := false
		for i := 2; i < rowCount; i++ {
			row := allRows[i]
			isData = false
			values := make(map[string]string, 0)
			for c, value := range row {
				if cName, ok := colsMap[c]; ok {
					if c == 0 && value != "" {
						isData = true
					}
					values[cName] = value
				}
			}
			if isData {
				datas = append(datas, values)
			} else {
				break
			}
		}
		currentPart.Data = datas
		if len(datas) > 0 {
			entityDatas = append(entityDatas, currentPart)
		}
	}
	return entityDatas, nil
}

func (s *ExcelSv) getCellValue(value string, idMap map[string]string) string {
	if value == "" {
		return value
	}
	if strings.HasSuffix(value, "*") && strings.HasPrefix(value, "*") {
		if v, ok := idMap[value]; ok {
			value = v
			return v
		} else {
			v = utils.GUID()
			idMap[value] = v
			return v
		}
	}
	return value
}
func (s *ExcelSv) ToExcel(data *ToExcel) (*FileData, error) {

	xlsx := excelize.NewFile()
	sheetName := xlsx.GetSheetName(xlsx.GetActiveSheetIndex())
	colMap := make(map[string]ExcelColumn)
	//增加系统默认导出列
	columns := make([]ExcelColumn, 0)
	for _, c := range data.Columns {
		columns = append(columns, c)
	}
	startIndex := 2
	for i, c := range columns {
		cName := excelize.ToAlphaString(i)
		columns[i].excelName = cName
		colMap[c.Name] = columns[i]

		xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", cName, 1), c.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", cName, 2), c.Title)
	}
	xlsx.SetRowVisible(sheetName, 0, false)
	//设置数据列宽度
	xlsx.SetColWidth("Sheet1", "A", columns[len(columns)-1].excelName, 20)
	//border
	if style, err := xlsx.NewStyle(`{"border":[{"type":"left","color":"666666","style":1},{"type":"top","color":"666666","style":1},{"type":"bottom","color":"666666","style":1},{"type":"right","color":"666666","style":1}]}`); err == nil {
		xlsx.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", columns[len(columns)-1].excelName, len(data.Datas)+startIndex), style)
	}
	//header
	if style, err := xlsx.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","pattern":1,"color":["#dddddd"]},"border":[{"type":"left","color":"666666","style":1},{"type":"top","color":"666666","style":1},{"type":"bottom","color":"666666","style":1},{"type":"right","color":"666666","style":1}]}`); err == nil {
		xlsx.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", columns[len(columns)-1].excelName, startIndex), style)
	}
	if style, err := xlsx.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","pattern":1,"color":["#dddddd"]},"border":[{"type":"left","color":"666666","style":1},{"type":"top","color":"666666","style":1},{"type":"bottom","color":"666666","style":1},{"type":"right","color":"666666","style":1}]}`); err == nil {
		xlsx.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", "A", len(data.Datas)+startIndex), style)
	}

	for r, row := range data.Datas {
		for _, cell := range row {
			if c, ok := colMap[cell.Name]; ok {
				xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", c.excelName, r+startIndex+1), cell.Value)
			}
		}
	}
	fileData := FileData{Dir: data.Dir, FileName: data.FileName}
	if fileData.FileName == "" {
		fileData.FileName = fmt.Sprintf("%s.%s", utils.GUID(), "xlsx")
	}
	if fileData.Dir == "" {
		fileData.Dir = path.Join(utils.DefaultConfig.App.Storage, "exports", utils.TimeNow().Format("200601"))
	}
	utils.CreatePath(fileData.Dir)
	fileData.FullPath = utils.JoinCurrentPath(path.Join(fileData.Dir, fileData.FileName))
	if err := xlsx.SaveAs(fileData.FullPath); err != nil {
		glog.Error(err)
		return nil, err
	}
	return &fileData, nil
}
