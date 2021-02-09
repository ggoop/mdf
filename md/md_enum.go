package md

import (
	"strings"
	"sync"

	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"
)

/**
枚举类型
*/
type MDEnumType struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Name      string     `gorm:"size:100"`
	Domain    string     `gorm:"size:50" json:"domain"`
	Enums     []MDEnum   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:EntityID"`
}

/**
枚举值
*/
type MDEnum struct {
	EntityID  string     `gorm:"size:50;primary_key:uix;morph:limit" json:"entity_id"`
	ID        string     `gorm:"size:50;primary_key:uix" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Name      string     `gorm:"size:50" json:"name"`
	Sequence  int        `json:"sequence"`
	SrcID     string     `gorm:"size:50" json:"src_id"`
}

func (t MDEnum) TableName() string {
	return "md_enums"
}
func (s *MDEnum) MD() *Mder {
	return &Mder{ID: "md.enum", Domain: md_domain, Name: "枚举", Type: TYPE_ENUM}
}

var enumCache map[string]MDEnum

type EntitySv struct {
	repo *repositories.MysqlRepo
	*sync.Mutex
}

/**
* 创建服务实例
 */
func NewEntitySv(repo *repositories.MysqlRepo) *EntitySv {
	return &EntitySv{repo: repo, Mutex: &sync.Mutex{}}
}
func GetEnum(typeId string, values ...string) *MDEnum {
	if enumCache == nil || typeId == "" || values == nil || len(values) == 0 {
		return nil
	}
	for _, v := range values {
		if v, ok := enumCache[strings.ToLower(typeId+":"+v)]; ok {
			return &v
		}
	}
	return nil
}

func (s *EntitySv) InitCache() {
	enumCache = make(map[string]MDEnum)
	items, _ := s.GetEnums()
	for _, v := range items {
		enumCache[strings.ToLower(v.EntityID+":"+v.ID)] = v
		enumCache[strings.ToLower(v.EntityID+":"+v.Name)] = v
	}
}

func (s *EntitySv) GetEnums() ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := s.repo.Model(&MDEnum{}).Where("entity_id in (?)", s.repo.Model(MDEntity{}).Select("id").Where("type=?", "enum").SubQuery()).Order("entity_id").Order("sequence").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (s *EntitySv) GetEnumBy(typeId string) ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := s.repo.Model(&MDEnum{}).Where("entity_id=?", typeId).Order("sequence,id").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *EntitySv) UpdateEntity(item MDEntity) error {
	if item.ID == "" {
		return nil
	}
	old := MDEntity{}
	s.repo.Model(old).Where("id=?", item.ID).Take(&old)
	if old.ID != "" {
		updates := utils.Map{}
		if old.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if old.Code != item.Code && item.Code != "" {
			updates["Code"] = item.Code
		}
		if old.Domain != item.Domain && item.Domain != "" {
			updates["Domain"] = item.Domain
		}
		if old.Tags != item.Tags && item.Tags != "" {
			updates["Tags"] = item.Tags
		}
		if old.Memo != item.Memo && item.Memo != "" {
			updates["Memo"] = item.Memo
		}
		if old.Type != item.Type && item.Type != "" {
			updates["Type"] = item.Type
		}
		if old.TableName != item.TableName && item.TableName != "" {
			updates["TableName"] = item.TableName
		}
		if len(updates) > 0 {
			s.repo.Model(&old).Where("id=?", old.ID).Updates(updates)
		}
	} else {
		s.repo.Create(&item)
	}
	return nil
}
func (s *EntitySv) UpdateEnumType(enumType MDEnumType) error {
	if enumType.ID == "" {
		return nil
	}
	entity := MDEntity{}
	s.repo.Model(entity).Where("id=?", enumType.ID).Order("id").Take(&entity)
	if entity.ID == "" {
		entity.ID = enumType.ID
		entity.Code = enumType.ID
		entity.Name = enumType.Name
		entity.Type = TYPE_ENUM
		entity.Domain = enumType.Domain
		s.repo.Create(&entity)
	} else {
		updates := utils.Map{}
		if entity.Name != enumType.Name && enumType.Name != "" {
			updates["Name"] = enumType.Name
		}
		if entity.Domain != entity.Domain && enumType.Domain != "" {
			updates["Domain"] = enumType.Domain
		}
		if len(updates) > 0 {
			s.repo.Model(&entity).Updates(updates)
		}
	}
	if len(enumType.Enums) > 0 {
		for i, enum := range enumType.Enums {
			if enum.Sequence == 0 {
				enumType.Enums[i].Sequence = i
			}
			enumType.Enums[i].EntityID = entity.ID
			if _, err := s.UpdateOrCreateEnum(enumType.Enums[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *EntitySv) UpdateOrCreateEnum(enum MDEnum) (*MDEnum, error) {
	entity := MDEntity{}
	if enum.EntityID == "" {
		return nil, nil
	}
	s.repo.Model(entity).Where("id=?", enum.EntityID).Order("id").Take(&entity)
	if entity.ID == "" {
		return nil, glog.Error("找不到枚举类型！")
	}
	old := MDEnum{}
	if s.repo.Where("entity_id=? and id=?", enum.EntityID, enum.ID).Order("id").Take(&old).RecordNotFound() {
		s.repo.Create(&enum)
	} else {
		updates := utils.Map{}
		if old.Name != enum.Name && enum.Name != "" {
			updates["Name"] = enum.Name
		}
		if old.SrcID != enum.SrcID && enum.SrcID != "" {
			updates["SrcID"] = enum.SrcID
		}
		if old.Sequence != enum.Sequence && enum.Sequence >= 0 {
			updates["Sequence"] = enum.Sequence
		}
		if len(updates) > 0 {
			s.repo.Model(&old).Updates(updates)
		}
	}
	return &enum, nil
}
