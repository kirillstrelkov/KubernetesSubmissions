#!/usr/bin/env bash
set -e

newurl=$(curl -s -D - -o /dev/null https://en.wikipedia.org/wiki/Special:Random | grep -i location | awk '{print $2}')

curl -X POST -d "body=Read $newurl" $URL

