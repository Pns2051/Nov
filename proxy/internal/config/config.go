package config

const (
    ProxyAddr        = "127.0.0.1:8080"
    CACertFile       = "ca-cert.pem"
    CAKeyFile        = "ca-key.pem"
    BlocklistFile    = "blocklist.txt"
)

var (
    Version               = "1.0.0"

    PrimaryBlocklistURL   = "https://raw.githubusercontent.com/Pns2051/Nov/main/blocklist/blocklist.txt"
    FallbackBlocklistURL  = "https://github.com/Pns2051/Nov/releases/latest/download/blocklist.txt"
    PrimaryVersionURL     = "https://raw.githubusercontent.com/Pns2051/Nov/main/version.txt"
    FallbackVersionURL    = "https://github.com/Pns2051/Nov/releases/latest/download/version.txt"
    PrimaryBinaryURL      = "https://raw.githubusercontent.com/Pns2051/Nov/main/dist/adblock-proxy-%s-%s"
    FallbackBinaryURL     = "https://github.com/Pns2051/Nov/releases/latest/download/adblock-proxy-%s-%s"
)
