# HBM (Harbormaster)

HBM is a Docker authorization plugin that controls what Docker commands users are allowed to run. It intercepts every Docker API call and checks it against a configurable whitelist of allowed actions, images, volumes, ports, capabilities, and more.

Users are identified by the Common Name (CN) in their TLS client certificate. Access is granted through a role-based model: users belong to groups, groups are assigned policies, policies link to collections of resources.

## How access control works

Authorization is **disabled by default**. After running `hbm init`, all Docker commands are permitted until you explicitly enable enforcement with `hbm config set authorization true`.

Once enabled, HBM operates as a whitelist — everything is denied unless explicitly allowed by policy. Without a matching policy, the following are among the things that will be blocked:

- `--privileged`
- `--ipc=host`, `--net=host`, `--pid=host`, `--userns=host`, `--uts=host`
- `--cap-add` (any capability)
- `--device`
- `--dns`
- Port bindings (`-p`)
- Volume mounts (`-v`)
- `--log-driver` and `--log-opt`
- `--sysctl`, `--security-opt`
- Pulling images not on the whitelist

## Documentation

See the [docs/](docs/) directory:

- [How it works](docs/how-it-works.md)
- [Installation](docs/installation.md)
- [Getting started](docs/getting-started.md)
- [Security](docs/security.md)
- [CLI reference](docs/reference/cli.md)

## Supported versions

| HBM version         | Go    | Docker API types               | Tested Docker Engine   |
|---------------------|-------|--------------------------------|------------------------|
| 0.23.x              | 1.26+ | moby/moby/api v1.56            | 28.1, 28.5, 29.7, 29.9 |
| 0.21.x – 0.22.x     | 1.26+ | moby/moby/api v1.54            | —                      |
| 0.20.x              | 1.25+ | docker/docker v28.5 (API 1.51) | —                      |
| 0.19.x              | 1.21+ | docker/docker v24.0 (API 1.43) | —                      |

The Go toolchain preferred by `go.mod` is 1.27.2. It is fetched automatically where `GOTOOLCHAIN=auto`. Go 1.26 also builds HBM.

## Docker API compatibility

HBM does not pin a Docker API version. It strips the optional `/vX.Y` prefix from each request URI and matches the remaining path against its endpoint table. Any client API version the daemon accepts is handled the same way.

Endpoints missing from that table are denied for every user, including administrators, and even with authorization disabled. A Docker release that adds a new endpoint therefore needs an HBM update before clients can use it.

## Dependencies

- [docker/go-plugins-helpers](https://github.com/docker/go-plugins-helpers) — Docker authorization plugin framework
- [jinzhu/gorm](https://github.com/jinzhu/gorm) + [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) — ORM and SQLite driver (pure Go, no CGO required)
- [moby/moby/api](https://github.com/moby/moby/tree/master/api) — Docker API types
- [moby/moby/client](https://github.com/moby/moby/tree/master/client) — Docker client
- [spf13/cobra](https://github.com/spf13/cobra) — CLI framework
