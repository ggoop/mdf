package token

import (
	"github.com/ggoop/mdf/gin"
	"github.com/ggoop/mdf/utils"
)

const (
	AuthSessionKey    = "GGOOPAUTH"
	AuthTokenKey      = "token"
	DefaultContextKey = "context"
)

type Config struct {
}

func Default() gin.HandlerFunc {
	config := Config{}
	return New(config)
}
func New(config Config) gin.HandlerFunc {
	return newHandler(config).Handle
}

// Get returns the request identifier
func Get(c *gin.Context) *utils.TokenContext {
	if v, ok := c.Get(DefaultContextKey); ok {
		if vv, is := v.(*utils.TokenContext); is {
			return vv
		}
	}
	return utils.NewTokenContext()
}
