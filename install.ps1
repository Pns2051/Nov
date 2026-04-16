$ErrorActionPreference = "Stop"

$GITHUB_USER = "Pns2051"
$VERSION = "latest"
$INSTALL_DIR = "$env:LOCALAPPDATA\AdblockProxy"
$BINARY_NAME = "adblock-proxy-windows-amd64.exe"
$EXTENSION_ID = "aaaaaaaaaaaaaaaaaa"

function Download-File {
    param([string]$UrlRaw, [string]$UrlGithub, [string]$Destination)
    try {
        # Try GitHub Raw first
        Invoke-WebRequest -Uri $UrlRaw -OutFile $Destination -UseBasicParsing -TimeoutSec 60
    } catch {
        Write-Host "Primary download failed (GitHub Raw), trying fallback (GitHub Releases)..."
        try {
            Invoke-WebRequest -Uri $UrlGithub -OutFile $Destination -UseBasicParsing -TimeoutSec 60
        } catch {
             Write-Host "All downloads failed."
             throw
        }
    }
}

if (-not (Test-Path -Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}
Set-Location -Path $INSTALL_DIR

Write-Host "Downloading proxy binary..."
Download-File -UrlRaw "https://raw.githubusercontent.com/$GITHUB_USER/Nov/main/dist/$BINARY_NAME" -UrlGithub "https://github.com/$GITHUB_USER/Nov/releases/$VERSION/download/$BINARY_NAME" -Destination "$INSTALL_DIR\$BINARY_NAME"

Write-Host "Downloading blocklist..."
try {
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$GITHUB_USER/Nov/main/blocklist/blocklist.txt" -OutFile "$INSTALL_DIR\blocklist.txt" -UseBasicParsing
} catch {
    Write-Host "Failed to download blocklist. It will be generated automatically later."
}

Write-Host "Generating CA Certificate..."
& "$INSTALL_DIR\$BINARY_NAME" -generate-ca

Write-Host "Installing CA Certificate (Admin required)..."
# Request admin privileges if not elevated
$identity = [System.Security.Principal.WindowsIdentity]::GetCurrent()
$principal = New-Object System.Security.Principal.WindowsPrincipal($identity)
if (-not $principal.IsInRole([System.Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "Elevating to install certificate..."
    Start-Process powershell -ArgumentList "-NoProfile -ExecutionPolicy Bypass -Command `"Import-Certificate -FilePath '$INSTALL_DIR\ca-cert.pem' -CertStoreLocation Cert:\LocalMachine\Root`"" -Verb RunAs
} else {
    Import-Certificate -FilePath "$INSTALL_DIR\ca-cert.pem" -CertStoreLocation Cert:\LocalMachine\Root
}

Write-Host "Setting system proxy..."
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyEnable -Value 1
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyServer -Value "127.0.0.1:8080"
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyOverride -Value "<local>"

Write-Host "Creating Native Messaging Manifest..."
$NativeDir = "$env:LOCALAPPDATA\Google\Chrome\User Data\NativeMessagingHosts"
if (-not (Test-Path -Path $NativeDir)) {
    New-Item -ItemType Directory -Path $NativeDir | Out-Null
}
$ManifestPath = "$NativeDir\com.adblock.proxy.json"
$ManifestContent = @"
{
  "name": "com.adblock.proxy",
  "description": "AdBlocker System Proxy",
  "path": "${INSTALL_DIR}\\${BINARY_NAME}",
  "type": "stdio",
  "allowed_origins": ["chrome-extension://${EXTENSION_ID}/"]
}
"@
Set-Content -Path $ManifestPath -Value $ManifestContent -Encoding UTF8

Write-Host "Configuring Scheduled Task..."
$TaskName = "AdblockProxy"
$Action = New-ScheduledTaskAction -Execute "$INSTALL_DIR\$BINARY_NAME" -Argument "-mode=proxy"
$Trigger = New-ScheduledTaskTrigger -AtLogon
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -RestartCount 3 -RestartInterval (New-TimeSpan -Minutes 1)
Register-ScheduledTask -TaskName $TaskName -Action $Action -Trigger $Trigger -Settings $Settings -Force
Start-ScheduledTask -TaskName $TaskName

Write-Host ""
Write-Host "Installation complete! Return to the extension and click 'Check Connection'."
