#!/bin/bash
#
# HBM TLS Certificate Setup
#
# Generates a CA, a Docker server certificate, and one client certificate per
# user. The client CN is the username that HBM reads for identity.
#
# Usage:
#   sudo ./setup-certs.sh [OPTIONS]
#
# Options:
#   --server HOST       Docker server hostname / IP (default: localhost)
#   --ca-dir DIR        Directory to store CA files (default: /etc/docker/certs)
#   --client-dir DIR    Directory for client certs (default: ~/.docker)
#   --user USERNAME     Generate a client cert for this user (repeatable)
#   --days N            Certificate validity in days (default: 365)
#   --force             Overwrite existing CA and server certs
#   --help              Show this message
#
# After running, add to /etc/docker/daemon.json:
#   {
#     "hosts": ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2376"],
#     "tlsverify": true,
#     "tlscacert": "/etc/docker/certs/ca.pem",
#     "tlscert":   "/etc/docker/certs/server-cert.pem",
#     "tlskey":    "/etc/docker/certs/server-key.pem"
#   }
#
# Each user should set in their shell:
#   export DOCKER_HOST=tcp://HOST:2376
#   export DOCKER_TLS_VERIFY=1
#   export DOCKER_CERT_PATH=~/.docker
#
set -e

log_info()    { echo -e "\033[0;34m[INFO]\033[0m $1"; }
log_success() { echo -e "\033[0;32m[SUCCESS]\033[0m $1"; }
log_warn()    { echo -e "\033[1;33m[WARN]\033[0m $1"; }
log_error()   { echo -e "\033[0;31m[ERROR]\033[0m $1" >&2; }

# Defaults
SERVER_HOST="localhost"
CA_DIR="/etc/docker/certs"
CLIENT_DIR="${HOME}/.docker"
DAYS=365
FORCE=false
USERS=()

parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --server)   SERVER_HOST="$2"; shift 2 ;;
            --ca-dir)   CA_DIR="$2";      shift 2 ;;
            --client-dir) CLIENT_DIR="$2"; shift 2 ;;
            --user)     USERS+=("$2");    shift 2 ;;
            --days)     DAYS="$2";        shift 2 ;;
            --force)    FORCE=true;       shift ;;
            --help)
                grep '^#' "$0" | grep -v '#!/bin/bash' | sed 's/^# \?//'
                exit 0
                ;;
            *) log_error "Unknown option: $1"; exit 1 ;;
        esac
    done
}

setup_ca() {
    mkdir -p "$CA_DIR"
    chmod 0700 "$CA_DIR"

    if [ -f "$CA_DIR/ca.pem" ] && [ "$FORCE" = false ]; then
        log_info "CA already exists at $CA_DIR/ca.pem (use --force to regenerate)"
        return
    fi

    log_info "Generating CA key and certificate..."
    openssl genrsa -out "$CA_DIR/ca-key.pem" 4096 2>/dev/null
    openssl req -new -x509 -days "$DAYS" \
        -key "$CA_DIR/ca-key.pem" \
        -sha256 \
        -subj "/CN=docker-ca" \
        -out "$CA_DIR/ca.pem"

    chmod 0400 "$CA_DIR/ca-key.pem"
    chmod 0444 "$CA_DIR/ca.pem"
    log_success "CA created: $CA_DIR/ca.pem"
}

setup_server() {
    if [ -f "$CA_DIR/server-cert.pem" ] && [ "$FORCE" = false ]; then
        log_info "Server cert already exists (use --force to regenerate)"
        return
    fi

    log_info "Generating server certificate for host: $SERVER_HOST"

    # Build SAN extension — include both IP and DNS forms.
    local ext_file
    ext_file=$(mktemp)
    {
        echo "subjectAltName = DNS:${SERVER_HOST},IP:127.0.0.1"
        echo "extendedKeyUsage = serverAuth"
    } > "$ext_file"

    openssl genrsa -out "$CA_DIR/server-key.pem" 4096 2>/dev/null
    openssl req -new -sha256 \
        -key "$CA_DIR/server-key.pem" \
        -subj "/CN=${SERVER_HOST}" \
        -out "$CA_DIR/server.csr"

    openssl x509 -req -days "$DAYS" -sha256 \
        -in  "$CA_DIR/server.csr" \
        -CA  "$CA_DIR/ca.pem" \
        -CAkey "$CA_DIR/ca-key.pem" \
        -CAcreateserial \
        -extfile "$ext_file" \
        -out "$CA_DIR/server-cert.pem"

    rm -f "$ext_file" "$CA_DIR/server.csr"
    chmod 0400 "$CA_DIR/server-key.pem"
    chmod 0444 "$CA_DIR/server-cert.pem"
    log_success "Server cert: $CA_DIR/server-cert.pem"
}

# Generate a client certificate for one user.
# The CN becomes the identity HBM reads from the TLS handshake.
setup_user_cert() {
    local username="$1"
    local dest

    # If the user exists on the system, put cert in their home .docker directory.
    if id "$username" &>/dev/null; then
        local home
        home=$(getent passwd "$username" | cut -d: -f6)
        dest="${home}/.docker"
    else
        dest="${CLIENT_DIR}"
        log_warn "OS user '$username' not found — placing cert in $dest"
    fi

    mkdir -p "$dest"

    log_info "Generating client certificate for user: $username (CN=$username)"

    local ext_file
    ext_file=$(mktemp)
    echo "extendedKeyUsage = clientAuth" > "$ext_file"

    openssl genrsa -out "$dest/key.pem" 4096 2>/dev/null
    openssl req -new -sha256 \
        -key "$dest/key.pem" \
        -subj "/CN=${username}" \
        -out "$dest/client.csr"

    openssl x509 -req -days "$DAYS" -sha256 \
        -in  "$dest/client.csr" \
        -CA  "$CA_DIR/ca.pem" \
        -CAkey "$CA_DIR/ca-key.pem" \
        -CAcreateserial \
        -extfile "$ext_file" \
        -out "$dest/cert.pem"

    # Copy CA cert for client-side verification.
    cp "$CA_DIR/ca.pem" "$dest/ca.pem"

    rm -f "$dest/client.csr" "$ext_file"
    chmod 0400 "$dest/key.pem" "$dest/ca.pem"
    chmod 0444 "$dest/cert.pem"

    if id "$username" &>/dev/null; then
        chown -R "${username}:${username}" "$dest"
    fi
    chmod 0700 "$dest"

    log_success "Client cert for '$username': $dest/"
    cat <<EOF
  → Add to ${username}'s shell profile:
      export DOCKER_HOST=tcp://${SERVER_HOST}:2376
      export DOCKER_TLS_VERIFY=1
      export DOCKER_CERT_PATH=${dest}
EOF
}

print_daemon_config() {
    cat <<EOF

Suggested /etc/docker/daemon.json:
{
  "hosts": ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2376"],
  "tlsverify": true,
  "tlscacert": "${CA_DIR}/ca.pem",
  "tlscert":   "${CA_DIR}/server-cert.pem",
  "tlskey":    "${CA_DIR}/server-key.pem",
  "authorization-plugins": ["hbm"]
}

Restart Docker after updating daemon.json:
  sudo systemctl restart docker
EOF
}

main() {
    parse_args "$@"

    if [[ $EUID -ne 0 ]]; then
        log_error "This script must run as root (CA and server keys go under $CA_DIR)"
        exit 1
    fi

    if ! command -v openssl &>/dev/null; then
        log_error "openssl is required but not found"
        exit 1
    fi

    setup_ca
    setup_server

    if [ ${#USERS[@]} -eq 0 ]; then
        log_warn "No --user specified; skipping client certificate generation"
        log_info "Run again with --user USERNAME to generate a client cert"
    else
        for u in "${USERS[@]}"; do
            setup_user_cert "$u"
        done
    fi

    print_daemon_config
    log_success "Certificate setup complete"
}

main "$@"
