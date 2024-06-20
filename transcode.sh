#!/bin/bash

# Function to convert HTTP URL to S3 URI
function http_to_s3_uri {
    local http_url=$1
    # Extract the bucket name and key from the HTTP URL
    local bucket_name=$(echo "$http_url" | sed -e 's#https://s3\.[^.]*\.amazonaws\.com/##' | cut -d'/' -f1)
    local key=$(echo "$http_url" | sed -e 's#https://s3\.[^.]*\.amazonaws\.com/##' | cut -d'/' -f2-)
    echo "s3://${bucket_name}/${key}"
}

# Check if required environment variables are set
if [ "$SOURCE_IMAGE" = "default" ] || [ "$DESTINATION_1080" = "default" ] || [ "$DESTINATION_720" = "default" ] || [ "$DESTINATION_360" = "default" ]; then
    echo "SOURCE_IMAGE, DESTINATION_1080, DESTINATION_720, and DESTINATION_360 environment variables must be set."
    exit 1
fi

# Transcode and upload functions
function transcode_and_upload {
    local resolution=$1
    local destination_var="DESTINATION_$resolution"
    local destination=${!destination_var}
    
    # Determine scale dimensions based on aspect ratio
    scale_filter="scale='if(gt(iw,ih),-2,$resolution)':'if(gt(iw,ih),$resolution,-2)'"
    
    ffmpeg -i "$SOURCE_IMAGE" -vf "$scale_filter" -c:v libx264 -c:a aac -f mp4 -movflags frag_keyframe+empty_moov - | aws s3 cp - "$destination"
}

# Transcode to different resolutions and upload
transcode_and_upload 1080
transcode_and_upload 720
transcode_and_upload 360

# Convert HTTP URL to S3 URI and delete the source file from S3
s3_uri=$(http_to_s3_uri "$SOURCE_IMAGE")
echo "Attempting to delete: $s3_uri" # Debug statement

# Manually decoding the key for verification
decoded_key=$(echo "$s3_uri" | sed -e 's/%40/@/g')
echo "Decoded Key: $decoded_key"

# Try deleting using the decoded key
delete_output=$(aws s3 rm "$decoded_key" 2>&1)
echo "Delete Output: $delete_output"

# Check if the deletion was successful
delete_status=$?
if [ $delete_status -eq 0 ]; then
    echo "Source file successfully deleted."
else
    echo "Failed to delete source file. Reason: $delete_output"
    exit 1
fi
