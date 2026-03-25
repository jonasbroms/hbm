# HBM Documentation

HBM (Harbormaster) is a Docker authorization plugin that controls what Docker commands users are allowed to run. It intercepts every Docker API call and checks it against a whitelist of allowed actions, images, volumes, ports, and more.

## How it works

HBM registers itself as a Docker authorization plugin via a Unix socket. When a user runs a Docker command, Docker asks HBM whether to allow or deny it before executing. HBM identifies users from the **Common Name (CN)** field of their TLS client certificate.

Users in the built-in `administrators` group bypass all checks. Everyone else is subject to the policies you configure.

## Documentation

- [Installation](installation.md) — Build, install, and wire up HBM with Docker
- [Getting Started](getting-started.md) — Set up your first users, groups, and policies
- [How It Works](how-it-works.md) — Architecture, identity model, authorization flow
- [Configuration Guide](configuration/)
  - [Users & Groups](configuration/users-and-groups.md)
  - [Resources](configuration/resources.md)
  - [Collections](configuration/collections.md)
  - [Policies](configuration/policies.md)
- [Security](security.md) — TLS setup, identity model, container ownership
- [Operations](operations.md) — Maintenance, migration, troubleshooting
- [CLI Reference](reference/cli.md) — All commands and flags
