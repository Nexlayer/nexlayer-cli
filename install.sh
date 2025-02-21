#!/bin/bash

set -e

# Constants
VERSION="v0.1.0"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'
LOG_FILE="$HOME/nexlayer_install.log"

# Progress indicator
show_progress() {
    local pid=$1
    local delay=0.1
    local spinstr='|/-\'
    while ps -p "$pid" > /dev/null 2>>"$LOG_FILE"; do
        local temp=${spinstr#?}
        printf " [%c]  " "$spinstr"
        spinstr=$temp${spinstr%"$temp"}
        sleep "$delay"
        printf "\b\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
}

# Check system requirements
check_system_requirements() {
    echo "🔍 Checking system requirements..." | tee -a "$LOG_FILE"
    
    # Check Go
    if ! command -v go >/dev/null 2>>"$LOG_FILE"; then
        echo -e "${RED}Error: Go is not installed${NC}" | tee -a "$LOG_FILE"
        if confirm "Install Go automatically?"; then
            install_dependency "go" "golang" "https://golang.org/dl/"
        else
            exit 1
        fi
    fi
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    MIN_GO_VERSION="1.23.0"
    if [ "$(printf '%s\n' "$MIN_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_GO_VERSION" ]; then
        echo -e "${RED}Error: Go version must be $MIN_GO_VERSION or higher (found $GO_VERSION)${NC}" | tee -a "$LOG_FILE"
        exit 1
    fi

    # Check Git
    if ! command -v git >/dev/null 2>>"$LOG_FILE"; then
        echo -e "${RED}Error: Git is not installed${NC}" | tee -a "$LOG_FILE"
        if confirm "Install Git automatically?"; then
            install_dependency "git" "git" "https://git-scm.com/downloads"
        else
            exit 1
        fi
    fi

    # Check disk space (cross-platform fallback)
    AVAILABLE_SPACE=$(df -m . 2>/dev/null | awk 'NR==2 {print $4}' || du -sm . 2>/dev/null | awk '{print $1}')
    if [ "$AVAILABLE_SPACE" -lt 100 ]; then
        echo -e "${RED}Error: Insufficient disk space. Need at least 100MB (found ${AVAILABLE_SPACE}MB)${NC}" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Install dependency (Go or Git)
install_dependency() {
    local cmd=$1
    local pkg=$2
    local url=$3
    echo "📥 Installing $cmd..." | tee -a "$LOG_FILE"
    if command -v brew >/dev/null 2>>"$LOG_FILE"; then
        brew install "$pkg" >>"$LOG_FILE" 2>&1 &
        show_progress $!
    elif command -v apt >/dev/null 2>>"$LOG_FILE"; then
        sudo apt update >>"$LOG_FILE" 2>&1 && sudo apt install "$pkg" -y >>"$LOG_FILE" 2>&1 &
        show_progress $!
    else
        echo "Please install $cmd manually from $url" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Backup existing installation
backup_existing() {
    if command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        if confirm "Backup existing Nexlayer CLI installation?"; then
            echo "📦 Backing up existing installation..." | tee -a "$LOG_FILE"
            BACKUP_DIR="$HOME/.nexlayer/backup/$(date +%Y%m%d_%H%M%S)"
            mkdir -p "$BACKUP_DIR" 2>>"$LOG_FILE"
            cp "$(which nexlayer)" "$BACKUP_DIR/" 2>>"$LOG_FILE"
            echo "✅ Backup created at $BACKUP_DIR" | tee -a "$LOG_FILE"
        fi
    fi
}

# Configure shell environment
configure_shell() {
    SHELL_RC=""
    case "$SHELL" in
        *zsh) SHELL_RC="$HOME/.zshrc" ;;
        *bash) SHELL_RC="$HOME/.bashrc" ;;
        *fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
        *) echo "⚠️ Unsupported shell detected" | tee -a "$LOG_FILE"; return ;;
    esac

    if [ -n "$SHELL_RC" ]; then
        local path_line
        [ "$SHELL" = *fish ] && path_line="set -x PATH \$PATH ~/go/bin" || path_line="export PATH=\$PATH:~/go/bin"
        if ! grep -q "~/go/bin" "$SHELL_RC" 2>>"$LOG_FILE"; then
            echo "$path_line" >> "$SHELL_RC"
            echo "✅ Added Go bin to PATH in $SHELL_RC" | tee -a "$LOG_FILE"
            if confirm "Source $SHELL_RC now?"; then
                source "$SHELL_RC" 2>>"$LOG_FILE" || echo "⚠️ Please restart your terminal" | tee -a "$LOG_FILE"
            fi
        fi
    fi
}

# Verify installation
verify_installation() {
    if ! command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        echo -e "${RED}❌ Installation failed: nexlayer not found${NC}" | tee -a "$LOG_FILE"
        echo "Check $LOG_FILE for details" | tee -a "$LOG_FILE"
        exit 1
    fi
    INSTALLED_VERSION=$(nexlayer version 2>>"$LOG_FILE")
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Verification failed${NC}" | tee -a "$LOG_FILE"
        exit 1
    fi
    echo -e "${GREEN}✅ Installed: $INSTALLED_VERSION${NC}" | tee -a "$LOG_FILE"
}

# Cleanup
cleanup() {
    [ -d "nexlayer-cli" ] && { echo "🧹 Cleaning up..." | tee -a "$LOG_FILE"; rm -rf nexlayer-cli 2>>"$LOG_FILE"; }
}
trap cleanup EXIT

# User confirmation
confirm() {
    read -p "$1 [y/N]: " response
    case "$response" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}

# ASCII art
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

echo "🚀 Installing Nexlayer CLI $VERSION..." | tee "$LOG_FILE"

# Main installation
check_system_requirements
backup_existing

echo "Choose an installation method:" | tee -a "$LOG_FILE"
echo "1. Simple install (recommended)" | tee -a "$LOG_FILE"
echo "2. Build from source" | tee -a "$LOG_FILE"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" == "1" ]; then
    echo "📦 Installing via go install..." | tee -a "$LOG_FILE"
    go install github.com/Nexlayer/nexlayer-cli@latest >>"$LOG_FILE" 2>&1 &
    show_progress $!
    configure_shell
else
    echo "📦 Cloning repository..." | tee -a "$LOG_FILE"
    git clone https://github.com/Nexlayer/nexlayer-cli.git >>"$LOG_FILE" 2>&1 &
    show_progress $!
    cd nexlayer-cli || exit 1
    echo "🔨 Building from source..." | tee -a "$LOG_FILE"
    go mod download >>"$LOG_FILE" 2>&1 &
    show_progress $!
    go build -o nexlayer . >>"$LOG_FILE" 2>&1 &
    show_progress $!
    echo "📥 Installing to /usr/local/bin..." | tee -a "$LOG_FILE"
    sudo mv nexlayer /usr/local/bin/ >>"$LOG_FILE" 2>&1 || {
        echo -e "${RED}❌ Failed to install. Run with sudo or check permissions${NC}" | tee -a "$LOG_FILE"
        exit 1
    }
    cd .. || exit 1
fi

verify_installation

echo -e "${GREEN}✨ Nexlayer CLI installed successfully!${NC}" | tee -a "$LOG_FILE"
echo -e "\n🎉 Welcome to Nexlayer CLI!\nDeploy full-stack apps in seconds with AI-powered detection and real-time monitoring.\n" | tee -a "$LOG_FILE"
echo "🎯 Core Commands:" | tee -a "$LOG_FILE"
echo "   nexlayer init                # Initialize a new project" | tee -a "$LOG_FILE"
echo "   nexlayer deploy              # Deploy your application" | tee -a "$LOG_FILE"
echo "   nexlayer list               # List all deployments" | tee -a "$LOG_FILE"
echo "   nexlayer info <ns> <app>    # Get deployment info" | tee -a "$LOG_FILE"
echo -e "\n🤖 AI Commands:" | tee -a "$LOG_FILE"
echo "   nexlayer ai detect          # Detect project type with AI" | tee -a "$LOG_FILE"
echo "   nexlayer ai generate        # Generate deployment template" | tee -a "$LOG_FILE"
echo -e "\n🔧 Configuration Commands:" | tee -a "$LOG_FILE"
echo "   nexlayer domain set         # Configure custom domain" | tee -a "$LOG_FILE"
echo "   nexlayer feedback           # Send feedback" | tee -a "$LOG_FILE"
echo "   nexlayer completion         # Generate shell completions" | tee -a "$LOG_FILE"
echo -e "\n🐚 Shell Completion Setup:" | tee -a "$LOG_FILE"
echo "   nexlayer completion bash > ~/.bash_completion" | tee -a "$LOG_FILE"
echo "   nexlayer completion zsh > ${fpath[1]}/_nexlayer" | tee -a "$LOG_FILE"
echo "   nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish" | tee -a "$LOG_FILE"
echo -e "\n📚 Learn more:" | tee -a "$LOG_FILE"
echo "   nexlayer help              # Show detailed help" | tee -a "$LOG_FILE"
echo "   nexlayer --version         # Show version info" | tee -a "$LOG_FILE"
echo -e "\n💡 For developers:" | tee -a "$LOG_FILE"
echo "   Run 'make setup' in the repo to set up the dev environment." | tee -a "$LOG_FILE"
echo "   See contribution guidelines: https://github.com/Nexlayer/nexlayer-cli/blob/main/CONTRIBUTING.md" | tee -a "$LOG_FILE"

if confirm "Open Nexlayer CLI docs?"; then
    if command -v xdg-open >/dev/null 2>>"$LOG_FILE"; then
        xdg-open https://nexlayer.dev/docs >>"$LOG_FILE" 2>&1 &
    elif command -v open >/dev/null 2>>"$LOG_FILE"; then
        open https://nexlayer.dev/docs >>"$LOG_FILE" 2>&1 &
    else
        echo "Visit: https://nexlayer.dev/docs" | tee -a "$LOG_FILE"
    fi
fi