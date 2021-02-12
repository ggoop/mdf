package services

import (
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/utils"
)

type TokenSv struct {
	repo *repositories.MysqlRepo
}

/**
* 创建服务实例
 */
func NewTokenSv(repo *repositories.MysqlRepo) *TokenSv {
	return &TokenSv{repo: repo}
}

func (s *TokenSv) Create(token *model.Token) *model.Token {
	token.ID = utils.GUID()
	token.Token = utils.GUID()
	s.repo.Create(token)
	return token
}
func (s *TokenSv) Get(token string) (*model.Token, bool) {
	t := &model.Token{}
	if err := s.repo.Model(&model.Token{}).Where("token=?", token).Take(t).Error; err != nil {
		return t, false
	}
	return t, true
}

func (s *TokenSv) Delete(ids []string) error {
	if err := s.repo.Delete(model.Token{}, "id in (?)", ids).Error; err != nil {
		return err
	}
	return nil
}
func (s *TokenSv) GetAndUse(token string) (*model.Token, bool) {
	t := &model.Token{}
	if err := s.repo.Model(&model.Token{}).Where("token=?", token).Take(t).Error; err != nil {
		return t, false
	}
	s.repo.Where("id = ?", t.ID).Delete(&model.Token{})
	return t, true
}
