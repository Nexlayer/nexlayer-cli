#!/bin/bash

set -e

# Constants
VERSION="v0.1.0"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
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

# Print colored message
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}" | tee -a "$LOG_FILE"
}

# Check system requirements
check_system_requirements() {
    print_message "$BLUE" "🔍 Checking system requirements..."
    
    # Check Go
    if ! command -v go >/dev/null 2>>"$LOG_FILE"; then
        print_message "$RED" "Error: Go is not installed"
        if confirm "Install Go automatically?"; then
            install_dependency "go" "golang" "https://golang.org/dl/"
        else
            exit 1
        fi
    fi
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    MIN_GO_VERSION="1.23.0"
    if [ "$(printf '%s\n' "$MIN_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_GO_VERSION" ]; then
        print_message "$RED" "Error: Go version must be $MIN_GO_VERSION or higher (found $GO_VERSION)"
        exit 1
    fi

    # Check Git
    if ! command -v git >/dev/null 2>>"$LOG_FILE"; then
        print_message "$RED" "Error: Git is not installed"
        if confirm "Install Git automatically?"; then
            install_dependency "git" "git" "https://git-scm.com/downloads"
        else
            exit 1
        fi
    fi

    # Check disk space (cross-platform fallback)
    AVAILABLE_SPACE=$(df -m . 2>/dev/null | awk 'NR==2 {print $4}' || du -sm . 2>/dev/null | awk '{print $1}')
    if [ "$AVAILABLE_SPACE" -lt 100 ]; then
        print_message "$RED" "Error: Insufficient disk space. Need at least 100MB (found ${AVAILABLE_SPACE}MB)"
        exit 1
    fi
}

# Configure shell environment
configure_shell() {
    print_message "$BLUE" "🔧 Configuring shell environment..."

    # Detect shell configuration file
    SHELL_RC=""
    case "$SHELL" in
        *zsh)
            SHELL_RC="$HOME/.zshrc"
            SHELL_NAME="Zsh"
            ;;
        *bash)
            if [[ "$OSTYPE" == "darwin"* ]]; then
                SHELL_RC="$HOME/.bash_profile"
            else
                SHELL_RC="$HOME/.bashrc"
            fi
            SHELL_NAME="Bash"
            ;;
        *fish)
            SHELL_RC="$HOME/.config/fish/config.fish"
            SHELL_NAME="Fish"
            ;;
        *)
            print_message "$YELLOW" "⚠️ Unsupported shell detected ($SHELL)"
            print_message "$YELLOW" "Please manually add '~/go/bin' to your PATH"
            return
            ;;
    esac

    # Add to PATH if not already present
    if [ -n "$SHELL_RC" ]; then
        mkdir -p "$(dirname "$SHELL_RC")"
        if ! grep -q "~/go/bin" "$SHELL_RC" 2>>"$LOG_FILE"; then
            case "$SHELL_NAME" in
                "Fish")
                    echo "set -x PATH \$PATH ~/go/bin" >> "$SHELL_RC"
                    ;;
                *)
                    echo 'export PATH=$PATH:~/go/bin' >> "$SHELL_RC"
                    ;;
            esac
            print_message "$GREEN" "✅ Added Go bin to PATH in $SHELL_RC"
            
            if confirm "Source $SHELL_RC now?"; then
                case "$SHELL_NAME" in
                    "Fish")
                        source "$SHELL_RC" 2>>"$LOG_FILE" || true
                        ;;
                    *)
                        . "$SHELL_RC" 2>>"$LOG_FILE" || true
                        ;;
                esac
                print_message "$GREEN" "✅ Shell configuration reloaded"
            else
                print_message "$YELLOW" "⚠️ Please restart your terminal or run: source $SHELL_RC"
            fi
        else
            print_message "$GREEN" "✅ PATH already configured in $SHELL_RC"
        fi
    fi
}

# Install dependency (Go or Git)
install_dependency() {
    local cmd=$1
    local pkg=$2
    local url=$3
    print_message "$BLUE" "📥 Installing $cmd..."
    
    if command -v brew >/dev/null 2>>"$LOG_FILE"; then
        brew install "$pkg" >>"$LOG_FILE" 2>&1 &
        show_progress $!
        print_message "$GREEN" "✅ Installed $cmd using Homebrew"
    elif command -v apt >/dev/null 2>>"$LOG_FILE"; then
        sudo apt update >>"$LOG_FILE" 2>&1 && sudo apt install "$pkg" -y >>"$LOG_FILE" 2>&1 &
        show_progress $!
        print_message "$GREEN" "✅ Installed $cmd using apt"
    else
        print_message "$YELLOW" "⚠️ Please install $cmd manually from $url"
        exit 1
    fi
}

# Backup existing installation
backup_existing() {
    if command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        if confirm "Backup existing Nexlayer CLI installation?"; then
            print_message "$BLUE" "📦 Backing up existing installation..."
            BACKUP_DIR="$HOME/.nexlayer/backup/$(date +%Y%m%d_%H%M%S)"
            mkdir -p "$BACKUP_DIR" 2>>"$LOG_FILE"
            cp "$(which nexlayer)" "$BACKUP_DIR/" 2>>"$LOG_FILE"
            print_message "$GREEN" "✅ Backup created at $BACKUP_DIR"
        fi
    fi
}

# Verify installation
verify_installation() {
    print_message "$BLUE" "🔍 Verifying installation..."
    
    if ! command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        print_message "$RED" "❌ Installation failed: nexlayer not found"
        print_message "$YELLOW" "Check $LOG_FILE for details"
        exit 1
    fi

    INSTALLED_VERSION=$(nexlayer version 2>>"$LOG_FILE")
    if [ $? -ne 0 ]; then
        print_message "$RED" "❌ Verification failed"
        exit 1
    fi
    print_message "$GREEN" "✅ Installed: $INSTALLED_VERSION"
}

# User confirmation
confirm() {
    read -p "$1 [y/N]: " response
    case "$response" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}

# Cleanup
cleanup() {
    [ -d "nexlayer-cli" ] && {
        print_message "$BLUE" "🧹 Cleaning up..."
        rm -rf nexlayer-cli 2>>"$LOG_FILE"
    }
}
trap cleanup EXIT

# Print welcome message
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

print_message "$BLUE" "🚀 Installing Nexlayer CLI $VERSION..."

# Main installation
check_system_requirements
backup_existing

print_message "$BLUE" "Choose an installation method:"
echo "1. Simple install (recommended)"
echo "2. Build from source"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" == "1" ]; then
    print_message "$BLUE" "📦 Installing via go install..."
    go install github.com/Nexlayer/nexlayer-cli@latest >>"$LOG_FILE" 2>&1 &
    show_progress $!
    configure_shell
else
    print_message "$BLUE" "📦 Cloning repository..."
    git clone https://github.com/Nexlayer/nexlayer-cli.git >>"$LOG_FILE" 2>&1 &
    show_progress $!
    cd nexlayer-cli || exit 1
    print_message "$BLUE" "🔨 Building from source..."
    go mod download >>"$LOG_FILE" 2>&1 &
    show_progress $!
    go build -o nexlayer . >>"$LOG_FILE" 2>&1 &
    show_progress $!
    print_message "$BLUE" "📥 Installing to /usr/local/bin..."
    sudo mv nexlayer /usr/local/bin/ >>"$LOG_FILE" 2>&1 || {
        print_message "$RED" "❌ Failed to install. Run with sudo or check permissions"
        exit 1
    }
    cd .. || exit 1
fi

verify_installation

print_message "$GREEN" "✨ Nexlayer CLI installed successfully!"
print_message "$BLUE" "\n🎉 Welcome to Nexlayer CLI!"
print_message "$NC" "Deploy full-stack apps in seconds with real-time monitoring.\n"

print_message "$BLUE" "🎯 Core Commands:"
echo "   nexlayer init                # Initialize a new project"
echo "   nexlayer deploy              # Deploy your application"
echo "   nexlayer list               # List all deployments"
echo "   nexlayer info <ns> <app>    # Get deployment info"
echo "   nexlayer domain set         # Configure custom domain"
echo "   nexlayer watch              # Watch for changes"

print_message "$BLUE" "\n🔧 Configuration Commands:"
echo "   nexlayer completion         # Generate shell completions"

print_message "$BLUE" "\n🐚 Shell Completion Setup:"
echo "   nexlayer completion bash > ~/.bash_completion"
echo "   nexlayer completion zsh > ${fpath[1]}/_nexlayer"
echo "   nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish"

print_message "$BLUE" "\n📚 Learn more:"
echo "   nexlayer help              # Show detailed help"
echo "   nexlayer --version         # Show version info"

print_message "$BLUE" "\n💡 For developers:"
echo "   Run 'make setup' in the repo to set up the dev environment."
echo "   See contribution guidelines: https://github.com/Nexlayer/nexlayer-cli/blob/main/CONTRIBUTING.md"

if confirm "Open Nexlayer CLI docs?"; then
    if command -v xdg-open >/dev/null 2>>"$LOG_FILE"; then
        xdg-open https://nexlayer.dev/docs >>"$LOG_FILE" 2>&1 &
    elif command -v open >/dev/null 2>>"$LOG_FILE"; then
        open https://nexlayer.dev/docs >>"$LOG_FILE" 2>&1 &
    else
        print_message "$BLUE" "Visit: https://nexlayer.dev/docs"
    fi
fi