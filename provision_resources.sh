#!/bin/bash

echo "Let's set up your resources"
sleep 1

# Read AWS credentials
read -p "Enter your AWS Access Key ID: " AWS_ACCESS_KEY_ID
read -sp "Enter your AWS Secret Access Key: " AWS_SECRET_ACCESS_KEY
echo  # For newline after entering secret key

echo "Warning: You need to give proper permission to this IAM user."

# Export AWS credentials as environment variables
export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" 
export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" 

# Ask for the AWS region
read -p "In which AWS region do you want to setup resources? (e.g., us-east-1): " AWS_REGION

# Change directories and update region in Pulumi configurations
cd pulumi

for dir in setup-s3-temp setup-s3-permanent setup-lambda setup-ecr setup-ecs-task-definition setup-ecs-cluster; do
  cd $dir
  pulumi config set aws:region "$AWS_REGION"
  cd ..
done

# Get temporary and permanent bucket names
read -p "Enter the name for the temporary S3 bucket: " TEMP_BUCKET_NAME
read -p "Enter the name for the permanent S3 bucket: " PERMANENT_BUCKET_NAME

# Export bucket names as environment variables
export TEMP_BUCKET_NAME="$TEMP_BUCKET_NAME"
export PERMANENT_BUCKET_NAME="$PERMANENT_BUCKET_NAME"

# Deploy Pulumi stacks for S3 and Lambda
cd setup-s3-temp
pulumi refresh --yes  # Sync the state before deploying
pulumi up --yes
cd ..

cd setup-s3-permanent
pulumi refresh --yes  # Sync the state before deploying
pulumi up --yes
cd ..

cd setup-lambda
pulumi refresh --yes  # Sync the state before deploying
pulumi up --yes
cd ..

# Ask for and validate ECR repository name
while true; do
  read -p "Enter a valid name for the ECR repository (lowercase, use - or _ if needed): " ECR_REPOSITORY_NAME

  if [[ $ECR_REPOSITORY_NAME =~ ^[a-z0-9]+([._-][a-z0-9]+)*$ ]]; then
    break
  else
    echo "Invalid repository name. Please ensure it is lowercase and follows AWS ECR naming rules."
  fi
done

# Export the repository name as an environment variable
export ECR_REPOSITORY_NAME="$ECR_REPOSITORY_NAME"

# Deploy Pulumi stack for ECR
cd setup-ecr
pulumi refresh --yes  # Sync the state before deploying
pulumi up --yes
cd ..

# Wait for 2 seconds
sleep 2
echo "You need to push a Docker image to this ECR repository in order to move further."

# Check if the user has pushed the Docker image
while true; do
  read -p "Have you pushed the Docker image to the ECR repository? (yes/no): " PUSHED_IMAGE

  if [[ $PUSHED_IMAGE == "yes" ]]; then
    break
  elif [[ $PUSHED_IMAGE == "no" ]]; then
    echo "Please push the Docker image and try again."
    exit 1
  else
    echo "Invalid response. Please enter 'yes' or 'no'."
  fi
done

# Get the image URI from the user
read -p "Enter the Docker image URI (you can find it in the ECR dashboard): " IMAGE_URI

# Export the image URI as an environment variable
export IMAGE_URI="$IMAGE_URI"

# Deploy Pulumi stacks for ECS task definition and cluster
cd setup-ecs-task-definition
pulumi refresh --yes  # Sync the state before deploying
pulumi up --yes
cd ..

cd setup-ecs-cluster
pulumi refresh --yes  # Sync the state before deploying
pulumi up --yes
cd ..

# Save the environment variables to a file
echo "Saving environment variables to env_vars.sh..."

{
  echo "export AWS_ACCESS_KEY_ID=\"$AWS_ACCESS_KEY_ID\""
  echo "export AWS_SECRET_ACCESS_KEY=\"$AWS_SECRET_ACCESS_KEY\""
  echo "export AWS_REGION=\"${AWS_REGION:-us-east-1}\""
  echo "export TEMP_BUCKET_NAME=\"${TEMP_BUCKET_NAME:-my-temp-bucket}\""
  echo "export PERMANENT_BUCKET_NAME=\"${PERMANENT_BUCKET_NAME:-my-permanent-bucket}\""
  echo "export ECR_REPOSITORY_NAME=\"${ECR_REPOSITORY_NAME:-my-ecr-repo}\""
  echo "export IMAGE_URI=\"${IMAGE_URI:-your-image-uri}\""
} > env_vars.sh || { echo "Failed to write to env_vars.sh"; }

# Source the env_vars.shs
source env_vars.sh

# Success message
echo "All your resources are setup successfully."
