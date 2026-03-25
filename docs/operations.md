# Operations

## Checking Status

```bash
hbm info
```

Shows the current HBM version, database path, and config values.

## Config Management

```bash
# View all config keys and their current values
hbm config ls

# Get a specific value
hbm config get authorization

# Set a value
hbm config set authorization true
hbm config set default-allow-action-error false
```

## Cleaning Up Orphaned Records

When Docker containers are removed outside of normal operation (e.g., Docker daemon crash, manual database edits), the container ownership records in HBM's database may become stale. Currently this requires manual cleanup via SQLite:

```bash
systemctl stop hbm
sqlite3 /var/lib/hbm/data.db "DELETE FROM container_owner WHERE container_id NOT IN (SELECT id FROM containers);"
systemctl start hbm
```

## Database Backup

The entire HBM state is in a single SQLite file:

```bash
cp /var/lib/hbm/data.db /var/lib/hbm/data.db.bak
```

Restore by stopping HBM, replacing the file, and restarting:

```bash
systemctl stop hbm
cp /var/lib/hbm/data.db.bak /var/lib/hbm/data.db
systemctl start hbm
```

## Database Migration

If you are upgrading from an older version of HBM with a legacy database schema, `hbm init` will detect and migrate the schema automatically. It preserves all existing users, groups, policies, collections, and resources.

To test migration against a specific database file before upgrading:

```bash
cp /var/lib/hbm/data.db /tmp/test-migrate.db
HBM_APP_PATH=/tmp/test-migrate hbm init
```

## Troubleshooting

### All Docker commands are denied

Check that authorization is enabled and that the user is in the database:

```bash
hbm config get authorization
hbm user find <username>
hbm user ls
```

Also verify the user's TLS certificate CN matches exactly what HBM has on record.

### User is not being recognized

HBM reads the username from Docker's TLS client certificate CN. Check:

```bash
# Inspect the CN field of a certificate
openssl x509 -in ~/.docker/cert.pem -noout -subject
```

The CN must match the username registered with `hbm user add`.

### HBM is not intercepting requests

Check that the plugin socket exists and HBM is running:

```bash
systemctl status hbm
ls -la /run/docker/plugins/hbm.sock
cat /etc/docker/plugins/hbm.spec
```

Verify Docker's daemon config includes the authorization plugin:

```bash
cat /etc/docker/daemon.json | grep authorization
```

### Stale container ownership errors

A user gets denied access to their own container after a Docker restart or migration. See [Cleaning Up Orphaned Records](#cleaning-up-orphaned-records) above.

## Log Inspection

HBM logs all authorization decisions to the system journal:

```bash
journalctl -u hbm -f
```

Denied requests are logged at `WARN` level and include the username, action, and the specific resource that was denied (resource type, value, and denial reason).
