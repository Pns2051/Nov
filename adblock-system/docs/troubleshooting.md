# Troubleshooting

## Certificate NOT Trusted
The installation script attempts to install the CA certificate. If you see certificate warnings in Chrome, you may need to install the certificate manually.
- The certificate is generated at: `~/.adblock-proxy/ca-cert.pem`
- On macOS, open Keychain Access, import it into `System` and explicitly set it to "Always Trust".
- On Windows, import into `Trusted Root Certification Authorities`.

## Proxy Not Running
Check if the proxy is running in the background. 
- macOS: Activity Monitor (look for `adblock-proxy`)
- Windows: Task Manager
- Linux: `systemctl --user status adblock-proxy`

## Extension Not Connecting
If the extension displays "Disconnected", the native messaging bridge is failing. Ensure that `com.adblock.proxy.json` is correctly installed in your browser's native messaging hosts directory and points to the right path.
