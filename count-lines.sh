#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <extension>"
    echo "Example: $0 js"
    echo "Example: $0 go"
    exit 1
fi

find . -name "*.$1" -type f | xargs wc -l 2>/dev/null | tail -1
