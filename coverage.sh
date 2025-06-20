#!/bin/bash

# Code Coverage Generation Script
# This script generates comprehensive test coverage reports for the microservices application

set -e  # Exit on any error

echo "🚀 Starting comprehensive test coverage analysis..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go to run tests."
    exit 1
fi

print_status "Running unit tests with coverage..."

# Run tests with coverage, excluding main.go and test files
go test -coverprofile=coverage.out -covermode=atomic ./...

# Check if tests passed
if [ $? -ne 0 ]; then
    print_error "Tests failed. Please fix the failing tests before generating coverage report."
    exit 1
fi

print_success "All tests passed!"

# Filter out main.go and test files from coverage
print_status "Filtering coverage data..."
grep -v "main.go\|_test.go" coverage.out > coverage_filtered.out

# Generate HTML coverage report
print_status "Generating HTML coverage report..."
go tool cover -html=coverage_filtered.out -o coverage.html

# Generate function coverage report
print_status "Generating function coverage report..."
go tool cover -func=coverage_filtered.out > coverage_func.txt

# Calculate coverage percentage
COVERAGE=$(go tool cover -func=coverage_filtered.out | tail -1 | awk '{print $3}' | sed 's/%//')
COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')

print_status "Coverage analysis complete!"

# Display coverage summary
echo ""
echo "📊 Coverage Summary:"
echo "===================="
echo "Overall Coverage: ${COVERAGE}"

# Check coverage threshold (80%)
if (( $(echo "$COVERAGE_NUM >= 80" | bc -l) )); then
    print_success "✅ Coverage target (80%) achieved!"
else
    print_warning "⚠️  Coverage below target (80%). Current: ${COVERAGE}"
fi

# Display detailed coverage by package
echo ""
echo "📋 Package Coverage Details:"
echo "============================"
go tool cover -func=coverage_filtered.out | grep -E "total:|src/" | while read line; do
    if [[ $line == *"total:"* ]]; then
        echo "$line"
    else
        echo "  $line"
    fi
done

# Generate coverage badges (if available)
if command -v gocov-badge &> /dev/null; then
    print_status "Generating coverage badge..."
    gocov-badge coverage_filtered.out > coverage.svg
    print_success "Coverage badge generated: coverage.svg"
fi

# Clean up temporary files
print_status "Cleaning up temporary files..."
rm -f coverage.out

# Display file locations
echo ""
echo "📁 Generated Files:"
echo "==================="
echo "• coverage.html - HTML coverage report (open in browser)"
echo "• coverage_func.txt - Function-level coverage report"
echo "• coverage_filtered.out - Raw coverage data (filtered)"

# Open HTML report if possible
if command -v open &> /dev/null; then
    print_status "Opening HTML coverage report..."
    open coverage.html
elif command -v xdg-open &> /dev/null; then
    print_status "Opening HTML coverage report..."
    xdg-open coverage.html
else
    print_status "HTML report generated. Open coverage.html in your browser to view detailed coverage."
fi

print_success "Coverage analysis completed successfully!"
echo ""
echo "🎯 Next Steps:"
echo "• Review coverage.html for detailed coverage information"
echo "• Focus on uncovered areas to improve test coverage"
echo "• Aim for 80%+ coverage across all packages"
echo "• Run integration tests: ./Test/integration/run-integration-test.bash" 