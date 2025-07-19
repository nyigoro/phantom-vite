#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status.

PROJECT_ROOT="$(dirname "$0")"/..
cd "$PROJECT_ROOT"

EXECUTABLE_NAME="phantom-vite"

# Build the phantom-vite executable
echo "Building phantom-vite executable..."
go build -o "$EXECUTABLE_NAME" ./cmd

# Run Agent Command Test 1
echo "Running Agent Command Test 1: summarize this document"
OUTPUT=$(./"$EXECUTABLE_NAME" agent "summarize this document")
EXPECTED_OUTPUT="This page is a sample domain used for illustrative purposes."
if [[ "$OUTPUT" == *"$EXPECTED_OUTPUT"* ]]; then
  echo "Test 1 Passed"
else
  echo "Test 1 Failed: Expected \"$EXPECTED_OUTPUT\" in output, got:\n$OUTPUT"
  exit 1
fi

# Run Agent Command Test 2
echo "Running Agent Command Test 2: hello world"
OUTPUT=$(./"$EXECUTABLE_NAME" agent "hello world")
EXPECTED_OUTPUT="[Agent Reply] I received: hello world"
if [[ "$OUTPUT" == *"$EXPECTED_OUTPUT"* ]]; then
  echo "Test 2 Passed"
else
  echo "Test 2 Failed: Expected \"$EXPECTED_OUTPUT\" in output, got:\n$OUTPUT"
  exit 1
fi

# Clean up the executable
echo "Cleaning up executable..."
rm "$EXECUTABLE_NAME"

echo "All integration tests passed!"

