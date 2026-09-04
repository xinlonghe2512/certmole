#!/usr/bin/env bash

set -e

# Get the latest git tag description
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

echo "📝 Generating release notes since version: ${LAST_TAG:-'Initial Commit'}"
echo "--------------------------------------------------"

if [ -z "$LAST_TAG" ]; then
    # If no tag exists, show all commits
    git log --oneline --pretty=format:"* %s (%an)"
else
    # Show commits since last release tag
    git log "${LAST_TAG}..HEAD" --oneline --pretty=format:"* %s (%an)"
fi
echo -e "\n--------------------------------------------------"