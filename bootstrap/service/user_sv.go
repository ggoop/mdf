package service

import (
	"fmt"

	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"
)

type UserSv struct {
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewUserSv(repo *repositories.MysqlRepo) *UserSv {
	return &UserSv{repo: repo}
}

func (s *UserSv) GetUserBy(id string) (*model.User, error) {
	item := model.User{}
	if err := s.repo.Model(item).Where("id=?", id).Take(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *UserSv) UpdateUser(id string, datas map[string]interface{}) error {
	item, err := s.GetUserBy(id)
	if err != nil {
		return err
	}
	updates := make(map[string]interface{})
	for k, v := range datas {
		if k == "Name" && item.Name != v {
			updates[k] = v
		}
		if k == "AvatarUrl" && item.AvatarUrl != v {
			updates[k] = v
		}
	}
	if len(updates) > 0 {
		s.repo.Model(&model.User{}).Where("id=?", id).Updates(updates)
	}
	return nil
}

//账号
func (s *UserSv) GetUserByAccount(account string) (*model.User, error) {
	acc := model.User{}
	if err := s.repo.Model(acc).Where("email=? or mobile=? or account=?", account, account, account).Take(&acc).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}
func (s *UserSv) SetPassword(userID, password string) error {
	account := model.User{}
	if err := s.repo.Take(&account, "id=?", userID).Error; err != nil {
		return err
	}
	longPassword, _ := utils.AesCFBEncrypt(password, account.Openid)
	s.repo.Model(account).Where("id=?", account.ID).Update("Password", longPassword)
	return nil
}
func (s *UserSv) ExistsEmail(email string, excludeIds ...string) bool {
	count := 0
	q := s.repo.Model(model.User{}).Where("email=?", email)
	if excludeIds != nil && len(excludeIds) > 0 {
		q = q.Where("id not in (?)", excludeIds)
	}
	q.Count(&count)

	return count > 0
}
func (s *UserSv) ExistsMobile(mobile string, excludeIds ...string) bool {
	count := 0
	q := s.repo.Model(model.User{}).Where("mobile=?", mobile)
	if excludeIds != nil && len(excludeIds) > 0 {
		q = q.Where("id not in (?)", excludeIds)
	}
	q.Count(&count)

	return count > 0
}
func (s *UserSv) IssueUser(account *model.User) (*model.User, error) {
	oldAcc := model.User{}
	tag := false
	//如果传入ID，则先按ID查询用户
	if account.ID != "" {
		s.repo.Model(oldAcc).Where("id=?", account.ID).Take(&oldAcc)
	}
	//如果存在openid，则先按Openid查询用户
	if account.Openid != "" && oldAcc.ID == "" {
		s.repo.Model(oldAcc).Where("openid=?", account.Openid).Take(&oldAcc)
	}
	if account.Mobile != "" && oldAcc.ID == "" {
		tag = true
		s.repo.Model(oldAcc).Where("mobile=?", account.Mobile).Take(&oldAcc)
	}
	if account.Email != "" && oldAcc.ID == "" {
		tag = true
		s.repo.Model(oldAcc).Where("email=?", account.Email).Take(&oldAcc)
	}
	if account.Account != "" && oldAcc.ID == "" {
		tag = true
		s.repo.Model(oldAcc).Where("account=?", account.Account).Take(&oldAcc)
	}
	//如果是新创建，则
	if oldAcc.ID == "" && !tag {
		return nil, fmt.Errorf("mobile or email or src_id is need!")
	}
	if oldAcc.ID != "" {
		updates := make(map[string]interface{})
		if oldAcc.Name != account.Name && account.Name != "" {
			updates["Name"] = account.Name
		}
		if oldAcc.AvatarUrl != account.AvatarUrl && account.AvatarUrl != "" {
			updates["AvatarUrl"] = account.AvatarUrl
		}
		if oldAcc.Email != account.Email && account.Email != "" {
			if ex := s.ExistsEmail(account.Email, oldAcc.ID); ex {
				return nil, fmt.Errorf("电子邮件 %s 已经被注册", account.Email)
			}
			updates["Email"] = account.Email
		}
		if oldAcc.Mobile != account.Mobile && account.Mobile != "" {
			if ex := s.ExistsMobile(account.Mobile, oldAcc.ID); ex {
				return nil, fmt.Errorf("电话 %s 已经被注册", account.Mobile)
			}
			updates["Mobile"] = account.Mobile
		}
		if account.Password != "" && oldAcc.IsSystem.IsFalse() {
			updates["Password"], _ = utils.AesCFBEncrypt(account.Password, oldAcc.Openid)
		}
		if len(updates) > 0 {
			s.repo.Model(oldAcc).Where("id=?", oldAcc.ID).Updates(updates)
			s.repo.Where("id=?", oldAcc.ID).Take(&oldAcc)
		}
		account = &oldAcc
	} else {
		account.ID = utils.GUID()
		if account.Openid == "" {
			account.Openid = utils.GUID()
		}
		if account.Password != "" {
			account.Password, _ = utils.AesCFBEncrypt(account.Password, account.Openid)
		}
		if account.AvatarUrl == "" {
			account.AvatarUrl = "/img/avatar/0.jpg"
		}
		if account.Mobile != "" {
			if ex := s.ExistsMobile(account.Mobile, oldAcc.ID); ex {
				return nil, fmt.Errorf("电话 %s 已经被注册", account.Mobile)
			}
		}
		if account.Email != "" {
			if ex := s.ExistsEmail(account.Email, oldAcc.ID); ex {
				return nil, fmt.Errorf("电子邮件 %s 已经被注册", account.Email)
			}
		}
		if err := s.repo.Create(account).Error; err != nil {
			return nil, err
		}
	}
	return s.GetUserBy(account.ID)
}
