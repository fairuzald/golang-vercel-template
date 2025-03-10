#!/bin/bash

# Set default output directory to "out" if not specified
OUTPUT_DIR="${1:-out}"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

echo "Will exclude: $OUTPUT_DIR/, node_modules/, .git/, vendor/, and some common binary/generated files"

# First, properly sanitize the OUTPUT_DIR to avoid issues with find
CLEAN_OUTPUT_DIR=$(realpath "$OUTPUT_DIR")
echo "Sanitized output directory path: $CLEAN_OUTPUT_DIR"

# Find all files, with proper exclusion of output directory
echo "Finding files to copy..."
FILES_TO_COPY=$(find . \
    -type f \
    -not -path "$CLEAN_OUTPUT_DIR/*" \
    -not -path "*/node_modules/*" \
    -not -path "*/.git/*" \
    -not -path "*/vendor/*" \
    -not -path "*/tmp/*" \
    -not -path "*/build/*" \
    -not -path "*/dist/*" \
    -not -path "*/.idea/*" \
    -not -path "*/.vscode/*" \
    -not -path "*/coverage/*" \
    -not -name "*.exe" \
    -not -name "*.dll" \
    -not -name "*.so" \
    -not -name "*.dylib" \
    -not -name "*.o" \
    -not -name "*.a" \
    -not -name "*.test" \
    -not -name "*.out" \
    -not -name ".DS_Store" \
    -not -name ".env*" \
    -not -name "*.log")

# Count total files to copy
TOTAL_FILES=$(echo "$FILES_TO_COPY" | wc -l)
echo "Found $TOTAL_FILES files to copy"

# Show sample of files that will be copied (up to 5)
echo "Sample files that will be copied:"
echo "$FILES_TO_COPY" | head -5

# Copy all files
echo "Copying files from current directory to '$OUTPUT_DIR'..."

# Track which files we've processed (for logging purposes only)
declare -A processed_files

echo "$FILES_TO_COPY" | while read -r file; do
    # Skip if the file is empty or in the output directory
    if [ -z "$file" ] || [[ "$(realpath "$file")" == "$CLEAN_OUTPUT_DIR"* ]]; then
        echo "Skipping: $file (in output directory)"
        continue
    fi

    # Get the filename without the path
    filename=$(basename "$file")

    # Check if this filename has been processed before
    if [[ -n "${processed_files[$filename]}" ]]; then
        echo "Overwriting: $OUTPUT_DIR/$filename (previously from ${processed_files[$filename]}, now from $file)"
    else
        echo "Copying: $file to $OUTPUT_DIR/$filename"
    fi

    # Always copy the file, overwriting if it already exists
    cp -f "$file" "$OUTPUT_DIR/$filename"
    processed_files[$filename]="$file"
done

# Count how many files were actually copied
FILE_COUNT=$(find "$OUTPUT_DIR" -type f | wc -l)
echo "Successfully copied $FILE_COUNT files to '$OUTPUT_DIR'"

echo "Operation complete. All files have been flattened and copied to '$OUTPUT_DIR'"

# Include GitHub workflow files if requested
if [[ "$2" == "include-workflows" ]]; then
    echo "Including GitHub workflow files..."
    mkdir -p "$OUTPUT_DIR/.github/workflows"
    cp -rv .github/workflows/* "$OUTPUT_DIR/.github/workflows/" 2>/dev/null || echo "No workflow files found"
fi
