package test

import (
	"fmt"
	"net"
	"testing"
)

func TestServerClient(t *testing.T) {
	_, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("connect err", err)
		return
	}

	fmt.Println("Connect success")
}
