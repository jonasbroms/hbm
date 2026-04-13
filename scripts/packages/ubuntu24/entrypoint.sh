#!/usr/bin/env bash
set -e

VERSION=$1
RELEASE=$2

# Debian requires the version field to start with a digit.
# When there is no git tag the version script produces a commit hash (e.g.
# d57eb0d-dirty); use the 0~git. prefix which is the conventional snapshot
# format and sorts below any real release.
DEB_VERSION="${VERSION}"
if [[ "${DEB_VERSION}" != [0-9]* ]]; then
	DEB_VERSION="0~git.${DEB_VERSION}"
fi

PKGNAME="hbm_${DEB_VERSION}-${RELEASE}_amd64"
PKGDIR="/tmp/pkg/${PKGNAME}"

mkdir -p "${PKGDIR}/DEBIAN"
mkdir -p "${PKGDIR}/usr/sbin"
mkdir -p "${PKGDIR}/etc/systemd/system"
mkdir -p "${PKGDIR}/usr/share/bash-completion/completions"
mkdir -p "${PKGDIR}/usr/share/man/man8"

install -p -m 755 /usr/local/src/hbm/hbm "${PKGDIR}/usr/sbin/"
install -p -m 644 /usr/local/src/hbm/hbm.service "${PKGDIR}/etc/systemd/system/"
install -p -m 644 /usr/local/src/hbm/hbm.socket "${PKGDIR}/etc/systemd/system/"
install -p -m 644 /usr/local/src/hbm/shellcompletion/bash \
	"${PKGDIR}/usr/share/bash-completion/completions/hbm"

for f in /usr/local/src/hbm/man/man8/*.8; do
	gzip -9c "$f" > "${PKGDIR}/usr/share/man/man8/$(basename "$f").gz"
done

cat > "${PKGDIR}/DEBIAN/control" <<EOF
Package: hbm
Version: ${DEB_VERSION}-${RELEASE}
Architecture: amd64
Maintainer: Jonas Bröms <https://github.com/jonasbroms>
Description: Docker Engine Access Authorization Plugin
 HBM is an authorization plugin for Docker that intercepts every Docker API
 call and validates it against a whitelist of allowed actions.
Section: admin
Priority: optional
EOF

cd /tmp/pkg
fakeroot dpkg-deb --build "${PKGNAME}"

mkdir -p /tmp/dist
cp /tmp/pkg/*.deb /tmp/dist/
