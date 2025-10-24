package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	mgr := NewTokenMgr("123456")

	token := mgr.Gen("archimeta", "test", "test", time.Second*30, jwt.SigningMethodHS256)

	fmt.Println("token: ", token)

	ok, err := mgr.Verify(token)
	if err != nil {
		fmt.Println("token verify err.", err)
		return
	}

	fmt.Println("token verify result:", ok)
}
