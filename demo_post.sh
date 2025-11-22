#!/bin/bash

echo "=========================================="
echo "HSP POST with JSON Body Demo"
echo "=========================================="
echo ""
echo "Creating a new post via JSONPlaceholder API"
echo ""

# Simulate user input for a POST request with JSON
(
  echo "https://jsonplaceholder.typicode.com/posts"  # URL
  echo "2"                                             # POST (choice 2)
  echo "n"                                             # No headers
  echo "n"                                             # No query params
  echo "y"                                             # Add body
  echo "1"                                             # JSON format
  echo '{'                                             # JSON body
  echo '  "title": "HSP Test Post",'
  echo '  "body": "This was created using HSP interactive builder",'
  echo '  "userId": 1'
  echo '}'
  echo ""                                              # Empty line
  echo ""                                              # Second empty line
  echo "y"                                             # Pretty print
  echo "y"                                             # Send
) | ./hsp request
