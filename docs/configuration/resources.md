# Resources

A resource is a single permission entry — one allowed action, image, volume path, port, etc. Resources are grouped into [collections](collections.md), which are then assigned to groups via [policies](policies.md).

## Adding a Resource

```bash
hbm resource add <name> -t <type> -v <value> [-o <option>]
```

- `<name>` — a unique identifier for this resource (your choice)
- `-t` / `--type` — the resource type (default: `action`)
- `-v` / `--value` — the value for this resource
- `-o` / `--option` — optional modifier (can be specified multiple times)

## Resource Types

### action

Allows a specific Docker command. Action names correspond to Docker API operations.

```bash
hbm resource add allow-ps      -t action -v container_list
hbm resource add allow-run     -t action -v container_create
hbm resource add allow-start   -t action -v container_start
hbm resource add allow-stop    -t action -v container_stop
hbm resource add allow-rm      -t action -v container_remove
hbm resource add allow-logs    -t action -v container_logs
hbm resource add allow-pull    -t action -v image_create
hbm resource add allow-images  -t action -v image_list
hbm resource add allow-exec    -t action -v container_exec_create
```

Use `hbm info` to see the full list of available action names.

### image

Allows pulling and running a specific image. The value can be a full image reference or a short name.

```bash
# Exact image name (any tag)
hbm resource add allow-nginx   -t image -v nginx

# With specific tag
hbm resource add allow-nginx   -t image -v nginx:1.25

# Allow all images under a registry path
hbm resource add allow-internal -t image -v registry.internal/myteam/ -o subimages
```

The `subimages` option allows any image whose name starts with the given prefix (prefix matching with trailing `/`).

### volume

Allows mounting a specific path as a bind mount.

```bash
# Exact path
hbm resource add allow-tmp    -t volume -v /tmp

# Allow mounting /data and any subdirectory
hbm resource add allow-data   -t volume -v /data -o recursive

# Allow mounting /data with nosuid enforcement
hbm resource add allow-data   -t volume -v /data -o nosuid
```

Options:
- `recursive` — also allows mounting any subdirectory under the given path
- `nosuid` — only allows the mount if nosuid is set

### port

Allows exposing a specific port.

```bash
hbm resource add allow-http   -t port -v 80
hbm resource add allow-https  -t port -v 443
hbm resource add allow-dev    -t port -v 8080
```

### capability

Allows adding a specific Linux capability with `--cap-add`.

```bash
hbm resource add allow-net-admin -t capability -v NET_ADMIN
hbm resource add allow-sys-time  -t capability -v SYS_TIME
```

Valid values include: `AUDIT_CONTROL`, `AUDIT_WRITE`, `CHOWN`, `DAC_OVERRIDE`, `DAC_READ_SEARCH`, `FOWNER`, `FSETID`, `IPC_LOCK`, `KILL`, `MKNOD`, `NET_ADMIN`, `NET_BIND_SERVICE`, `NET_RAW`, `SETGID`, `SETUID`, `SYS_ADMIN`, `SYS_CHROOT`, `SYS_MODULE`, `SYS_PTRACE`, `SYS_TIME`, and others.

### dns

Allows specifying a custom DNS server with `--dns`.

```bash
hbm resource add allow-dns-google -t dns -v 8.8.8.8
hbm resource add allow-dns-cf     -t dns -v 1.1.1.1
```

### device

Allows passing a host device with `--device`.

```bash
hbm resource add allow-gpu -t device -v /dev/nvidia0
```

### logdriver

Allows using a specific log driver with `--log-driver`.

```bash
hbm resource add allow-json-log -t logdriver -v json-file
hbm resource add allow-journald -t logdriver -v journald
```

### logopt

Allows setting a specific log driver option with `--log-opt`.

```bash
hbm resource add allow-log-maxsize -t logopt -v max-size
```

### volumedriver

Allows using a specific volume driver with `--volume-driver`.

```bash
hbm resource add allow-local-driver -t volumedriver -v local
```

### registry

Allows pushing or pulling from a specific registry.

```bash
hbm resource add allow-docker-hub  -t registry -v registry-1.docker.io
hbm resource add allow-internal    -t registry -v registry.internal
```

### runtime

Allows specifying a container runtime with `--runtime`.

```bash
hbm resource add allow-runc  -t runtime -v runc
```

## Managing Resources

```bash
# Remove a resource
hbm resource rm <name>

# List all resources
hbm resource ls

# Filter by name or type
hbm resource ls -f name=allow-nginx
hbm resource ls -f type=volume

# Check if a resource exists
hbm resource find <name>
```
