#!/bin/bash

PROJECT_DIR=$(dirname "$0")
GENERATE_DIR="$PROJECT_DIR/gen"

cd "$GENERATE_DIR" || exit

echo "Start Generating"
go run .
