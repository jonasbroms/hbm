# Collections

A collection is a named set of resources. It represents a permission profile — for example, "basic Docker user" or "CI/CD runner".

Collections are the middle layer in the permission model:

```
Group → Policy → Collection → Resources
```

## Managing Collections

```bash
# Create a collection
hbm collection add <name>

# Remove a collection (fails if referenced by a policy)
hbm collection rm <name>

# List all collections
hbm collection ls

# Check if a collection exists
hbm collection find <name>
```

## Adding and Removing Resources

```bash
# Add a resource to a collection
hbm resource member <collection> <resource> --add

# Remove a resource from a collection
hbm resource member <collection> <resource> --remove
```

## Example

```bash
hbm collection add ci-runner

hbm resource member ci-runner allow-pull    --add
hbm resource member ci-runner allow-run     --add
hbm resource member ci-runner allow-start   --add
hbm resource member ci-runner allow-stop    --add
hbm resource member ci-runner allow-rm      --add
hbm resource member ci-runner allow-logs    --add
hbm resource member ci-runner allow-nginx   --add
hbm resource member ci-runner allow-alpine  --add
```

This collection can then be assigned to a group via a [policy](policies.md).
