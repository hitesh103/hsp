#!/bin/bash

echo "=========================================="
echo "HSP Advanced Demo - GET with Query Params"
echo "=========================================="
echo ""
echo "Testing: GitHub API with query parameters"
echo ""

# Simulate user input
(
  echo "https://api.github.com/search/repositories"  # URL
  echo "1"                                             # GET (choice 1)
  echo "y"                                             # Add headers
  echo "User-Agent"                                    # Header name
  echo "HSP-Client/1.0"                               # Header value
  echo "done"                                          # Done with headers
  echo "y"                                             # Add query params
  echo "q"                                             # Query param: q
  echo "language:go stars:>5000"                       # Value
  echo "sort"                                          # Query param: sort
  echo "stars"                                         # Value
  echo "order"                                         # Query param: order
  echo "desc"                                          # Value
  echo "done"                                          # Done with params
  echo "y"                                             # Pretty print
  echo "y"                                             # Send
) | ./hsp request
