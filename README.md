# mt-multiserver-chatcommands
mt-multiserver-chatcommands provides a useful chat command interface for mt-multiserver-proxy.

## Commands

> `shutdown`
```
Usage: `shutdown`
Perm: cmd_shutdown
Description: Disconnect all clients and stop the server.
```

> `find`
```
Usage: `find <name>`
Perm: cmd_find
Description: Check whether a player is connected and report their upstream server if they are.
```

> `addr`
```
Usage: `addr <name>`
Perm: cmd_addr
Description: Find the network address of a player if they're connected.
```

> `alert`
```
Usage: `alert <message>`
Perm: cmd_alert
Description: Send a message to all connected clients regardless of their upstream server.
```


> `send`
```
Usage: `send <player <server> <name> | current <server> | all <server>>`
Perm: cmd_send
Description: Send player(s) to a new server. player causes a single player to be redirected, current affects all players that are on your current server and all affects everyone.
```

> `players`
```
Usage: `players`
Perm: cmd_players
Description: Show the player list of every server.
```

> `reload`
```
Usage: `reload`
Perm: cmd_reload
Description: Reload the configuration file. You should restart the proxy instead if possible.
```

> `group`
```
Usage: `group [name]`
Perm: cmd_group
Description: Display the group of a player. Display your group if no player name is specified.
```

> `perms`
```
Usage: `perms [name]`
Perm: cmd_perms
Description: Show the permissions of a player. Show your permissions if no player name is specified.
```

> `gperms`
```
Usage: `gperms <group>`
Perm: cmd_gperms
Description: Show the permissions of a group.
```

> `server`
```
Usage: `server [server]`
Perm: cmd_server
Description: Display your current upstream server and all other configured servers. If a valid server name is specified, switch to that server.
```

> `kick`
```
Usage: `kick <name> [reason]`
Perm: cmd_kick
Description: Disconnect a player with an optional reason.
```

> `ban`
```
Usage: `ban <name>`
Perm: cmd_ban
Description: Ban a player from using the proxy.
```

> `unban`
```
Usage: `unban <name | address>`
Perm: cmd_unban
Description: Remove a player from the ban list. Accepts addresses and names.
```

> `uptime`
```
Usage: `uptime`
Perm: cmd_uptime
Description: Show the uptime of the proxy.
```

> `help`
```
Usage: `help [command]`
Perm: cmd_help
Description: Show help for a command (all commands if unspecified).
```

> `usage`
```
Usage: `usage [command]`
Perm: cmd_usage
Description: Show the usage string of a command (all commands if unspecified).
```
