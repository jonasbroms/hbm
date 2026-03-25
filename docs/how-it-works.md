# How HBM Works

## Docker Authorization Plugin Model

Docker has a built-in authorization plugin API. When a plugin is registered, Docker forwards every API request to the plugin before executing it. The plugin responds with allow or deny.

HBM registers itself by writing a spec file at `/etc/docker/plugins/hbm.spec` and listening on a Unix socket at `/run/docker/plugins/hbm.sock`. Docker reads the spec file on startup and routes authorization requests to HBM.

```
User runs: docker run ...
     │
     ▼
Docker daemon receives API request
     │
     ▼
Docker asks HBM: allow this?  ──► HBM checks policies ──► allow / deny
     │
     ▼ (if allowed)
Docker executes the request
```

## User Identity

HBM identifies users from the **Common Name (CN) field** of the TLS client certificate used to connect to Docker. This requires Docker to listen on a TLS-enabled TCP port (typically `2376`) rather than just the Unix socket.

When Docker is configured for TLS, each user has their own client certificate with their username as the CN. HBM reads `req.User` from the authorization request, which Docker populates from the certificate CN.

If no TLS certificate is present (e.g., connecting via the Unix socket), the username defaults to `root`.

## Authorization Model

HBM uses a layered whitelist model:

```
User → Group → Policy → Collection → Resources
```

- **Users** are individuals identified by their TLS cert CN.
- **Groups** are sets of users. A user can belong to multiple groups.
- **Collections** are named sets of resources (what's allowed).
- **Policies** link a group to a collection.
- **Resources** are the individual permissions: allowed actions, images, volumes, ports, etc.

A user is allowed to perform an action if any of their groups has a policy granting access to a collection that contains the required resource.

### The administrators group

Users in the `administrators` group bypass all authorization checks. The group is created automatically by `hbm init` and cannot be deleted.

### When authorization is disabled

The `authorization` config key controls whether checks are enforced. When set to `false`, all requests are allowed regardless of policy. This is the default after `hbm init` — you must explicitly enable it with:

```
hbm config set authorization true
```

## Data Storage

All configuration (users, groups, policies, collections, resources) is stored in a SQLite database at `/var/lib/hbm/`. No static config files — everything is managed via the `hbm` CLI or updated live while the server is running.
