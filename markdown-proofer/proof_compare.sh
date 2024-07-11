# #!/bin/bash

# Usage: ./proof_compare.sh <file> <start_line> <end_line>

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <file> <start_line> <end_line>"
    exit 1
fi

FILE=$1
START_LINE=$2
END_LINE=$3
TEMP_FILE="temp_${FILE}"
PROOFED_FILE="proofed_${FILE}"

# Check if the file exists
if [ ! -f "$FILE" ]; then
    echo "File $FILE does not exist."
    exit 1
fi

# Extract the specified lines from the original file to create a temporary file
sed -n "${START_LINE},${END_LINE}p" "$FILE" > "$TEMP_FILE"

# Run the proofing command on the temporary file
./mdp -input "$TEMP_FILE" -output "$PROOFED_FILE"

# Compare the original selected lines to the proofed version
diff "$TEMP_FILE" "$PROOFED_FILE"

# Clean up temporary files
rm "$TEMP_FILE" "$PROOFED_FILE"
