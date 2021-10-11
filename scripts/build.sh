#!/bin/sh

# Git project root
project_root="$(git rev-parse --show-toplevel)"
front_pub="$project_root/frontend/public"
back_pub="$project_root/backend/public"

# Build frontend
echo "Building frontend..."
(cd "$project_root/frontend" && npm run build)

# Build backend
echo "Building backend..."
[ -d "$back_pub" ] && rm -r "$back_pub"
cp -r "$front_pub" "$back_pub"
(cd "$project_root/backend" && go build -o ./build)