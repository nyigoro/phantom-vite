#!/bin/bash
# scripts/dev-setup.sh - Local development setup

set -e

echo "🚀 Setting up Phantom Vite for local development..."

# Check prerequisites
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        exit 1
    fi
    echo "✅ $1 found"
}

echo "📋 Checking prerequisites..."
check_command "node"
check_command "npm"
check_command "go"

# Install Node.js dependencies
echo "📦 Installing Node.js dependencies..."
npm install

# Build TypeScript scripts
echo "🔨 Building TypeScript scripts..."
npx vite build

# Build Go CLI for current platform
echo "🏗️  Building Go CLI..."
go build -o phantom-vite ./cmd

# Make executable (Unix systems)
if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "win32" ]]; then
    chmod +x phantom-vite
fi

# Test the build
echo "🧪 Testing the build..."
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    ./phantom-vite.exe dist/example.js
else
    ./phantom-vite dist/example.js
fi

echo "✅ Development setup complete!"
echo ""
echo "🎯 Available commands:"
echo "  ./phantom-vite dist/example.js       - Run example script"
echo "  ./phantom-vite dist/puppet-test.js   - Run puppet test"
echo "  ./phantom-vite open https://google.com - Open URL"
echo "  ./phantom-vite build                 - Build project"
echo "  ./phantom-vite --help                - Show help"
echo ""
echo "🔧 To rebuild after changes:"
echo "  npm run build  # Rebuild TS scripts"
echo "  go build -o phantom-vite ./cmd  # Rebuild Go CLI"
