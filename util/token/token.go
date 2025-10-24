package token

import (
	"github.com/LeeroyLin/goengine/core/elog"
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
func (m *TokenMgr) Gen(iss, sub, aud string, expireAfter time.Duration, method jwt.SigningMethod) string {
	token := jwt.NewWithClaims(method, jwt.MapClaims{
		"iss": iss,                                // 发行者
		"sub": sub,                                // 主题
		"aud": aud,                                // 受众
		"exp": time.Now().Add(expireAfter).Unix(), // 过期时间，例如72小时后过期
		"iat": time.Now().Unix(),                  // 签发时间
	})

	// 生成签名的token
	tokenStr, err := token.SignedString(m.signingKey)
	if err != nil {
		elog.Error("[Token] gen token err.", err)
	}

	return tokenStr
}

func (m *TokenMgr) Verify(tokenStr string) (bool, error) {
	// 解析token
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 这里我们验证签名的算法是否是我们所期望的算法，这里是HS256算法，并且验证密钥是否正确。
		return m.signingKey, nil
	})

	if err != nil {
		return false, err
	}
	//
	//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	//	println(claims["sub"]) // 输出主题（或其他claims）的信息
	//} else {
	//	panic(err) // 如果解析失败，打印错误信息
	//}

	return true, nil
}
