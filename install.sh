#!/bin/bash
set -e

GITHUB_USER="Pns2051"
VERSION="latest"
INSTALL_DIR="$HOME/.adblock-proxy"
BINARY_NAME="adblock-proxy"
EXTENSION_ID="aaaaaaaaaaaaaaaaaa"

# Detect OS
OS=$(uname -s)
if [ "$OS" = "Darwin" ]; then
    OS="darwin"
elif [ "$OS" = "Linux" ]; then
    OS="linux"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# Detect Architecture
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
else
    echo "Unsupported Architecture: $ARCH"
    exit 1
fi

BINARY_FILENAME="$BINARY_NAME-$OS-$ARCH"

download() {
    local url_jsdelivr="https://cdn.jsdelivr.net/gh/$GITHUB_USER/Nov@main/$1"
    local url_github="https://github.com/$GITHUB_USER/Nov/releases/$VERSION/download/$1"
    echo "Downloading $1..."
    curl -fsSL --connect-timeout 10 --max-time 60 "$url_jsdelivr" -o "$2" || \
    curl -fsSL --connect-timeout 10 --max-time 60 "$url_github" -o "$2"
}

mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

download "dist/$BINARY_FILENAME" "$BINARY_NAME"
chmod +x "$BINARY_NAME"
download "blocklist/blocklist.txt" "blocklist.txt" || echo "Could not download blocklist, will generate later."

echo "Generating CA Certificate..."
./$BINARY_NAME -generate-ca

if [ "$OS" = "darwin" ]; then
    # Native messaging manifest
    NATIVE_DIR="$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
    mkdir -p "$NATIVE_DIR"
    cat <<EOF > "$NATIVE_DIR/com.adblock.proxy.json"
{
  "name": "com.adblock.proxy",
  "description": "AdBlocker System Proxy",
  "path": "$INSTALL_DIR/$BINARY_NAME",
  "type": "stdio",
  "allowed_origins": ["chrome-extension://$EXTENSION_ID/"]
}
EOF

    # Chromium manifest
    CHROMIUM_NATIVE_DIR="$HOME/Library/Application Support/Chromium/NativeMessagingHosts"
    mkdir -p "$CHROMIUM_NATIVE_DIR"
    cp "$NATIVE_DIR/com.adblock.proxy.json" "$CHROMIUM_NATIVE_DIR/"

    # Set system proxy
    echo "Configuring system proxy for active network service..."
    SERVICE=$(networksetup -listnetworkserviceorder | grep -B 1 $(route -n get default | grep interface | awk '{print $2}') | head -n 1 | sed 's/([^)]*)//g' | xargs)
    if [ -n "$SERVICE" ]; then
        networksetup -setwebproxy "$SERVICE" 127.0.0.1 8080
        networksetup -setsecurewebproxy "$SERVICE" 127.0.0.1 8080
    else
        echo "Could not auto-detect active network service. Please set HTTP/HTTPS proxy to 127.0.0.1:8080 manually."
    fi

    # Install CA
    echo "Installing CA Certificate. You may be prompted for your password."
    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "$INSTALL_DIR/ca-cert.pem"

    # LaunchAgent
    PLIST_PATH="$HOME/Library/LaunchAgents/com.user.adblockproxy.plist"
    cat <<EOF > "$PLIST_PATH"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.adblockproxy</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/$BINARY_NAME</string>
        <string>-mode=proxy</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$INSTALL_DIR/log.txt</string>
    <key>StandardErrorPath</key>
    <string>$INSTALL_DIR/error.txt</string>
</dict>
</plist>
EOF
    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    launchctl load "$PLIST_PATH"
    launchctl start com.user.adblockproxy

elif [ "$OS" = "linux" ]; then
    # Native messaging manifest
    NATIVE_DIR="$HOME/.config/google-chrome/NativeMessagingHosts"
    mkdir -p "$NATIVE_DIR"
    cat <<EOF > "$NATIVE_DIR/com.adblock.proxy.json"
{
  "name": "com.adblock.proxy",
  "description": "AdBlocker System Proxy",
  "path": "$INSTALL_DIR/$BINARY_NAME",
  "type": "stdio",
  "allowed_origins": ["chrome-extension://$EXTENSION_ID/"]
}
EOF

    CHROMIUM_NATIVE_DIR="$HOME/.config/chromium/NativeMessagingHosts"
    mkdir -p "$CHROMIUM_NATIVE_DIR"
    cp "$NATIVE_DIR/com.adblock.proxy.json" "$CHROMIUM_NATIVE_DIR/"

    # Proxy settings
    if command -v gsettings >/dev/null 2>&1; then
        gsettings set org.gnome.system.proxy mode 'manual'
        gsettings set org.gnome.system.proxy.http host '127.0.0.1'
        gsettings set org.gnome.system.proxy.http port 8080
        gsettings set org.gnome.system.proxy.https host '127.0.0.1'
        gsettings set org.gnome.system.proxy.https port 8080
        echo "Gnome proxy configured."
    else
        echo "Please configure your system proxy manually to 127.0.0.1:8080."
    fi

    # Install CA
    echo "Installing CA Certificate. You may be prompted for your password."
    if [ -d "/usr/local/share/ca-certificates/" ]; then
        sudo cp "$INSTALL_DIR/ca-cert.pem" /usr/local/share/ca-certificates/adblock-proxy.crt
        sudo update-ca-certificates
    elif [ -d "/usr/share/pki/trust/anchors/" ]; then
        sudo cp "$INSTALL_DIR/ca-cert.pem" /usr/share/pki/trust/anchors/adblock-proxy.pem
        sudo update-ca-trust
    else
        echo "Could not identify CA cert dir. Please install $INSTALL_DIR/ca-cert.pem manually."
    fi

    # Systemd Service
    SYSTEMD_DIR="$HOME/.config/systemd/user"
    mkdir -p "$SYSTEMD_DIR"
    cat <<EOF > "$SYSTEMD_DIR/adblock-proxy.service"
[Unit]
Description=AdBlocker System Proxy

[Service]
ExecStart=$INSTALL_DIR/$BINARY_NAME -mode=proxy
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
EOF
    systemctl --user daemon-reload
    systemctl --user enable adblock-proxy.service
    systemctl --user restart adblock-proxy.service
fi

echo ""
echo "Installation complete!"
echo "Return to the extension and click 'Check Connection'."
