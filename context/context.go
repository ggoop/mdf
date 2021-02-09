package context

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ggoop/mdf/configs"
	"github.com/ggoop/mdf/utils"

	"github.com/dgrijalva/jwt-go"
)

func New() Context {
	return Context{data: make(map[string]string)}
}

type Token struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (s *Token) String() string {
	return fmt.Sprintf("%s %s", s.Type, s.AccessToken)
}

type Context struct {
	data map[string]string
}
type Config struct {
	Expiration bool
	Credential bool
	User       bool
	Ent        bool
}

const (
	AuthSessionKey    = "GGOOPAUTH"
	DefaultContextKey = "context"
)

func (s *Context) Copy() Context {
	c := New()
	if s.data != nil {
		for k, v := range s.data {
			c.data[k] = v
		}
	}
	return c
}
func (s *Context) ValueReplace(value string) string {
	if s.data != nil {
		const REGEXP_FIELD_EXP string = `{([A-Za-z._]+[0-9A-Za-z]*)}`
		r, _ := regexp.Compile(REGEXP_FIELD_EXP)
		matched := r.FindAllStringSubmatch(value, -1)
		for _, match := range matched {
			v := s.GetValue(utils.SnakeString(match[1]))
			value = strings.ReplaceAll(value, match[0], v)
		}
	}
	return value
}
func (s *Context) SetID(id string) {
	s.SetValue("id", id)
}
func (s *Context) ID() string {
	return s.GetValue("id")
}

func (s *Context) SetEntID(ent string) {
	s.SetValue("ent_id", ent)
}
func (s *Context) SetClientID(ent string) {
	s.SetValue("client_id", ent)
}
func (s *Context) SetUserID(user string) {
	s.SetValue("user_id", user)
}
func (s *Context) EntID() string {
	return s.GetValue("ent_id")
}
func (s *Context) ClientID() string {
	return s.GetValue("client_id")
}
func (s *Context) UserID() string {
	return s.GetValue("user_id")
}
func (s *Context) ExpiresAt() int64 {
	str := s.GetValue("expires_at")
	exp, _ := strconv.ParseInt(str, 10, 64)
	return exp
}

// 验证
func (s *Context) Valid(user, ent bool) error {
	if user && s.VerifyUser() != nil {
		return utils.NewError(fmt.Errorf("用户验证失败"), 401)
	}
	if ent && s.VerifyEnt() != nil {
		return utils.NewError(fmt.Errorf("企业验证失败"), 401)
	}
	return nil
}
func (c *Context) VerifyClient() error {
	if c.ClientID() == "" {
		return utils.NewError(fmt.Errorf("客户端验证失败"), 401)
	}
	return nil
}
func (c *Context) VerifyUser() error {
	if c.UserID() == "" {
		return utils.NewError(fmt.Errorf("用户验证失败"), 401)
	}
	return nil
}
func (c *Context) VerifyEnt() error {
	if c.EntID() == "" {
		return utils.NewError(fmt.Errorf("企业验证失败"), 401)
	}
	return nil
}
func (c *Context) VerifyExpiresAt(cmp int64, required bool) bool {
	exp := c.ExpiresAt()
	if exp == 0 {
		return !required
	}
	return cmp <= exp
}

//
func (s *Context) Clean() {
	s.data = make(map[string]string)
}
func (s *Context) GetValue(name string) string {
	if s.data == nil {
		s.data = make(map[string]string)
		return ""
	}
	name = utils.SnakeString(name)
	return s.data[name]
}
func (s *Context) GetIntValue(name string) int {
	v, _ := strconv.Atoi(s.GetValue(name))
	return v
}
func (s *Context) GetInt64Value(name string) int64 {
	v, _ := strconv.ParseInt(s.GetValue(name), 10, 64)
	return v
}
func (s *Context) SetValue(name, value string) *Context {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	name = utils.SnakeString(name)
	if name == "ent_id" {
		s.data[name] = value
	} else if name == "user_id" {
		s.data[name] = value
	} else {
		s.data[name] = value
	}
	return s
}
func (s *Context) ToToken() *Token {
	tokenStr := strings.Split(s.ToTokenString(), " ")
	if len(tokenStr) < 2 {
		return nil
	}
	t := Token{AccessToken: tokenStr[1], Type: tokenStr[0], ExpiresAt: s.ExpiresAt()}
	return &t
}
func (s *Context) ToTokenString() string {
	claim := jwt.MapClaims{}
	for k, v := range s.data {
		claim[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(configs.Default.App.Token))
	if err != nil {
		return ""
	}
	return "bearer " + tokenString
}
func (s *Context) FromTokenString(token string) (*Context, error) {
	ctx := New()
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
		token = tokenParts[1]
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Default.App.Token), nil
	})
	if err != nil {
		return &ctx, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		for k, v := range claims {
			if vstr, ok := v.(string); ok {
				ctx.SetValue(k, vstr)
			}
		}
	}
	return &ctx, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Context) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" || string(bytes) == "" {
		return nil
	}
	data := make(map[string]string)
	json.Unmarshal(bytes, &data)
	nc := Context{data: data}
	*d = nc
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Context) MarshalJSON() ([]byte, error) {
	if d.data == nil || len(d.data) == 0 {
		return []byte(""), nil
	}
	bytes, err := json.Marshal(d.data)
	return bytes, err
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Context) Scan(value interface{}) error {
	if value == nil {
		*d = New()
		return nil
	}
	switch v := value.(type) {
	case float32:
		*d = New()
		return nil

	case float64:
		*d = New()
		return nil

	case int64:
		*d = New()
		return nil
	case string:
		nd := New()
		json.Unmarshal([]byte(v), &nd)
		*d = nd
		return nil
	default:
		*d = New()
		return nil
	}
}

// Value implements the driver.Valuer interface for database serialization.
func (d Context) Value() (driver.Value, error) {
	data, err := d.MarshalJSON()
	return string(data), err
}

func (t Context) IsValid() bool {
	return t.data != nil && len(t.data) > 0
}
