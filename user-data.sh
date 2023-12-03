#!/bin/bash

sudo apt update -y

# Install Docker
sudo apt install docker.io -y
sudo apt install awscli -y

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add the current user to the docker group
sudo usermod -aG docker ubuntu

# Login to AWS ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 072422391281.dkr.ecr.us-east-1.amazonaws.com

# Pull and run your Docker image
docker pull 072422391281.dkr.ecr.us-east-1.amazonaws.com/okta-event-hooks:latest
docker run -d -p 8080:8080 072422391281.dkr.ecr.us-east-1.amazonaws.com/okta-event-hooks:latest