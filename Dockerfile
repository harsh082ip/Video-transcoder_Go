FROM ubuntu:latest

# Set environment variables
ENV AWS_ACCESS_KEY_ID=<your-access-key-id>
ENV AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
ENV AWS_DEFAULT_REGION=<your-region>
ENV SOURCE_IMAGE=default
ENV DESTINATION_1080=default
ENV DESTINATION_720=default
ENV DESTINATION_360=default

# Update and install necessary packages
RUN apt-get update && apt-get install -y \
    ffmpeg \
    unzip \
    curl

# Install AWS CLI
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install

# Copy the transcoding script into the container
COPY transcode.sh /usr/local/bin/transcode_and_upload.sh

# Make the script executable
RUN chmod +x /usr/local/bin/transcode_and_upload.sh

# Set the script as the entry point
ENTRYPOINT ["/usr/local/bin/transcode_and_upload.sh"]
