#!/bin/bash
# scripts/build-all.sh - Build for all platforms

set -e

echo "🌍 Building Phantom Vite for all platforms..."

# Build TypeScript scripts first
echo "🔨 Building TypeScript scripts..."
npx vite build

# Create build directory
mkdir -p builds

# Build configurations
declare -A builds=(
    ["linux-amd64"]="linux amd64"
    ["linux-arm64"]="linux arm64"
    ["windows-amd64"]="windows amd64"
    ["windows-arm64"]="windows arm64"
    ["darwin-amd64"]="darwin amd64"
    ["darwin-arm64"]="darwin arm64"
)

# Build for each platform
for platform in "${!builds[@]}"; do
    read -r GOOS GOARCH <<< "${builds[$platform]}"
    
    echo "🏗️  Building for $platform ($GOOS/$GOARCH)..."
    
    # Set output filename
    output="builds/phantom-vite-$platform"
    if [ "$GOOS" = "windows" ]; then
        output="$output.exe"
    fi
    
    # Build
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
        -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty)" \
        -o "$output" \
        ./cmd
    
    echo "✅ Built: $output"
done

# Create release packages
echo "📦 Creating release packages..."
for platform in "${!builds[@]}"; do
    read -r GOOS GOARCH <<< "${builds[$platform]}"
    
    # Create platform directory
    platform_dir="builds/phantom-vite-$platform-release"
    mkdir -p "$platform_dir"
    
    # Copy executable
    if [ "$GOOS" = "windows" ]; then
        cp "builds/phantom-vite-$platform.exe" "$platform_dir/"
    else
        cp "builds/phantom-vite-$platform" "$platform_dir/"
    fi
    
    # Copy required files
    cp -r dist "$platform_dir/"
    cp -r plugins "$platform_dir/" 2>/dev/null || echo "No plugins directory"
    cp phantomvite.config.json "$platform_dir/" 2>/dev/null || echo "No config file"
    cp README.md "$platform_dir/"
    
    # Create archive
    cd builds
    if [ "$GOOS" = "windows" ]; then
        zip -r "phantom-vite-$platform.zip" "phantom-vite-$platform-release/"
        echo "📦 Created: builds/phantom-vite-$platform.zip"
    else
        tar -czf "phantom-vite-$platform.tar.gz" "phantom-vite-$platform-release/"
        echo "📦 Created: builds/phantom-vite-$platform.tar.gz"
    fi
    cd ..
done

echo "🎉 All builds complete!"
echo "📁 Check the 'builds/' directory for executables and release packages"
echo ""
echo "🚀 To create a GitHub release:"
echo "  git tag v1.0.0"
echo "  git push origin v1.0.0"
