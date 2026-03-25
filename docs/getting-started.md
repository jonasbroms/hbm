# Getting Started

This guide walks through setting up a basic policy that allows a group of users to run containers from approved images, with access to specific ports and volumes.

## Prerequisite

HBM is installed and running. Authorization is enabled:

```bash
hbm config set authorization true
```

## 1. Add Users

Add each user that will connect via Docker. The name must match the CN in their TLS client certificate:

```bash
hbm user add alice
hbm user add bob
```

## 2. Create a Group

```bash
hbm group add developers
hbm user member developers alice --add
hbm user member developers bob --add
```

## 3. Define Resources

Resources are the individual permissions. Create one resource per permission:

```bash
# Allow basic Docker actions
hbm resource add allow-info        -t action -v container_info
hbm resource add allow-ps          -t action -v container_list
hbm resource add allow-run         -t action -v container_create
hbm resource add allow-start       -t action -v container_start
hbm resource add allow-stop        -t action -v container_stop
hbm resource add allow-rm          -t action -v container_remove
hbm resource add allow-logs        -t action -v container_logs
hbm resource add allow-pull        -t action -v image_create
hbm resource add allow-images      -t action -v image_list

# Allow specific images
hbm resource add allow-nginx       -t image -v nginx
hbm resource add allow-alpine      -t image -v alpine

# Allow a volume path
hbm resource add allow-tmp         -t volume -v /tmp

# Allow a port
hbm resource add allow-port-8080   -t port -v 8080
```

## 4. Create a Collection

A collection groups resources together:

```bash
hbm collection add basic-docker

hbm resource member basic-docker allow-info    --add
hbm resource member basic-docker allow-ps      --add
hbm resource member basic-docker allow-run     --add
hbm resource member basic-docker allow-start   --add
hbm resource member basic-docker allow-stop    --add
hbm resource member basic-docker allow-rm      --add
hbm resource member basic-docker allow-logs    --add
hbm resource member basic-docker allow-pull    --add
hbm resource member basic-docker allow-images  --add
hbm resource member basic-docker allow-nginx   --add
hbm resource member basic-docker allow-alpine  --add
hbm resource member basic-docker allow-tmp     --add
hbm resource member basic-docker allow-port-8080 --add
```

## 5. Create a Policy

A policy links a group to a collection:

```bash
hbm policy add dev-policy -g developers -c basic-docker
```

## 6. Verify

Check what's configured:

```bash
hbm policy ls
hbm collection ls
hbm resource ls
```

Alice and Bob can now:
- Pull and run `nginx` and `alpine` images
- Mount `/tmp`
- Expose port `8080`
- But not: run `--privileged`, mount arbitrary paths, pull other images, or expose other ports

## What's Blocked by Default

Without explicit resources, the following are always denied:
- `--privileged`
- `--net=host`
- `--pid=host`
- Mounting paths not in the volume whitelist
- Pulling images not in the image whitelist
- Exposing ports not in the port whitelist
- Adding Linux capabilities not in the capability whitelist
