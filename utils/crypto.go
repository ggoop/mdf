package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/ggoop/mdf/framework/glog"
)

func getKey(key string) []byte {
	//AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	l := 32
	if len(key) >= 32 {
		l = 32
	} else if len(key) >= 24 {
		l = 24
	} else {
		l = 16
	}
	ctx := md5.New()
	ctx.Write([]byte(key))
	return ([]byte(hex.EncodeToString(ctx.Sum(nil))))[:l]
}
func AesCFBEncrypt(text, key string) (string, error) {
	skey := getKey(key)
	var iv = skey[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(skey)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}
func AesCFBDecrypt(encrypted, key string) (string, error) {
	skey := getKey(key)
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = skey[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher(skey)
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}

//获取签名算法为HS256的token
func CreateJWTToken(dic map[string]interface{}, SIGNED_KEY string) string {
	claims := jwt.MapClaims{}
	for k, v := range dic {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密算法是HS256时，这里的SignedString必须是[]byte（）类型
	ss, err := token.SignedString([]byte(SIGNED_KEY))
	if err != nil {
		glog.Infof("token生成签名错误,err=%v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(ss))
}

//解析签名算法为HS256的token
func ParseJWTToken(tokenString string, SIGNED_KEY string) (jwt.MapClaims, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}
	tokenString = string(decodeBytes)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNED_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("ParseHStoken:claims类型转换失败")
	}
	return claims, nil
}

//解析签名算法为HS256的token
func ParseJWTTokenNotValidation(tokenString string, SIGNED_KEY string) (jwt.MapClaims, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}
	tokenString = string(decodeBytes)

	claims := jwt.MapClaims{}
	parse := new(jwt.Parser)
	parse.SkipClaimsValidation = true
	_, _, err = parse.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

/*
  md5 sign: "123" -> "202cb962ac59075b964b07152d234b70"
*/
func Md5Signer(message string) string {
	data := []byte(message)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
