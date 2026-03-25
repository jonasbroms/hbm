# CLI Reference

## Global Commands

### `hbm init`

Initialize the HBM database and set default config values. Safe to run multiple times — skips keys that already exist.

```
hbm init [--action] [--config]
```

Creates `/var/lib/hbm/` and the SQLite database if they don't exist. Migrates legacy schemas automatically.

---

### `hbm server`

Start the HBM authorization plugin server. Writes `/etc/docker/plugins/hbm.spec` and listens on `/run/docker/plugins/hbm.sock`.

```
hbm server
```

Typically run via systemd. Handles graceful shutdown on SIGTERM/SIGINT.

---

### `hbm info`

Display HBM version, database path, and current config values.

```
hbm info
```

---

### `hbm version`

Show the HBM version string.

```
hbm version
```

---

## User Commands

### `hbm user add`

Add a user to the whitelist. The name must match the CN in the user's TLS client certificate.

```
hbm user add <name>
```

### `hbm user rm`

Remove a user from the whitelist.

```
hbm user rm <name>
hbm user remove <name>
```

### `hbm user ls`

List all whitelisted users.

```
hbm user ls [-f <filter>]
hbm user list [-f <filter>]
```

### `hbm user find`

Check if a user exists.

```
hbm user find <name>
```

### `hbm user member`

Add or remove a user from a group.

```
hbm user member <group> <user> --add
hbm user member <group> <user> --remove
```

---

## Group Commands

### `hbm group add`

Create a group.

```
hbm group add <name>
```

### `hbm group rm`

Remove a group. Fails if the group is referenced by a policy.

```
hbm group rm <name>
hbm group remove <name>
```

### `hbm group ls`

List all groups.

```
hbm group ls [-f <filter>]
hbm group list [-f <filter>]
```

### `hbm group find`

Check if a group exists.

```
hbm group find <name>
```

---

## Policy Commands

### `hbm policy add`

Create a policy linking a group to a collection.

```
hbm policy add <name> -g <group> -c <collection>
```

Flags:
- `-g` / `--group` — the group name
- `-c` / `--collection` — the collection name

### `hbm policy rm`

Remove a policy.

```
hbm policy rm <name>
hbm policy remove <name>
```

### `hbm policy ls`

List all policies.

```
hbm policy ls [-f <filter>]
hbm policy list [-f <filter>]
```

### `hbm policy find`

Check if a policy exists.

```
hbm policy find <name>
```

---

## Collection Commands

### `hbm collection add`

Create a collection.

```
hbm collection add <name>
```

### `hbm collection rm`

Remove a collection. Fails if referenced by a policy.

```
hbm collection rm <name>
hbm collection remove <name>
```

### `hbm collection ls`

List all collections.

```
hbm collection ls [-f <filter>]
hbm collection list [-f <filter>]
```

### `hbm collection find`

Check if a collection exists.

```
hbm collection find <name>
```

---

## Resource Commands

### `hbm resource add`

Add a resource to the whitelist.

```
hbm resource add <name> [-t <type>] [-v <value>] [-o <option>]
```

Flags:
- `-t` / `--type` — resource type (default: `action`). See [Resources](../configuration/resources.md) for all types.
- `-v` / `--value` — the resource value
- `-o` / `--option` — option modifier, can be repeated

### `hbm resource rm`

Remove a resource.

```
hbm resource rm <name>
hbm resource remove <name>
```

### `hbm resource ls`

List all resources.

```
hbm resource ls [-f <filter>]
hbm resource list [-f <filter>]
```

### `hbm resource find`

Check if a resource exists.

```
hbm resource find <name>
```

### `hbm resource member`

Add or remove a resource from a collection.

```
hbm resource member <collection> <resource> --add
hbm resource member <collection> <resource> --remove
```

---

## Config Commands

### `hbm config set`

Set a config value.

```
hbm config set <key> <value>
```

Available keys:

| Key | Values | Description |
|-----|--------|-------------|
| `authorization` | `true` / `false` | Enable or disable authorization enforcement |
| `default-allow-action-error` | `true` / `false` | Allow requests when an internal error occurs |

### `hbm config get`

Get the current value of a config key.

```
hbm config get <key>
hbm config find <key>
```

### `hbm config ls`

List all config keys and their values.

```
hbm config ls [-f <filter>]
hbm config list [-f <filter>]
```

---

## System Commands

`hbm system` currently has no subcommands beyond `hbm init`, `hbm info`, and `hbm version` (listed above).
