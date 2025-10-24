package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	mgr := NewTokenMgr("123456")

	token, err := mgr.Gen("archimetagame", "test", "player", time.Second*30, jwt.SigningMethodHS256)

	if err != nil {
		fmt.Println("token gen err.", err)
		return
	}

	fmt.Println("token: ", token)

	expireAt, err := mgr.GetExpireAt(token)
	if err != nil {
		fmt.Println("token get expire at err.", err)
		return
	}
	fmt.Println("token expireAt: ", expireAt)

	ok, err := mgr.Verify(token)
	if err != nil {
		fmt.Println("token verify err.", err)
		return
	}

	fmt.Println("token verify result:", ok)
}
