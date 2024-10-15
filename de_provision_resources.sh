#!/bin/bash

# Function to confirm if the user wants to change an environment variable
confirm_change() {
  read -p "Do you want to change $1? (yes/no): " CHANGE
  if [[ $CHANGE == "yes" ]]; then
    return 0
  else
    return 1
  fi
}

# Function to prompt for AWS credentials and display current environment variables
get_aws_credentials() {
  echo "AWS Credentials Setup"
  read -p "Enter your AWS Access Key ID: " AWS_ACCESS_KEY_ID
  read -sp "Enter your AWS Secret Access Key: " AWS_SECRET_ACCESS_KEY
  echo  # Newline after entering the secret key

  echo "Current environment variables:"
  echo "AWS_REGION=${AWS_REGION:-not set}"
  echo "TEMP_BUCKET_NAME=${TEMP_BUCKET_NAME:-not set}"
  echo "PERMANENT_BUCKET_NAME=${PERMANENT_BUCKET_NAME:-not set}"
  echo "ECR_REPOSITORY_NAME=${ECR_REPOSITORY_NAME:-not set}"
  echo "IMAGE_URI=${IMAGE_URI:-not set}"
  echo "KEY_PAIR_NAME=${KEY_PAIR_NAME:-not set}"

  # Check if the user wants to change any variables
  if confirm_change "AWS_REGION"; then
    read -p "Enter the AWS Region: " AWS_REGION
  fi

  if confirm_change "TEMP_BUCKET_NAME"; then
    read -p "Enter the temporary S3 bucket name: " TEMP_BUCKET_NAME
  fi

  if confirm_change "PERMANENT_BUCKET_NAME"; then
    read -p "Enter the permanent S3 bucket name: " PERMANENT_BUCKET_NAME
  fi

  if confirm_change "ECR_REPOSITORY_NAME"; then
    read -p "Enter the ECR repository name: " ECR_REPOSITORY_NAME
  fi

  if confirm_change "IMAGE_URI"; then
    read -p "Enter the Docker image URI: " IMAGE_URI
  fi

  if confirm_change "KEY_PAIR_NAME"; then
    read -p "Enter the key pair name: " KEY_PAIR_NAME
  fi

  # Save the environment variables to a file
  echo "Saving environment variables to env_vars.sh..."
  {
    echo "export AWS_REGION=${AWS_REGION:-us-east-1}"
    echo "export TEMP_BUCKET_NAME=${TEMP_BUCKET_NAME:-my-temp-bucket}"
    echo "export PERMANENT_BUCKET_NAME=${PERMANENT_BUCKET_NAME:-my-permanent-bucket}"
    echo "export ECR_REPOSITORY_NAME=${ECR_REPOSITORY_NAME:-my-ecr-repo}"
    echo "export IMAGE_URI=${IMAGE_URI:-your-image-uri}"
  } > env_vars.sh
}

# Call the function to get AWS credentials
get_aws_credentials

# Confirm with the user before proceeding to deprovision resources
read -p "Are you sure you want to deprovision the resources? (yes/no): " CONFIRM
if [[ $CONFIRM != "yes" ]]; then
  echo "Operation aborted."
  exit 1
fi

# Set the base directory for resource setup
BASE_DIR="$(pwd)"  # Get the current working directory

# Navigate to each directory and deprovision resources
for dir in setup-s3-temp setup-s3-permanent setup-lambda setup-ecr setup-ecs-task-definition setup-ecs-cluster setup-ec2; do
  cd "$BASE_DIR/pulumi/$dir"  # Use the path relative to the BASE_DIR
  pulumi down --yes
  cd "$BASE_DIR"  # Go back to the base directory
done

# Success message
echo "All resources have been successfully deprovisioned."
