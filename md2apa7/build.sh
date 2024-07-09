#!/bin/bash

# Download the APA CSL file
# ./download-csl.sh

# Build the Docker image
docker build -t md2apa7-converter .

cp md2apa7 ../bin/
