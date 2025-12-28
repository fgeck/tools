#!/bin/bash
# Install Git hooks for the project

set -e

echo "🔧 Installing Git hooks..."
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GIT_HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

# Check if we're in a git repository
if [ ! -d "$PROJECT_ROOT/.git" ]; then
    echo -e "${YELLOW}⚠️  Not a git repository${NC}"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p "$GIT_HOOKS_DIR"

# Install pre-commit hook
echo "📝 Installing pre-commit hook..."
cp "$PROJECT_ROOT/scripts/pre-commit.sh" "$GIT_HOOKS_DIR/pre-commit"
chmod +x "$GIT_HOOKS_DIR/pre-commit"

echo ""
echo -e "${GREEN}✅ Git hooks installed successfully!${NC}"
echo ""
echo "The following hooks are now active:"
echo "  • pre-commit: Runs go mod tidy, fmt, vet, and lint"
echo ""
echo "To bypass the hook temporarily, use:"
echo "  git commit --no-verify"
echo ""
echo "To uninstall, run:"
echo "  rm .git/hooks/pre-commit"
echo ""
