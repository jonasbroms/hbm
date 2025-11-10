#!/bin/bash
#
# Migrate logrus to slog/logging package
# This script replaces logrus imports and usage with our slog wrapper
#

set -e

# Find all Go files (excluding vendor)
find . -name "*.go" -not -path "./vendor/*" -not -path "./pkg/logging/*" | while read -r file; do
    # Check if file uses logrus
    if grep -q "sirupsen/logrus" "$file"; then
        echo "Migrating: $file"

        # Replace logrus import with logging package
        sed -i 's|log "github.com/sirupsen/logrus"|"github.com/jonasbroms/hbm/pkg/logging"|g' "$file"
        sed -i 's|"github.com/sirupsen/logrus"|"github.com/jonasbroms/hbm/pkg/logging"|g' "$file"

        # Replace log.Fields with logging.Fields
        sed -i 's|log\.Fields{|logging.Fields{|g' "$file"

        # Replace log.WithFields with logging.WithFields
        sed -i 's|log\.WithFields(|logging.WithFields(|g' "$file"

        # Replace log.Fatal with logging.Fatal
        sed -i 's|log\.Fatal(|logging.Fatal(|g' "$file"

        # Replace log.Error with logging.Error
        sed -i 's|log\.Error(|logging.Error(|g' "$file"

        # Replace log.Warn with logging.Warn
        sed -i 's|log\.Warn(|logging.Warn(|g' "$file"

        # Replace log.Info with logging.Info
        sed -i 's|log\.Info(|logging.Info(|g' "$file"

        # Replace log.Debug with logging.Debug
        sed -i 's|log\.Debug(|logging.Debug(|g' "$file"
    fi
done

echo "Migration complete!"
echo "Please review the changes and run: go build ./..."
