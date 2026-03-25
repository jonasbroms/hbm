# HBM (Harbormaster)

HBM is a Docker authorization plugin that controls what Docker commands users are allowed to run. It intercepts every Docker API call and checks it against a configurable whitelist of allowed actions, images, volumes, ports, capabilities, and more.

Users are identified by the Common Name (CN) in their TLS client certificate. Access is granted through a role-based model: users belong to groups, groups are assigned policies, policies link to collections of resources.

## What it blocks by default

Without explicit policy, the following are always denied:

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

| HBM Version | Docker Version | Docker API |
|-------------|----------------|------------|
| 0.19.x      | 27.x           | 1.47       |

## Dependencies

- [docker/docker](https://github.com/docker/docker) — Docker client and API types
- [docker/go-connections](https://github.com/docker/go-connections) — Docker network helpers
- [docker/go-plugins-helpers](https://github.com/docker/go-plugins-helpers) — Docker authorization plugin framework
- [jinzhu/gorm](https://github.com/jinzhu/gorm) + [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) — ORM and SQLite driver
- [spf13/cobra](https://github.com/spf13/cobra) — CLI framework
- [go-utils](https://github.com/jonasbroms/go-utils) (forked from [juliengk/go-utils](https://github.com/juliengk/go-utils))
- [go-mount](https://github.com/kassisol/go-mount) (forked from [juliengk/go-mount](https://github.com/juliengk/go-mount))
