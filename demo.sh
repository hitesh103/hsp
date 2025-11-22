#!/bin/bash

# Demo script for HSP interactive request builder
# This script demonstrates the workflow

echo "=========================================="
echo "HSP Interactive Request Builder Demo"
echo "=========================================="
echo ""
echo "Let's make a GET request to the GitHub API"
echo ""

# Simulate user input for a GitHub API request
(
  echo "https://api.github.com/users/golang"  # URL
  echo "1"                                     # GET (choice 1)
  echo "y"                                     # Add headers
  echo "User-Agent"                            # Header name
  echo "HSP-Client/1.0"                        # Header value
  echo "done"                                  # Done with headers
  echo "n"                                     # No query params
  echo "y"                                     # Pretty print
  echo "y"                                     # Send
) | ./hsp request

echo ""
echo "=========================================="
echo "Demo Complete! Request saved to ~/.hsp/history/"
echo "=========================================="
