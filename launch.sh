#!/usr/bin/env sh

# Build container
docker build -t image-api .

# Run Docker Container
docker run -p 8080:8080 image-api --name image-api
