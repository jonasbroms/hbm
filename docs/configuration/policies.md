# Policies

A policy links a group to a collection. It is the final step that activates a set of permissions for a group of users.

```
Group → Policy → Collection → Resources
```

## Managing Policies

```bash
# Create a policy
hbm policy add <name> -g <group> -c <collection>

# Remove a policy
hbm policy rm <name>

# List all policies
hbm policy ls

# Check if a policy exists
hbm policy find <name>
```

## Example

```bash
hbm policy add dev-policy -g developers -c basic-docker
hbm policy add ci-policy  -g ci-runners  -c ci-runner
```

## How Policies Are Evaluated

When a user makes a Docker request:

1. HBM looks up which groups the user belongs to.
2. For each group, it finds all policies linked to that group.
3. For each policy, it checks if the requested resource is in the linked collection.
4. If any group/policy/collection chain grants access, the request is allowed.

A user with multiple group memberships gets the union of all their groups' permissions.

## Notes

- A group can have multiple policies (linked to different collections).
- A collection can be referenced by multiple policies.
- Removing a collection or group that is referenced by a policy will fail — remove the policy first.
