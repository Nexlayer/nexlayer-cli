#!/bin/bash

set -e

# Version and colors
VERSION="v0.1.0"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Progress indicator
show_progress() {
    local pid=$1
    local delay=0.1
    local spinstr='|/-\'
    while ps -p $pid > /dev/null; do
        local temp=${spinstr#?}
        printf " [%c]  " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
}

# System requirements check
check_system_requirements() {
    # Check minimum Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if ! command -v awk &> /dev/null; then
        echo -e "${RED}Error: awk is not installed${NC}"
        exit 1
    }
    
    MIN_GO_VERSION="1.23.0"
    if [ "$(printf '%s\n' "$MIN_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_GO_VERSION" ]; then
        echo -e "${RED}Error: Go version must be $MIN_GO_VERSION or higher${NC}"
        exit 1
    }

    # Check available disk space (minimum 100MB)
    if command -v df &> /dev/null; then
        AVAILABLE_SPACE=$(df -m . | awk 'NR==2 {print $4}')
        if [ "$AVAILABLE_SPACE" -lt 100 ]; then
            echo -e "${RED}Error: Insufficient disk space. Need at least 100MB${NC}"
            exit 1
        fi
    fi
}

# Backup existing installation
backup_existing() {
    if command -v nexlayer &> /dev/null; then
        echo "📦 Backing up existing installation..."
        BACKUP_DIR="$HOME/.nexlayer/backup/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        cp "$(which nexlayer)" "$BACKUP_DIR/"
        echo "✅ Backup created at $BACKUP_DIR"
    fi
}

# Configure shell environment
configure_shell() {
    # Detect shell
    SHELL_RC=""
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        SHELL_RC="$HOME/.bashrc"
    fi

    if [ -n "$SHELL_RC" ]; then
        if ! grep -q 'export PATH=$PATH:~/go/bin' "$SHELL_RC"; then
            echo 'export PATH=$PATH:~/go/bin' >> "$SHELL_RC"
            echo "✅ Added Go bin to PATH in $SHELL_RC"
        fi
    fi
}

# Verify installation
verify_installation() {
    if ! command -v nexlayer &> /dev/null; then
        echo -e "${RED}❌ Installation failed: nexlayer command not found${NC}"
        exit 1
    fi

    # Verify version
    INSTALLED_VERSION=$(nexlayer version)
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Installation verification failed${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Verified installation: $INSTALLED_VERSION${NC}"
}

# Cleanup function
cleanup() {
    if [ -d "nexlayer-cli" ]; then
        echo "🧹 Cleaning up temporary files..."
        rm -rf nexlayer-cli
    fi
}
trap cleanup EXIT

# Display ASCII art logo
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

echo "🚀 Installing Nexlayer CLI ${VERSION}..."

# Check dependencies
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: Git is not installed${NC}"
    echo "Please install Git from https://git-scm.com/downloads"
    exit 1
fi

# Check system requirements
check_system_requirements

# Backup existing installation
backup_existing

# Installation method prompt
echo "Choose an installation method:"
echo "1. Simple install (recommended)"
echo "2. Build from source"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" == "1" ]; then
    echo "📦 Installing via go install..."
    go install github.com/Nexlayer/nexlayer-cli@latest &
    show_progress $!
    echo "✅ Installation complete!"
    
    # Configure shell
    configure_shell
else
    echo "📦 Cloning Nexlayer CLI repository..."
    git clone https://github.com/Nexlayer/nexlayer-cli.git &
    show_progress $!
    cd nexlayer-cli
    
    echo "🔨 Building from source..."
    go mod download &
    show_progress $!
    go build -o nexlayer . &
    show_progress $!
    
    echo "📥 Installing to /usr/local/bin..."
    if ! sudo mv nexlayer /usr/local/bin/; then
        echo -e "${RED}❌ Failed to install to /usr/local/bin${NC}"
        echo "Try running with sudo or install manually."
        exit 1
    fi
    cd ..
fi

# Verify installation
verify_installation

echo -e "${GREEN}✨ Nexlayer CLI successfully installed!${NC}"
echo
echo "🎉 Welcome to Nexlayer CLI! 🎉"
echo "Deploy full-stack apps in seconds with AI-powered detection and real-time monitoring."
echo
echo "🎯 Core Commands:"
echo "   nexlayer init                # Initialize a new project"
echo "   nexlayer deploy              # Deploy your application"
echo "   nexlayer list               # List all deployments"
echo "   nexlayer info <ns> <app>    # Get deployment info"
echo
echo "🤖 AI Commands:"
echo "   nexlayer ai detect          # Detect project type with AI"
echo "   nexlayer ai generate        # Generate deployment template"
echo
echo "🔧 Configuration Commands:"
echo "   nexlayer domain set         # Configure custom domain"
echo "   nexlayer feedback           # Send feedback"
echo "   nexlayer completion         # Generate shell completions"
echo
echo "🐚 Shell Completion Setup:"
echo "   nexlayer completion bash > ~/.bash_completion"
echo "   nexlayer completion zsh > ${fpath[1]}/_nexlayer"
echo "   nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish"
echo
echo "📚 Learn more:"
echo "   nexlayer help              # Show detailed help"
echo "   nexlayer --version         # Show version info"
echo
echo "💡 For developers:"
echo "   Run 'make setup' in the repo to set up the dev environment."
echo "   See contribution guidelines: https://github.com/Nexlayer/nexlayer-cli/blob/main/CONTRIBUTING.md"

# Prompt to open documentation
read -p "Would you like to open the Nexlayer CLI docs? (y/n): " answer
if [[ "$answer" == "y" || "$answer" == "Y" ]]; then
    if command -v xdg-open &> /dev/null; then
        xdg-open https://nexlayer.dev/docs
    elif command -v open &> /dev/null; then
        open https://nexlayer.dev/docs
    else
        echo "Couldn't open automatically. Visit: https://nexlayer.dev/docs"
    fi
fi