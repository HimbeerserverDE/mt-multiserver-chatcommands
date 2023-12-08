/*
mt-multiserver-chatcommands provides a useful chat command interface
for mt-multiserver-proxy.
*/
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	proxy "github.com/HimbeerserverDE/mt-multiserver-proxy"
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

			return "Player is not connected."
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

			return "Player is not connected."
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
			usage := func() string {
				if cc != nil {
					return "Usage: send <player <server> <name> | current <server> | all <server>>"
				} else {
					return "Usage: send <player <server> <name> | all <server>>"
				}
			}

			if len(args) < 2 {
				return usage()
			}

			if _, ok := proxy.Conf().Servers[args[1]]; !ok {
				return "Server does not exist."
			}

			switch args[0] {
			case "player":
				if len(args) != 3 {
					return usage()
				}

				clt := proxy.Find(args[2])
				if clt == nil {
					return "Player is not connected."
				}

				if args[1] == clt.ServerName() {
					return "Player is already connected to this server."
				}

				if err := clt.Hop(args[1]); err != nil {
					clt.Log("<-", err)

					if errors.Is(err, proxy.ErrNoSuchServer) {
						return "Server does not exist."
					} else if errors.Is(err, proxy.ErrNewMediaPool) {
						return "The new server belongs to a media pool that is not present on this client."
					}

					clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
					return "Could not switch servers. Error: " + err.Error()
				}
			case "current":
				if cc == nil {
					return usage()
				}

				if len(args) != 2 {
					return usage()
				}

				if args[1] == cc.ServerName() {
					return "Start and destination are identical."
				}

				for clt := range proxy.Clts() {
					if clt.ServerName() == cc.ServerName() && clt.ServerName() != args[1] {
						if err := clt.Hop(args[1]); err != nil {
							clt.Log("<-", err)

							if errors.Is(err, proxy.ErrNoSuchServer) {
								return "Server does not exist."
							}

							clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
						}
					}
				}
			case "all":
				if len(args) != 2 {
					return usage()
				}

				for clt := range proxy.Clts() {
					if clt.ServerName() != args[1] {
						if err := clt.Hop(args[1]); err != nil {
							clt.Log("<-", err)

							if errors.Is(err, proxy.ErrNoSuchServer) {
								return "Server does not exist."
							}

							clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
						}
					}
				}
			default:
				return usage()
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "gsend",
		Perm:  "cmd_gsend",
		Help:  "Send player(s) to a new server group. player causes a single player to be redirected, current affects all players that are on your current server and all affects everyone.",
		Usage: "gsend <player <group> <name> | current <server> | all <server>>",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			usage := "Usage: gsend <player <group> <name> | current <server> | all <server>>"

			if len(args) < 2 {
				return usage
			}

			switch args[0] {
			case "player":
				if len(args) != 3 {
					return usage
				}

				clt := proxy.Find(args[2])
				if clt == nil {
					return "Player is not connected."
				}

				for i := 0; i < 5; i++ {
					srv, ok := proxy.Conf().RandomGroupServer(args[1])
					if !ok {
						return "Group does not exist."
					}

					if srv == args[1] {
						return "Group is also a server."
					}

					if clt.ServerName() == srv {
						if i == 4 {
							return "Player is already connected to this server after 5 attempts."
						} else {
							continue
						}
					}

					if err := clt.HopGroup(args[1]); err != nil {
						clt.Log("<-", err)

						if errors.Is(err, proxy.ErrNoSuchServer) {
							return "Server does not exist."
						} else if errors.Is(err, proxy.ErrNewMediaPool) {
							return "The new server belongs to a media pool that is not present on this client."
						}

						clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
						return "Could not switch servers. Error: " + err.Error()
					}
				}
			case "current":
				if len(args) != 2 {
					return usage
				}

				for clt := range proxy.Clts() {
					if clt.ServerName() == cc.ServerName() {
						for i := 0; i < 5; i++ {
							srv, ok := proxy.Conf().RandomGroupServer(args[1])
							if !ok {
								return "Group does not exist."
							}

							if srv == args[1] {
								return "Group is also a server."
							}

							if srv == cc.ServerName() {
								return "Start and destination are identical."
							}

							if clt.ServerName() == srv {
								continue
							}

							if err := clt.HopGroup(args[1]); err != nil {
								clt.Log("<-", err)

								if errors.Is(err, proxy.ErrNoSuchServer) {
									return "Server does not exist."
								}

								clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
							}
						}
					}
				}
			case "all":
				if len(args) != 2 {
					return usage
				}

				for clt := range proxy.Clts() {
					for i := 0; i < 5; i++ {
						srv, ok := proxy.Conf().RandomGroupServer(args[1])
						if !ok {
							return "Group does not exist."
						}

						if srv == args[1] {
							return "Group is also a server."
						}

						if clt.ServerName() == srv {
							continue
						}

						if err := clt.HopGroup(args[1]); err != nil {
							clt.Log("<-", err)

							if errors.Is(err, proxy.ErrNoSuchServer) {
								return "Server does not exist."
							}

							clt.SendChatMsg("Could not switch servers. Reconnect if you encounter any problems. Error:", err.Error())
						}
					}
				}
			default:
				return usage
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
			for clt := range proxy.Clts() {
				srv := clt.ServerName()
				srvs[srv] = append(srvs[srv], clt)
			}

			b := &strings.Builder{}
			for srv, clts := range srvs {
				if srv != "" {
					b.WriteString(srv + ":\n")
				} else {
					b.WriteString("--- No server ---\n")
				}

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
		Name:  "group",
		Perm:  "cmd_group",
		Help:  "Display the group of a player. Display your group if no player name is specified.",
		Usage: "group [name]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) > 0 {
				if len(args) != 1 {
					return "Usage: group [name]"
				}

				grp, ok := proxy.Conf().UserGroups[args[0]]
				if !ok {
					grp = "default"
				}

				return "Group: " + grp
			}

			grp, ok := proxy.Conf().UserGroups[cc.Name()]
			if !ok {
				grp = "default"
			}

			return "Your group: " + grp
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "perms",
		Perm:  "cmd_perms",
		Help:  "Show the permissions of a player. Show your permissions if no player name is specified.",
		Usage: "perms [name]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) > 0 {
				if len(args) != 1 {
					return "Usage: perms [name]"
				}

				clt := proxy.Find(args[0])
				if clt == nil {
					return "Player is not connected."
				}

				return "Player permissions: " + strings.Join(clt.Perms(), ", ")
			}

			return "Your permissions: " + strings.Join(cc.Perms(), ", ")
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "gperms",
		Perm:  "cmd_gperms",
		Help:  "Show the permissions of a group.",
		Usage: "gperms [group]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) > 0 {
				if len(args) != 1 {
					return "Usage: gperms [group]"
				}

				perms, ok := proxy.Conf().Groups[args[0]]
				if !ok {
					return "Group does not exist."
				}
				return "Group permissions: " + strings.Join(perms, ", ")
			}

			grp, ok := proxy.Conf().UserGroups[cc.Name()]
			if !ok {
				grp = "default"
			}

			perms, _ := proxy.Conf().Groups[grp]
			return "Group permissions: " + strings.Join(perms, ", ")
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "server",
		Perm:  "cmd_server",
		Help:  "Display your current upstream server and all other configured servers. If a valid server name is specified, switch to that server.",
		Usage: "server [server]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 1 {
				if len(args) > 1 {
					return "Usage: server [server]"
				}

				srvs := make([]string, 0, len(proxy.Conf().Servers))
				for name := range proxy.Conf().Servers {
					srvs = append(srvs, name)
				}

				return fmt.Sprintf("Connected to: %s | Servers: %s", cc.ServerName(), strings.Join(srvs, ", "))
			}

			if cc.ServerName() == args[0] {
				return "Already connected to this server."
			}

			if err := cc.Hop(args[0]); err != nil {
				cc.Log("<-", err)

				if errors.Is(err, proxy.ErrNoSuchServer) {
					return "Server does not exist."
				} else if errors.Is(err, proxy.ErrNewMediaPool) {
					return "The new server belongs to a media pool that is not present on this client. Please reconnect to access it."
				}

				return ""
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "gserver",
		Perm:  "cmd_gserver",
		Help:  "Display the groups your current upstream server is in and all other configured groups. If a valid group name is specified, switch to a random server of that group.",
		Usage: "gserver [group]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 1 {
				if len(args) > 1 {
					return "Usage: gserver [group]"
				}

				conf := proxy.Conf()

				srv, ok := conf.Servers[cc.ServerName()]
				if !ok {
					return "Not connected to a server."
				}

				groups := conf.ServerGroups()
				srvGroups := strings.Join(srv.Groups, ", ")

				allGroups := make([]string, 0, len(groups))
				for group := range groups {
					allGroups = append(allGroups, group)
				}

				return fmt.Sprintf("Connected to: %s | Groups: %s", srvGroups, strings.Join(allGroups, ", "))
			}

			for i := 0; i < 5; i++ {
				srv, ok := proxy.Conf().RandomGroupServer(args[0])
				if !ok {
					return "Group does not exist."
				}

				if srv == args[0] {
					return "Group is also a server."
				}

				if cc.ServerName() == srv {
					if i == 4 {
						return "Already connected to this server after 5 attempts."
					} else {
						continue
					}
				}

				if err := cc.HopGroup(srv); err != nil {
					cc.Log("<-", err)

					if errors.Is(err, proxy.ErrNoSuchServer) {
						return "Server does not exist."
					} else if errors.Is(err, proxy.ErrNewMediaPool) {
						return "The new server belongs to a media pool that is not present on this client. Please reconnect to access it."
					}

					return ""
				}
			}

			return ""
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "kick",
		Perm:  "cmd_kick",
		Help:  "Disconnect a player with an optional reason.",
		Usage: "kick <name> [reason]",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) < 1 {
				return "Usage: kick <name> [reason]"
			}

			reason := "Kicked by proxy. "
			if len(args) >= 2 {
				reason += strings.Join(args[1:], " ")
			}
			reason = strings.Trim(reason, " ")

			clt := proxy.Find(args[0])
			if clt == nil {
				return "Player is not connected."
			}

			clt.Kick(reason)
			return "Player kicked."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "ban",
		Perm:  "cmd_ban",
		Help:  "Ban a player from using the proxy.",
		Usage: "ban <name>",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 1 {
				return "Usage: ban <name>"
			}

			clt := proxy.Find(args[0])
			if clt == nil {
				return "Player is not connected."
			}

			if err := clt.Ban(); err != nil {
				return "Could not ban. Error: " + err.Error()
			}

			return "Player banned."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "unban",
		Perm:  "cmd_unban",
		Help:  "Remove a player from the ban list. Accepts addresses and names.",
		Usage: "unban <name | address>",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			if len(args) != 1 {
				return "Usage: unban <name | address>"
			}

			if err := proxy.Unban(args[0]); err != nil {
				return "Could not unban. Error: " + err.Error()
			}

			return "Player unbanned."
		},
	})
	proxy.RegisterChatCmd(proxy.ChatCmd{
		Name:  "uptime",
		Perm:  "cmd_uptime",
		Help:  "Show the uptime of the proxy.",
		Usage: "uptime",
		Handler: func(cc *proxy.ClientConn, args ...string) string {
			return fmt.Sprintf("Uptime: %fs", proxy.Uptime().Seconds())
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
				return proxy.Colorize(name+": ", "#6FF") + cmds[name].Help
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
				return proxy.Colorize(name+": ", "#6F3") + cmds[name].Usage
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
