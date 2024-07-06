#!/bin/bash

# Download the APA CSL file
./download-csl.sh

# Build the Docker image
docker build -t apa7-converter .
