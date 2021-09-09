package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HimbeerserverDE/mt-multiserver-proxy"
)

func init() {
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "shutdown",
		Perm: "cmd_shutdown",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				return "Could not find process: " + err.Error()
			}

			proc.Signal(os.Interrupt)
			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "reload",
		Perm: "cmd_reload",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if err := proxy.LoadConfig(); err != nil {
				return "Configuration could not be reloaded. Old config is still active. Error: " + err.Error()
			}

			return "Configuration reloaded. You should restart if you encounter any problems."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "perms",
		Perm: "cmd_perms",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			return "Your permissions: " + strings.Join(cc.Perms(), " ")
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "server",
		Perm: "cmd_server",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) == 0 {
				srvs := make([]string, len(proxy.Conf().Servers))
				for i, srv := range proxy.Conf().Servers {
					srvs[i] = srv.Name
				}

				return fmt.Sprintf("Connected to: %s | Servers: %s", cc.ServerName(), strings.Join(srvs, " "))
			}

			srv := strings.Join(args, " ")
			if cc.ServerName() == srv {
				return "Already connected to this server."
			}

			if err := cc.Hop(srv); err != nil {
				return "Could not switch servers. Reconnect if you encounter any problems. Error: " + err.Error()
			}

			return ""
		},
	})
}
