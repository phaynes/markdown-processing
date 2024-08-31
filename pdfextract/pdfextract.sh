#!/bin/bash

set -e

# Function to process a single PDF file
process_file() {
    local input_file="$1"
    local output_dir="$2"
    local processed_dir="$3"
    local filename=$(basename "$input_file" .pdf)

    echo "Processing: $input_file"

    # Use curl to send the PDF to Grobid's REST API
    if curl -s -S -X POST -F "input=@$input_file" \
         -F "consolidateHeader=1" \
         -F "consolidateCitations=1" \
         -F "teiCoordinates=persName,figure,ref" \
         -H "Accept: application/xml" \
         http://192.168.86.222:8070/api/processFulltextDocument > "$output_dir/${filename}.tei.xml"; then

        if [ -s "$output_dir/${filename}.tei.xml" ]; then
            echo "Successfully processed: $input_file"
            echo "XML output: $output_dir/${filename}.tei.xml"
            # Move the successfully processed PDF to the processed directory
            mv "$input_file" "$processed_dir/"
        else
            echo "Failed to process: $input_file (Empty output file)"
            rm -f "$output_dir/${filename}.tei.xml"  # Remove empty file
        fi
    else
        echo "Failed to process: $input_file (Curl command failed)"
    fi
}

# Check if Docker is installed and running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if input is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <input_pdf_or_directory> [output_directory]"
    exit 1
fi

input="$1"
output_dir="${2:-output}"
processed_dir="${output_dir}_processed"

# Create output and processed directories if they don't exist
mkdir -p "$output_dir"
mkdir -p "$processed_dir"

# Process input (file or directory)
if [ -f "$input" ]; then
    process_file "$input" "$output_dir" "$processed_dir"
elif [ -d "$input" ]; then
    for pdf in "$input"/*.pdf; do
        if [ -f "$pdf" ]; then
            process_file "$pdf" "$output_dir" "$processed_dir" || true
        fi
    done
else
    echo "Error: Input is neither a file nor a directory"
    exit 1
fi

echo "Processing complete. Output files are in $output_dir"
echo "Successfully processed PDFs have been moved to $processed_dir"


# Check if input is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <input_pdf_or_directory> [output_directory]"
    exit 1
fi

input="$1"
output_dir="${2:-output}"

# Create output directory if it doesn't exist
mkdir -p "$output_dir"

# Process input (file or directory)
if [ -f "$input" ]; then
    process_file "$input" "$output_dir"
elif [ -d "$input" ]; then
    for pdf in "$input"/*.pdf; do
        if [ -f "$pdf" ]; then
            process_file "$pdf" "$output_dir"
        fi
    done
else
    echo "Error: Input is neither a file nor a directory"
    exit 1
fi

echo "Processing complete. Output files are in $output_dir"
