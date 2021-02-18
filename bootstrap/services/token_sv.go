package services

import (
	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/utils"
)

type ITokenSv interface {
}

type tokenSvImpl struct{}

var tokenSv ITokenSv = newTokenSvImpl()

func TokenSv() ITokenSv {
	return tokenSv
}

/**
* 创建服务实例
 */
func newTokenSvImpl() *tokenSvImpl {
	return &tokenSvImpl{}
}

func (s *tokenSvImpl) Create(token *model.Token) *model.Token {
	token.ID = utils.GUID()
	token.Token = utils.GUID()
	db.Default().Create(token)
	return token
}
func (s *tokenSvImpl) Get(token string) (*model.Token, bool) {
	t := &model.Token{}
	if err := db.Default().Model(&model.Token{}).Where("token=?", token).Take(t).Error; err != nil {
		return t, false
	}
	return t, true
}

func (s *tokenSvImpl) Delete(ids []string) error {
	if err := db.Default().Delete(model.Token{}, "id in (?)", ids).Error; err != nil {
		return err
	}
	return nil
}
func (s *tokenSvImpl) GetAndUse(token string) (*model.Token, bool) {
	t := &model.Token{}
	if err := db.Default().Model(&model.Token{}).Where("token=?", token).Take(t).Error; err != nil {
		return t, false
	}
	db.Default().Where("id = ?", t.ID).Delete(&model.Token{})
	return t, true
}
