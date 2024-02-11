# mt-multiserver-chatcommands
mt-multiserver-chatcommands provides a useful chat command interface for mt-multiserver-proxy.

## Commands

> `shutdown`
```
Permission: cmd_shutdown
Description: Disconnect all clients and stop the server.
Usage: `shutdown`
```

> `find`
```
Permission: cmd_find
Description: Check whether a player is connected and report their upstream server if they are.
Usage: `find <name>`
```

> `addr`
```
Permission: cmd_addr
Description: Find the network address of a player if they're connected.
Usage: `addr <name>`
```

> `alert`
```
Permission: cmd_alert
Description: Send a message to all connected clients regardless of their upstream server.
Usage: `alert <message>`
```


> `send`
```
Permission: cmd_send
Description: Send player(s) to a new server. player causes a single player to be redirected, current affects all players that are on your current server and all affects everyone.
Usage: `send <player <server> <name> | current <server> | all <server>>`
Example: `send player lobby bob`
```

> `gsend`
```
Permission: cmd_gsend
Description: Send player(s) to a new server group. player causes a single player to be redirected, current affects all players that are on your current server and all affects everyone.
Usage: `gsend <player <group> <name> | current <group> | all <group>>`
Example: `gsend player lobbies bob`
```

> `players`
```
Permission: cmd_players
Description: Show the player list of every server.
Usage: `players`
```

> `reload`
```
Permission: cmd_reload
Description: Reload the configuration file. You should restart the proxy instead if possible.
Usage: `reload`
```

> `group`
```
Permission: cmd_group
Description: Display the group of a player. Display your group if no player name is specified.
Usage: `group [name]`
```

> `perms`
```
Permission: cmd_perms
Description: Show the permissions of a player. Show your permissions if no player name is specified.
Usage: `perms [name]`
```

> `gperms`
```
Permission: cmd_gperms
Description: Show the permissions of a group.
Usage: `gperms <group>`
```

> `server`
```
Permission: cmd_server
Description: Display your current upstream server and all other configured servers. If a valid server name is specified, switch to that server.
Usage: `server [server]`
```

> `gserver`
```
Permission: cmd_gserver
Description: Display the groups your current upstream server is in
and all other configured groups. If a valid group name is specified,
switch to a random server of that group.
Usage: `gserver [group]`
```

> `kick`
```
Permission: cmd_kick
Description: Disconnect a player with an optional reason.
Usage: `kick <name> [reason]`
```

> `ban`
```
Permission: cmd_ban
Description: Ban a player from using the proxy.
Usage: `ban <name>`
```

> `unban`
```
Permission: cmd_unban
Description: Remove a player from the ban list. Accepts addresses and names.
Usage: `unban <name | address>`
```

> `uptime`
```
Permission: cmd_uptime
Description: Show the uptime of the proxy.
Usage: `uptime`
```

> `help`
```
Permission: cmd_help
Description: Show help for a command (all commands if unspecified).
Usage: `help [command]`
```

> `usage`
```
Permission: cmd_usage
Description: Show the usage string of a command (all commands if unspecified).
Usage: `usage [command]`
```
