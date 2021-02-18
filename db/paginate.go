package db

import (
	"github.com/ggoop/mdf/db/gorm"
	"github.com/ggoop/mdf/utils"
)

func Paginate(db *gorm.DB, out interface{}, page int, pageSize int) (*utils.Pager, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 30
	}
	rtn := utils.Pager{}
	totals := 0
	if pageSize > 0 {
		if err := db.Count(&totals).Error; err != nil {
			return nil, err
		}
		if totals > 0 {
			if err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(out).Error; err != nil {
				return nil, err
			}
		}
	} else {
		if err := db.Find(out).Error; err != nil {
			return nil, err
		}
	}
	rtn.Value = out
	rtn.Page = page
	rtn.PageSize = pageSize
	rtn.Total = totals
	return &rtn, nil
}
func PaginateScan(db *gorm.DB, out interface{}, page int, pageSize int) (*utils.Pager, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 30
	}
	rtn := utils.Pager{}
	totals := 0
	if pageSize > 0 {
		if err := db.Count(&totals).Error; err != nil {
			return nil, err
		}
		if totals > 0 {
			if err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(out).Error; err != nil {
				return nil, err
			}
		}
	} else {
		if err := db.Scan(out).Error; err != nil {
			return nil, err
		}
	}
	rtn.Value = out
	rtn.Page = page
	rtn.PageSize = pageSize
	rtn.Total = totals
	return &rtn, nil
}
