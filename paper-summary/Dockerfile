# Use the official Python image
FROM python:3.9

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy the current directory contents into the container
COPY . .

# Install the required libraries
RUN pip install openai==1.43.0 bibtexparser markdown

# Ensure the API key file is available
COPY openai_api_key.txt .

# Create input and output directories if they don't exist
RUN mkdir -p input output

# Set the entrypoint to run the script
ENTRYPOINT ["python", "./script.py"]
