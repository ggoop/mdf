package utils

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
)

/**
请求通用类
*/
type TokenContext struct {
	*BizObject
}

func NewTokenContext() *TokenContext {
	return &TokenContext{BizObject: NewBizObject()}
}
func (s *TokenContext) UserID() string {
	return s.GetString("UserID")
}

func (s *TokenContext) EntID() string {
	return s.GetString("EntID")
}

func (s *TokenContext) OrgID() string {
	return s.GetString("OrgID")
}

func (s *TokenContext) ID() string {
	return s.GetString("ID")
}

func (s TokenContext) ToTokenString() string {
	claim := jwt.MapClaims{}
	data := s.Data()
	for k, v := range data {
		claim[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(DefaultConfig.App.Token))
	if err != nil {
		return ""
	}
	return "bearer " + tokenString
}
func (s TokenContext) FromTokenString(token string) (*TokenContext, error) {
	ctx := NewTokenContext()
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
		token = tokenParts[1]
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(DefaultConfig.App.Token), nil
	})
	if err != nil {
		return ctx, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		for k, v := range claims {
			ctx.Set(k, v)
		}
		ctx.Set("id", GUID())
	}
	return ctx, nil
}
