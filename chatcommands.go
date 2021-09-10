/*
mt-multiserver-chatcommands provides a useful chat command interface
for mt-multiserver-proxy.
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HimbeerserverDE/mt-multiserver-proxy"
)

func init() {
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "shutdown",
		Perm:  "cmd_shutdown",
		Help:  "Disconnect all clients and stop the server.",
		Usage: "shutdown",
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
		Name:  "find",
		Perm:  "cmd_find",
		Help:  "Check whether a player is connected and report their upstream server if they are.",
		Usage: "find <name>",
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
		Name:  "addr",
		Perm:  "cmd_addr",
		Help:  "Find the network address of a player if they're connected.",
		Usage: "addr <name>",
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
		Name:  "alert",
		Perm:  "cmd_alert",
		Help:  "Send a message to all connected clients regardless of their upstream server.",
		Usage: "alert <message>",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) <= 0 {
				return "Usage: alert <message>"
			}

			msg := "[ALERT] "
			msg += strings.Join(args, " ")

			for clt := range proxy.Clts() {
				clt.SendChatMsg(msg)
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "send",
		Perm:  "cmd_send",
		Help:  "Send player(s) to a new server. player causes a single player to be redirected, current affects all players that are on your current server and all affects everyone.",
		Usage: "send <player <server> <name> | current <server> | all <server>>",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) < 2 {
				return "Usage: send <player <server> <name> | current <server> | all <server>>"
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
				if len(args) != 3 {
					return "Usage: send <player <server> <name> | current <server> | all <server>>"
				}

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
				if len(args) != 2 {
					return "Usage: send <player <server> <name> | current <server> | all <server>>"
				}

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
				if len(args) != 2 {
					return "Usage: send <player <server> <name> | current <server> | all <server>>"
				}

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
		Name:  "players",
		Perm:  "cmd_players",
		Help:  "Show the player list of every server.",
		Usage: "players",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			srvs := make(map[string][]*proxy.ClientConn)
			for _, srv := range proxy.Conf().Servers {
				srvs[srv.Name] = []*proxy.ClientConn{}
			}

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
		Name:  "reload",
		Perm:  "cmd_reload",
		Help:  "Reload the configuration file. You should restart the proxy instead if possible.",
		Usage: "reload",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if err := proxy.LoadConfig(); err != nil {
				return "Configuration could not be reloaded. Old config is still active. Error: " + err.Error()
			}

			return "Configuration reloaded. You should restart if you encounter any problems."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "perms",
		Perm:  "cmd_perms",
		Help:  "Show the permissions of a player. Show your permissions if no player name is specified.",
		Usage: "perms [name]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			return "Your permissions: " + strings.Join(cc.Perms(), ", ")
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "server",
		Perm:  "cmd_server",
		Help:  "Display your current upstream server and all other configured servers. If a valid server name is specified, switch to that server.",
		Usage: "server [server]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) == 0 {
				srvs := make([]string, len(proxy.Conf().Servers))
				for i, srv := range proxy.Conf().Servers {
					srvs[i] = srv.Name
				}

				return fmt.Sprintf("Connected to: %s | Servers: %s", cc.ServerName(), strings.Join(srvs, ", "))
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

	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "help",
		Perm:  "cmd_help",
		Help:  "Show help for a command (all commands if unspecified).",
		Usage: "help [command]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			cmds := proxy.ChatCmds()

			help := func(name string) string {
				str := proxy.Colorize(name+": ", "#6FF")
				str += cmds[name].Help

				return str
			}

			if len(args) != 1 {
				if len(args) > 1 {
					return "Usage: help [command]"
				}

				for cmd := range cmds {
					cc.SendChatMsg(help(cmd))
				}
			} else {
				if _, ok := cmds[args[0]]; !ok {
					return "Inexistent command."
				}

				return help(args[0])
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "usage",
		Perm:  "cmd_usage",
		Help:  "Show the usage string of a command (all commands if unspecified).",
		Usage: "usage [command]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			cmds := proxy.ChatCmds()

			usage := func(name string) string {
				str := proxy.Colorize(name+": ", "#6F3")
				str += cmds[name].Usage

				return str
			}

			if len(args) != 1 {
				if len(args) > 1 {
					return "Usage: usage [command]"
				}

				for cmd := range cmds {
					cc.SendChatMsg(usage(cmd))
				}
			} else {
				if _, ok := cmds[args[0]]; !ok {
					return "Inexistent command."
				}

				return usage(args[0])
			}

			return ""
		},
	})
}
