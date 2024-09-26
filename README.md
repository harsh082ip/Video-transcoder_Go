# Video Transcoding Service

A scalable video transcoding service built using Golang, Gin, AWS, Pulumi, MongoDB, and Redis. This service handles video uploads, transcodes them into multiple formats, and stores them securely on AWS S3. It showcases a robust architecture for high-demand video processing applications.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Setup](#setup)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [Scaling and Performance](#scaling-and-performance)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Authentication**: Signup and login functionalities with session-based authentication.
- **Secure Video Upload**: Uses AWS S3 pre-signed URLs for secure video uploads.
- **Event-Driven Processing**: AWS Lambda triggers upon video upload to process metadata.
- **Task Queue Management**: Redis is used to manage a queue for video processing tasks.
- **Concurrent Transcoding**: An EC2 worker retrieves tasks from Redis and spins up ECS containers to transcode videos into multiple formats (360p, 720p, 1080p).
- **Scalable Storage**: Transcoded videos are stored in a secure S3 bucket.
- **Highly Scalable**: The architecture supports dynamic scaling and modular management.

## Architecture

### Project Architecture Diagram

![image](https://github.com/user-attachments/assets/8be6429b-025b-42c6-ac37-5d7e202efc7b)



1. **User Uploads Video**: Users upload videos to an AWS S3 bucket (temporary storage) using a pre-signed URL.
2. **AWS Lambda Trigger**: A Lambda function triggers upon video upload and pushes metadata to a Redis-based queue.
3. **Queue Processing**: An EC2 instance running a worker retrieves data from the queue.
4. **Concurrent Transcoding**: The worker spins up ECS containers to transcode up to 5 videos concurrently into three formats (360p, 720p, 1080p).
5. **Permanent Storage**: The transcoded videos are stored in a permanent S3 bucket.
6. **Secure Access**: The service includes authentication and authorization mechanisms to ensure secure access.

## Tech Stack

- **Golang**: Backend development using the Gin framework.
- **AWS S3**: For temporary and permanent storage of video files.
- **AWS Lambda**: For event-driven processing.
- **AWS EC2 and ECS**: For worker instances and container orchestration.
- **Pulumi**: Infrastructure as Code (IaC) for managing cloud resources.
- **Redis**: For managing the task queue.
- **MongoDB**: For metadata storage and management.

## Setup

### Prerequisites

- [Golang](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Pulumi](https://www.pulumi.com/docs/get-started/install/) (for Infrastructure as Code)
- AWS account with access to S3, Lambda, EC2, and ECS
- [MongoDB](https://www.mongodb.com/try/download/community)

### Installation

1. **Clone the repository:**
    ```bash
    git clone https://github.com/harsh082ip/Video-transcoder_Go.git
    cd Video-transcoder_Go
    ```

2. **Set up environment variables:**
    Create a `.env` file in the project root directory and add the following:
    ```env
    AWS_ACCESS_KEY_ID=your_access_key
    AWS_SECRET_ACCESS_KEY=your_secret_key
    REDIS_URL=redis://localhost:6379
    MONGODB_URI=mongodb://localhost:27017
    AWS_DEFAULT_REGION=your_region
    SOURCE_IMAGE=your_source_image_path
    DESTINATION_1080=your_destination_path_for_1080p
    DESTINATION_720=your_destination_path_for_720p
    DESTINATION_360=your_destination_path_for_360p
    ```

3. **Deploy the infrastructure using Pulumi:**
    ```bash
    pulumi login
    pulumi stack init dev
    pulumi config set aws:region <your-region>
    pulumi up
    ```

## Usage

1. **Run the server:**
    ```bash
    go run cmd/getSignedUrl/main.go
    ```

2. **Authenticate Users:**
   - **Signup**: `/auth/signup`
   - **Login**: `/auth/login`

3. **Upload Videos:**
   - **Get Pre-Signed URL**: `/videos/getpresignedurl`  
   _(Pass session ID in the middleware for authorization.)_

## API Endpoints

- **POST** `/auth/signup`: Registers a new user.
- **POST** `/auth/login`: Logs in an existing user.
- **POST** `/videos/getpresignedurl`: Retrieves a pre-signed URL for uploading a video to S3.  
  _(Requires a valid session ID in the middleware.)_

## Environment Variables

The following environment variables must be set for ECS containers:

- `AWS_ACCESS_KEY_ID`: AWS access key ID.
- `AWS_SECRET_ACCESS_KEY`: AWS secret access key.
- `AWS_DEFAULT_REGION`: AWS region.
- `SOURCE_IMAGE`: Path of the source video image.
- `DESTINATION_1080`: Path to store the transcoded video in 1080p format.
- `DESTINATION_720`: Path to store the transcoded video in 720p format.
- `DESTINATION_360`: Path to store the transcoded video in 360p format.

## Scaling and Performance

- **Horizontal Scaling**: Utilizes ECS and Lambda for horizontal scaling to handle increasing video processing demands.
- **Cost-Efficiency**: Leverages AWS Lambda and on-demand ECS containers to optimize resource usage and reduce costs.
- **Fault Tolerance**: Designed to be fault-tolerant with retry mechanisms and proper error handling.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/YourFeature`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/YourFeature`
5. Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
