package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenMgr struct {
	signingKey []byte
}

func NewTokenMgr(signingKey string) *TokenMgr {
	m := &TokenMgr{
		signingKey: []byte(signingKey),
	}

	return m
}

// Gen 生成token
//
// parameters:
//
// iss: 发行者
// sub: 主题
// aud: 受众
// expireAfter: 该时间之后过期
// method: 使用的方法，例如：jwt.SigningMethodHS256
func (m *TokenMgr) Gen(iss, sub, aud string, expireAfter time.Duration, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, jwt.MapClaims{
		"iss": iss,                                // 发行者
		"sub": sub,                                // 主题
		"aud": aud,                                // 受众
		"exp": time.Now().Add(expireAfter).Unix(), // 过期时间，例如72小时后过期
		"iat": time.Now().Unix(),                  // 签发时间
	})

	// 生成签名的token
	return token.SignedString(m.signingKey)
}

func (m *TokenMgr) Verify(tokenStr string) (bool, error) {
	// 解析token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 这里我们验证签名的算法是否是我们所期望的算法，这里是HS256算法，并且验证密钥是否正确。
		return m.signingKey, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New(fmt.Sprintf("token not valid. token:'%s'", tokenStr))
	}

	return true, nil
}

func (m *TokenMgr) GetClaims(tokenStr string) (jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 这里我们验证签名的算法是否是我们所期望的算法，这里是HS256算法，并且验证密钥是否正确。
		return m.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New(fmt.Sprintf("token not valid. token:'%s'", tokenStr))
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New(fmt.Sprintf("get claims failed. token:'%s'", tokenStr))
	}

	return claims, nil
}

func (m *TokenMgr) GetExpireAt(tokenStr string) (int64, error) {
	claims, err := m.GetClaims(tokenStr)

	if err != nil {
		return 0, err
	}

	exp, ok := claims["exp"]

	if !ok {
		return 0, errors.New(fmt.Sprintf("find exp from claims err.token:'%s'", tokenStr))
	}

	return int64(exp.(float64)), nil
}
