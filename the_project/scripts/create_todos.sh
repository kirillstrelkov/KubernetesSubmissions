#!/bin/bash

URL="http://localhost:8081/posts"

echo "Starting to send 20 requests to $URL..."

for i in {1..20}
do
    CURRENT_TIME=$(date +"%H:%M:%S")
    MESSAGE="$CURRENT_TIME Here is #$i"
    
    echo "Sending request $i: body=$MESSAGE"
   
    curl -s -X POST "$URL" -d "body=$MESSAGE"
    
    echo ""
done

echo "---------------------------------"
echo "Finished sending 20 requests."