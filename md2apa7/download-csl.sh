#!/bin/bash

# Download the APA CSL file
wget https://raw.githubusercontent.com/citation-style-language/styles/master/apa.csl -O apa.csl

# Check if the download was successful
if [ $? -eq 0 ]; then
    echo "APA CSL file downloaded successfully."
else
    echo "Failed to download APA CSL file. Please check your internet connection and try again."
    exit 1
fi
