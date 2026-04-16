# Install Guide

To install the ad-blocker system on your machine:

1. **Install the Browser Extension**:
   - Load the unpacked extension from the `extension` directory into Chrome.
   - Note the Extension ID. (If not already `aaaaaaaaaaaaaaaaaa`, you may need to update `installer/install.sh` and repackage).

2. **Run the Installer**:
   - Open the extension popup.
   - You will see a welcoming screen with an auto-detected installation command for your operating system.
   - Copy the command and paste it into your preferred terminal (Bash for macOS/Linux, PowerShell for Windows).
   - Press Enter and follow any prompts. You may be asked for permissions to insert the CA certificate into your system's trust store.

3. **Verify Connection**:
   - Return to the browser extension and click **Check Connection**.
   - If successful, click **Finish Setup**. You are now protected by **Nov**!

## Manual Blocklist Updates
By default, the proxy fetches updates periodically. You can manually force an update from the extension's main interface.
