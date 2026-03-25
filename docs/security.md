# Security

## TLS and User Identity

HBM identifies users from the **Common Name (CN)** field of the TLS client certificate used to connect to Docker. This means Docker must be configured to listen on a TLS-enabled TCP port.

### Generating Certificates

#### CA and server certificates

Run the following **on the Docker host, as root**. This is a one-time setup.

```bash
# Generate CA key and certificate
openssl genrsa -out ca-key.pem 4096
openssl req -new -x509 -days 3650 -key ca-key.pem -out ca.pem \
  -subj "/CN=docker-ca"

# Generate server key and certificate
openssl genrsa -out server-key.pem 4096
openssl req -new -key server-key.pem -out server.csr \
  -subj "/CN=$(hostname)"
openssl x509 -req -days 3650 -in server.csr \
  -CA ca.pem -CAkey ca-key.pem -CAcreateserial \
  -out server-cert.pem

# Install
install -d /etc/docker/certs
install -m 600 ca-key.pem server-key.pem /etc/docker/certs/
install -m 644 ca.pem server-cert.pem /etc/docker/certs/
```

Keep `ca-key.pem` secure — it is used to sign all user certificates.

#### Client certificates

Run the following **on the Docker host, as root**, once per user. The CN must exactly match the username registered with `hbm user add`.

```bash
USERNAME=alice

openssl genrsa -out ${USERNAME}-key.pem 4096
openssl req -new -key ${USERNAME}-key.pem -out ${USERNAME}.csr \
  -subj "/CN=${USERNAME}"
openssl x509 -req -days 365 -in ${USERNAME}.csr \
  -CA ca.pem -CAkey ca-key.pem -CAcreateserial \
  -out ${USERNAME}-cert.pem

# Place the certificate files where the user can access them
mkdir -p /home/${USERNAME}/.docker
cp ca.pem                   /home/${USERNAME}/.docker/ca.pem
cp ${USERNAME}-key.pem      /home/${USERNAME}/.docker/key.pem
cp ${USERNAME}-cert.pem     /home/${USERNAME}/.docker/cert.pem
chown -R ${USERNAME}: /home/${USERNAME}/.docker
```

Add to the user's shell profile (on whichever machine they run Docker commands from):

```bash
export DOCKER_HOST=tcp://<docker-host>:2376
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=~/.docker
```

#### Automating certificate management

For larger deployments, managing certificates manually per user does not scale well. [TSA](https://github.com/kassisol/tsa) is a self-hosted Certificate Authority server designed specifically for Docker TLS setups, and [TWIC](https://github.com/kassisol/twic) is its companion client tool that users run to request and configure their own certificates automatically. Together they replace the manual openssl steps above.

### Fallback to root

If Docker is accessed via the Unix socket (e.g., directly on the host without TLS), or if no TLS certificate CN is present, HBM treats the request as coming from `root`. Root is always in the `administrators` group and bypasses all checks.

## Container Ownership

HBM tracks which user created each container. Ownership is stored in the database at container creation time and updated on rename and removal.

A user can only interact with containers they own. Attempting to stop, remove, exec into, or otherwise interact with another user's container will be denied, regardless of the action whitelist.

**Example:** even if Alice has `container_stop` in her collection, she cannot stop a container created by Bob.

This protection is always active for non-admin users and cannot be bypassed through policy configuration.

## What Administrators Bypass

Members of the `administrators` group skip all authorization checks:

- No action whitelist
- No image whitelist
- No volume, port, or capability restrictions
- No container ownership enforcement

Be careful about who you add to `administrators`.

## Volume Path Security

HBM validates volume paths before checking them against policy:

- Symlinks are resolved before policy lookup, preventing symlink-based traversal attacks
- Paths containing `..` traversal components are rejected
- Trailing slashes are normalized before matching

## Config Keys Affecting Security

| Key | Default | Description |
|-----|---------|-------------|
| `authorization` | `false` | Master switch for authorization enforcement |
| `default-allow-action-error` | `false` | Whether to allow requests when an internal error occurs |

Both default to `false` (secure). After enabling `authorization`, set `default-allow-action-error` based on your preference for fail-open vs fail-closed behavior on errors.
