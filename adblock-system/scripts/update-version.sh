#!/bin/bash
if [ -z "$1" ]; then
  echo "Usage: $0 <new-version>"
  exit 1
fi
NEW_VERSION="$1"
sed -i "s/Version *= *\".*\"/Version = \"$NEW_VERSION\"/" proxy/internal/config/config.go
git add proxy/internal/config/config.go
git commit -m "Bump version to $NEW_VERSION"
git tag "v$NEW_VERSION"
echo "Tagged v$NEW_VERSION. Push with: git push --tags"
