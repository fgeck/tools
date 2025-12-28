#!/bin/bash
# Pre-commit hook for Go projects
# This script runs go mod tidy, fmt, vet, and lint before each commit

set -e

echo "🔍 Running pre-commit checks..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    if [ "$status" = "success" ]; then
        echo -e "${GREEN}✅ ${message}${NC}"
    elif [ "$status" = "error" ]; then
        echo -e "${RED}❌ ${message}${NC}"
    elif [ "$status" = "info" ]; then
        echo -e "${YELLOW}ℹ️  ${message}${NC}"
    fi
}

# Check if running in git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_status "error" "Not a git repository"
    exit 1
fi

# Store the original stash state
STASH_NAME="pre-commit-$(date +%s)"
HAS_CHANGES=false

# Check if there are unstaged changes
if ! git diff-files --quiet; then
    HAS_CHANGES=true
    print_status "info" "Stashing unstaged changes..."
    git stash push -k -u -m "$STASH_NAME" > /dev/null 2>&1
fi

# Function to restore stashed changes on exit
cleanup() {
    if [ "$HAS_CHANGES" = true ]; then
        print_status "info" "Restoring unstaged changes..."
        git stash pop > /dev/null 2>&1 || true
    fi
}

# Set trap to always restore on exit
trap cleanup EXIT

# Track if any check fails
CHECKS_FAILED=false

echo "📦 Step 1/4: Running go mod tidy..."
if go mod tidy; then
    if ! git diff --exit-code go.mod go.sum > /dev/null 2>&1; then
        print_status "info" "go.mod or go.sum was modified by 'go mod tidy'"
        git add go.mod go.sum
        print_status "success" "Changes automatically staged"
    else
        print_status "success" "go mod tidy - no changes"
    fi
else
    print_status "error" "go mod tidy failed"
    CHECKS_FAILED=true
fi

echo ""
echo "🎨 Step 2/4: Running go fmt..."
# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -n "$STAGED_GO_FILES" ]; then
    # Format staged files
    echo "$STAGED_GO_FILES" | xargs -I {} go fmt {} > /dev/null

    # Check if formatting changed anything
    if ! git diff --exit-code $STAGED_GO_FILES > /dev/null 2>&1; then
        print_status "info" "Code was auto-formatted"
        echo ""
        echo "Formatted files:"
        git diff --name-only $STAGED_GO_FILES
        echo ""
        # Automatically add the formatted files
        git add $STAGED_GO_FILES
        print_status "success" "Changes automatically staged"
    else
        print_status "success" "go fmt - all files formatted correctly"
    fi
else
    print_status "success" "go fmt - no Go files to check"
fi

echo ""
echo "🔍 Step 3/4: Running go vet..."
if go vet ./...; then
    print_status "success" "go vet - no issues found"
else
    print_status "error" "go vet found issues"
    CHECKS_FAILED=true
fi

echo ""
echo "🧹 Step 4/4: Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run --timeout=5m ./...; then
        print_status "success" "golangci-lint - no issues found"
    else
        print_status "error" "golangci-lint found issues"
        CHECKS_FAILED=true
    fi
else
    print_status "info" "golangci-lint not installed - skipping"
    echo "   Install: brew install golangci-lint"
    echo "   Or: https://golangci-lint.run/usage/install/"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Final result
if [ "$CHECKS_FAILED" = true ]; then
    echo ""
    print_status "error" "Pre-commit checks FAILED"
    echo ""
    print_status "info" "Fix the issues and try again"
    print_status "info" "To skip this hook, use: git commit --no-verify"
    echo ""
    exit 1
else
    echo ""
    print_status "success" "All pre-commit checks PASSED"
    echo ""
    exit 0
fi
