#!/bin/bash

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
    
    ffmpeg -i "$SOURCE_IMAGE" -vf scale=-1:$resolution -c:v libx264 -c:a aac -f mp4 -movflags frag_keyframe+empty_moov - | aws s3 cp - "$destination"
}

# Transcode to different resolutions and upload
transcode_and_upload 1080
transcode_and_upload 720
transcode_and_upload 360
