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
		Name: "find",
		Perm: "cmd_find",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 1 {
				return "Usage: find <name>"
			}

			clt := proxy.Find(args[0])
			if clt != nil {
				return fmt.Sprintf("%s is connected to %s", clt.Name(), clt.ServerName())
			}

			return "Player not connected."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "addr",
		Perm: "cmd_addr",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 1 {
				return "Usage: addr <name>"
			}

			clt := proxy.Find(args[0])
			if clt != nil {
				return fmt.Sprintf("%s is at %s", clt.Name(), clt.RemoteAddr())
			}

			return "Player not connected."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "alert",
		Perm: "cmd_alert",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			msg := strings.Join(args, " ")
			for clt := range proxy.Clts() {
				clt.SendChatMsg(msg)
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "send",
		Perm: "cmd_send",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 2 {
				return "Usage: send <player | current | all> <server> [name]"
			}

			var found bool
			for _, srv := range proxy.Conf().Servers {
				if srv.Name == args[1] {
					found = true
					break
				}
			}

			if !found {
				return "Server not existent."
			}

			switch args[0] {
			case "player":
				if args[1] == cc.ServerName() {
					return "Player is already connected to this server."
				}

				clt := proxy.Find(args[2])
				if clt == nil {
					return "Player not connected."
				}

				if err := clt.Hop(args[1]); err != nil {
					clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
				}
			case "current":
				if args[1] == cc.ServerName() {
					return "Start and destination are identical."
				}

				for clt := range proxy.Clts() {
					if clt.ServerName() == cc.ServerName() && clt.ServerName() != args[1] {
						if err := clt.Hop(args[1]); err != nil {
							clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
						}
					}
				}
			case "all":
				for clt := range proxy.Clts() {
					if clt.ServerName() != args[1] {
						if err := clt.Hop(args[1]); err != nil {
							clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
						}
					}
				}
			default:
				return "Usage: send <player | current | all> <server> [name]"
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name: "players",
		Perm: "cmd_players",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			srvs := make(map[string][]*proxy.ClientConn)
			for clt := range proxy.Clts() {
				srv := clt.ServerName()
				srvs[srv] = append(srvs[srv], clt)
			}

			b := &strings.Builder{}
			for srv, clts := range srvs {
				b.WriteString(srv + ":\n")
				for _, clt := range clts {
					b.WriteString("- " + clt.Name() + "\n")
				}
			}

			return strings.Trim(b.String(), "\n")
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
