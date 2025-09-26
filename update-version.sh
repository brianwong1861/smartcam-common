#!/bin/bash

# SmartCam Common Package Version Update Script
# Updates logging time format and creates new version tag

set -e

echo "📦 Updating SmartCam Common Package..."

# Add changes to git
git add shared-logging/logger.go

# Commit changes
git commit -m "feat: update logging time format to YYYYMMDD:HHMMSS

- Add CustomTimeEncoder for better readability
- Change timestamp format from ISO8601 to YYYYMMDD:HHMMSS
- Update time key from 'timestamp' to 'time'

🤖 Generated with Claude Code"

# Create new version tag
NEW_VERSION="v1.3.0"
git tag -a $NEW_VERSION -m "Version $NEW_VERSION - Custom time format for logging

Features:
- Custom time encoder with YYYYMMDD:HHMMSS format
- Improved log readability for users
- Maintains all existing logging functionality

🤖 Generated with Claude Code"

# Push changes and tags
git push origin main
git push origin $NEW_VERSION

echo "✅ Successfully updated to version $NEW_VERSION"
echo "📝 Time format is now: YYYYMMDD:HHMMSS (e.g., 20250926:143022)"