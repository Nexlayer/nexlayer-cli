#!/bin/bash

set -eo pipefail  # Exit on error and pipe failures

# Constants
VERSION="v0.1.0-alpha.9"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'
LOG_FILE="$HOME/nexlayer_install.log"
MIN_GO_VERSION="1.23.0"
MIN_DISK_SPACE=100  # MB
DEFAULT_GLOBAL_DIR="$HOME/.local/bin"
LOCAL_DIR="./bin"

# Create log file directory if it doesn't exist
mkdir -p "$(dirname "$LOG_FILE")" 2>/dev/null || true
echo "Nexlayer CLI installation started at $(date)" > "$LOG_FILE"

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

    # Check OS compatibility
    OS=$(uname -s)
    case "$OS" in
        Darwin) OS_NAME="macOS" ;;
        Linux) OS_NAME="Linux" ;;
        MINGW*|MSYS*|CYGWIN*) OS_NAME="Windows" ;;
        *)
            print_message "$YELLOW" "⚠️ Unsupported OS: $OS"
            print_message "$YELLOW" "Installation may not work correctly."
            ;;
    esac
    print_message "$GREEN" "✅ Detected OS: $OS_NAME"

    # Check Go
    if ! command -v go >/dev/null 2>>"$LOG_FILE"; then
        print_message "$RED" "Error: Go is not installed"
        print_message "$YELLOW" "Please install Go from https://golang.org/dl/"
        if confirm "Install Go automatically?"; then
            install_dependency "go" "golang" "https://golang.org/dl/"
        else
            exit 1
        fi
    fi
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if ! check_version "$GO_VERSION" "$MIN_GO_VERSION"; then
        print_message "$RED" "Error: Go version must be $MIN_GO_VERSION or higher (found $GO_VERSION)"
        print_message "$YELLOW" "Please update Go from https://golang.org/dl/"
        exit 1
    fi
    print_message "$GREEN" "✅ Go version: $GO_VERSION"

    # Check Git
    if ! command -v git >/dev/null 2>>"$LOG_FILE"; then
        print_message "$RED" "Error: Git is not installed"
        print_message "$YELLOW" "Please install Git from https://git-scm.com/downloads"
        if confirm "Install Git automatically?"; then
            install_dependency "git" "git" "https://git-scm.com/downloads"
        else
            exit 1
        fi
    fi
    print_message "$GREEN" "✅ Git installed"

    # Check disk space
    check_disk_space
}

# Check version numbers
check_version() {
    local version=$1
    local required=$2
    version="${version#v}"
    required="${required#v}"
    if [ "$(printf '%s\n' "$required" "$version" | sort -V | head -n1)" != "$required" ]; then
        return 1
    else
        return 0
    fi
}

# Check available disk space
check_disk_space() {
    local available
    if [ "$OS" = "Darwin" ] || [ "$OS" = "Linux" ]; then
        available=$(df -m . | awk 'NR==2 {print $4}')
    else
        available=$(du -sm . 2>/dev/null | awk '{print $1}')
    fi
    if [ -z "$available" ] || [ "$available" -lt "$MIN_DISK_SPACE" ]; then
        print_message "$RED" "Error: Need at least ${MIN_DISK_SPACE}MB (found ${available:-unknown}MB)"
        exit 1
    fi
    print_message "$GREEN" "✅ Disk space: ${available}MB available"
}

# Configure shell environment
configure_shell() {
    local install_dir=$1
    print_message "$BLUE" "🔧 Configuring shell environment..."

    detect_shell_config
    if [ -n "$SHELL_RC" ]; then
        mkdir -p "$(dirname "$SHELL_RC")" 2>/dev/null || true
        if ! grep -q "$install_dir" "$SHELL_RC" 2>>"$LOG_FILE"; then
            case "$SHELL_NAME" in
                "Fish") echo "set -x PATH \$PATH $install_dir" >> "$SHELL_RC" ;;
                *) echo "export PATH=\$PATH:$install_dir" >> "$SHELL_RC" ;;
            esac
            print_message "$GREEN" "✅ Added $install_dir to PATH in $SHELL_RC"
            if confirm "Source $SHELL_RC now?"; then
                case "$SHELL_NAME" in
                    "Fish") source "$SHELL_RC" 2>>"$LOG_FILE" || true ;;
                    *) . "$SHELL_RC" 2>>"$LOG_FILE" || true ;;
                esac
                print_message "$GREEN" "✅ Shell reloaded"
            else
                print_message "$YELLOW" "⚠️ Run: source $SHELL_RC or restart your terminal"
            fi
        else
            print_message "$GREEN" "✅ PATH already configured in $SHELL_RC"
        fi
    else
        print_message "$YELLOW" "⚠️ Add $install_dir to your PATH manually"
    fi
}

# Detect shell configuration file
detect_shell_config() {
    SHELL_RC=""
    case "$SHELL" in
        *zsh) SHELL_RC="$HOME/.zshrc"; SHELL_NAME="Zsh" ;;
        *bash)
            if [[ "$OSTYPE" == "darwin"* ]]; then
                SHELL_RC="$HOME/.bash_profile"
            else
                SHELL_RC="$HOME/.bashrc"
            fi
            SHELL_NAME="Bash" ;;
        *fish) SHELL_RC="$HOME/.config/fish/config.fish"; SHELL_NAME="Fish" ;;
        *)
            print_message "$YELLOW" "⚠️ Unsupported shell: $SHELL"
            return ;;
    esac
    print_message "$GREEN" "✅ Detected shell: $SHELL_NAME"
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
    elif command -v apt-get >/dev/null 2>>"$LOG_FILE"; then
        sudo apt-get update >>"$LOG_FILE" 2>&1 && sudo apt-get install "$pkg" -y >>"$LOG_FILE" 2>&1 &
        show_progress $!
    elif command -v dnf >/dev/null 2>>"$LOG_FILE"; then
        sudo dnf install "$pkg" -y >>"$LOG_FILE" 2>&1 &
        show_progress $!
    elif command -v yum >/dev/null 2>>"$LOG_FILE"; then
        sudo yum install "$pkg" -y >>"$LOG_FILE" 2>&1 &
        show_progress $!
    elif command -v pacman >/dev/null 2>>"$LOG_FILE"; then
        sudo pacman -S --noconfirm "$pkg" >>"$LOG_FILE" 2>&1 &
        show_progress $!
    else
        print_message "$YELLOW" "⚠️ Install $cmd manually from $url"
        exit 1
    fi
    print_message "$GREEN" "✅ Installed $cmd"
}

# Backup existing installation
backup_existing() {
    if command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        if confirm "Backup existing Nexlayer CLI?"; then
            print_message "$BLUE" "📦 Backing up..."
            BACKUP_DIR="$HOME/.nexlayer/backup/$(date +%Y%m%d_%H%M%S)"
            mkdir -p "$BACKUP_DIR" 2>>"$LOG_FILE"
            cp "$(which nexlayer)" "$BACKUP_DIR/" 2>>"$LOG_FILE"
            print_message "$GREEN" "✅ Backup at $BACKUP_DIR"
        fi
    fi
}

# Verify installation
verify_installation() {
    local install_dir=$1
    print_message "$BLUE" "🔍 Verifying installation..."

    local nexlayer_path="$install_dir/nexlayer"
    if [ ! -x "$nexlayer_path" ] && ! command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        print_message "$RED" "❌ Failed: nexlayer not found in $install_dir or PATH"
        print_message "$YELLOW" "Fixes:"
        print_message "$YELLOW" "1. Add $install_dir to PATH"
        print_message "$YELLOW" "2. Run: source $SHELL_RC"
        print_message "$YELLOW" "3. Check $LOG_FILE"
        exit 1
    fi

    # Try to get version, but don't fail if version command is not available
    INSTALLED_VERSION=$("$nexlayer_path" version 2>>"$LOG_FILE" || nexlayer version 2>>"$LOG_FILE" || echo "Unknown version")
    if [ -x "$nexlayer_path" ] || command -v nexlayer >/dev/null 2>>"$LOG_FILE"; then
        print_message "$GREEN" "✅ Installed: $INSTALLED_VERSION"
    else
        print_message "$RED" "❌ Failed: nexlayer command error"
        print_message "$YELLOW" "Check $LOG_FILE"
        exit 1
    fi
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

# Welcome message
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

print_message "$BLUE" "Choose installation method:"
echo "1. Simple install (recommended) - Uses 'go install'"
echo "2. Build from source - Clones and builds manually"
read -p "Enter choice (1 or 2): " choice

print_message "$BLUE" "Where to install?"
echo "1. Globally ($DEFAULT_GLOBAL_DIR)"
echo "2. Locally in project ($LOCAL_DIR)"
read -p "Enter choice (1 or 2): " location_choice

if [ "$location_choice" == "1" ]; then
    INSTALL_DIR="$DEFAULT_GLOBAL_DIR"
else
    INSTALL_DIR="$LOCAL_DIR"
fi
mkdir -p "$INSTALL_DIR" 2>>"$LOG_FILE" || {
    print_message "$RED" "❌ Failed to create directory $INSTALL_DIR"
    print_message "$YELLOW" "Try running with sudo or choose a different location"
    exit 1
}
chmod 755 "$INSTALL_DIR" 2>>"$LOG_FILE" || {
    print_message "$YELLOW" "⚠️ Warning: Could not set permissions on $INSTALL_DIR"
}

if [ "$choice" == "1" ]; then
    print_message "$BLUE" "📦 Installing via go install..."
    GO111MODULE=on go install github.com/Nexlayer/nexlayer-cli@latest >>"$LOG_FILE" 2>&1 &
    show_progress $!
    if [ "$location_choice" == "2" ]; then
        cp "$HOME/go/bin/nexlayer" "$INSTALL_DIR/" 2>>"$LOG_FILE" || {
            print_message "$RED" "❌ Failed to copy binary to $INSTALL_DIR"
            print_message "$YELLOW" "Using binary from $HOME/go/bin instead"
            INSTALL_DIR="$HOME/go/bin"
        }
    else
        INSTALL_DIR="$HOME/go/bin"
    fi
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
    mv nexlayer "$INSTALL_DIR/" 2>>"$LOG_FILE" || {
        print_message "$RED" "❌ Failed to install to $INSTALL_DIR"
        exit 1
    }
    cd .. || exit 1
fi

configure_shell "$INSTALL_DIR"
verify_installation "$INSTALL_DIR"

print_message "$GREEN" "✨ Nexlayer CLI installed successfully!"
print_message "$BLUE" "\n🎉 Welcome to Nexlayer CLI!"
if [ "$location_choice" == "2" ]; then
    print_message "$NC" "Run it from your project with: $INSTALL_DIR/nexlayer"
else
    print_message "$NC" "Deploy apps globally with: nexlayer"
fi

print_message "$BLUE" "🎯 Core Commands:"
echo "   nexlayer init                # Initialize a new project"
echo "   nexlayer deploy              # Deploy your application"
echo "   nexlayer list                # List all deployments"
echo "   nexlayer info <ns> <app>     # Get deployment info"
echo "   nexlayer domain set          # Configure custom domain"
echo "   nexlayer watch               # Watch for changes"

print_message "$BLUE" "\n📚 Learn more:"
echo "   nexlayer help                # Show detailed help"
echo "   nexlayer --version           # Show version info"

echo "Installation completed at $(date)" >> "$LOG_FILE"

if confirm "Open Nexlayer CLI docs?"; then
    if command -v xdg-open >/dev/null 2>>"$LOG_FILE"; then
        xdg-open https://nexlayer.dev/docs >>"$LOG_FILE" 2>&1 &
    elif command -v open >/dev/null 2>>"$LOG_FILE"; then
        open https://nexlayer.dev/docs >>"$LOG_FILE" 2>&1 &
    else
        print_message "$BLUE" "Visit: https://nexlayer.dev/docs"
    fi
fi