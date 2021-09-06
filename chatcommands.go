package main

import (
	"strings"

	"github.com/HimbeerserverDE/mt-multiserver-proxy"
)

func init() {
	testCmd := func(cc *proxy.ClientConn, args ...string) string {
		return strings.Join(args, " ")
	}
	proxy.RegisterChatCmd("test", testCmd)
}
