# Users and Groups

## Users

A user in HBM corresponds to a person who connects to Docker. The username must match the **Common Name (CN)** in their TLS client certificate.

### Managing Users

```bash
# Add a user
hbm user add <name>

# Remove a user
hbm user rm <name>

# List all users
hbm user ls

# Filter by name
hbm user ls -f name=alice

# Check if a user exists
hbm user find <name>
```

## Groups

Groups collect users together so you can assign permissions to the group rather than each user individually.

### The administrators Group

The `administrators` group is created automatically by `hbm init` and cannot be deleted. Members of this group bypass all authorization checks — they can run any Docker command.

```bash
# Grant a user full admin access
hbm user member administrators <username> --add
```

### Managing Groups

```bash
# Create a group
hbm group add <name>

# Remove a group (fails if the group is referenced by a policy)
hbm group rm <name>

# List all groups
hbm group ls

# Check if a group exists
hbm group find <name>
```

### Managing Group Membership

```bash
# Add a user to a group
hbm user member <group> <user> --add

# Remove a user from a group
hbm user member <group> <user> --remove
```

A user can belong to multiple groups. If any group grants access to a resource, the request is allowed.
