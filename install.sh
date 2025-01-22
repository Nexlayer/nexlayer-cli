#!/bin/bash

set -e

# ASCII art and colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
cat << "EOF"
 _   _           _                       
| \ | |         | |                      
|  \| | _____  _| | __ _ _   _  ___ _ __ 
| . ` |/ _ \ \/ / |/ _` | | | |/ _ \ '__|
| |\  |  __/>  <| | (_| | |_| |  __/ |   
\_| \_/\___/_/\_\_|\__,_|\__, |\___|_|   
                          __/ |           
                         |___/            
EOF
echo -e "${NC}"

echo "🚀 Installing Nexlayer CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Install the CLI
echo "📦 Installing latest version..."
go install github.com/Nexlayer/nexlayer-cli@latest

# Verify installation
if command -v nexlayer &> /dev/null; then
    echo -e "${GREEN}✨ Nexlayer CLI successfully installed!${NC}"
    echo
    echo "🎯 Get started with:"
    echo "   nexlayer wizard"
    echo
    echo "📚 Learn more:"
    echo "   nexlayer help"
    echo "   nexlayer wizard advanced"
else
    echo -e "${RED}❌ Installation failed${NC}"
    echo "Please try again or visit https://docs.nexlayer.io/installation for help"
    exit 1
fi
