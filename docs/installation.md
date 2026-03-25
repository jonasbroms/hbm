# Installation

## Prerequisites

- Linux (RHEL/Fedora or compatible)
- Docker CE installed and running
- Go 1.25+ (for building from source)
- Docker configured for TLS (required for multi-user setups — see [Security](security.md))

## Build from Source

```bash
git clone https://github.com/jonasbroms/hbm.git
cd hbm
go build -o hbm .
install -m 755 hbm /usr/local/sbin/hbm
```

## Initialize the Database

```bash
hbm init
```

This creates `/var/lib/hbm/` (mode `0700`) and initializes the SQLite database with default config values and the `administrators` group.

## Systemd Setup

Create the socket unit at `/etc/systemd/system/hbm.socket`:

```ini
[Unit]
Description=HBM Docker Authorization Plugin Socket

[Socket]
ListenStream=/run/docker/plugins/hbm.sock
SocketMode=0660


[Install]
WantedBy=sockets.target
```

Create the service unit at `/etc/systemd/system/hbm.service`:

```ini
[Unit]
Description=HBM Docker Authorization Plugin
Documentation=https://github.com/jonasbroms/hbm
Before=docker.service
After=network.target hbm.socket
Requires=hbm.socket
Wants=docker.service

[Service]
Type=simple
ExecStartPre=-/usr/local/sbin/hbm init
ExecStart=/usr/local/sbin/hbm server
Restart=on-failure
RestartSec=10s

# Security
NoNewPrivileges=yes
PrivateTmp=yes
ProtectHome=yes
ReadWritePaths=/var/lib/hbm /run/docker/plugins /etc/docker/plugins

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
systemctl daemon-reload
systemctl enable --now hbm.socket hbm.service
```

## Register HBM with Docker

Add HBM to Docker's authorization plugins in `/etc/docker/daemon.json`:

```json
{
  "authorization-plugins": ["hbm"],
  "hosts": ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2376"],
  "tls": true,
  "tlsverify": true,
  "tlscacert": "/etc/docker/certs/ca.pem",
  "tlscert": "/etc/docker/certs/server-cert.pem",
  "tlskey": "/etc/docker/certs/server-key.pem"
}
```

The cert paths above (`/etc/docker/certs/`) refer to the CA and server certificates generated in the [TLS setup in the Security guide](security.md#tls-and-user-identity). Complete that step first, then come back here.

Then restart Docker:

```bash
systemctl restart docker
```

> **Note:** The TLS configuration is required for multi-user setups. Without TLS, all Docker connections are treated as `root`.

## Enable Authorization

After installation, authorization checking is disabled by default. Add the root user to administrators and enable it:

```bash
hbm user add root
hbm user member administrators root --add
hbm config set authorization true
```

At this point, `root` has full access and all other users are denied everything until you configure policies. See [Getting Started](getting-started.md).
