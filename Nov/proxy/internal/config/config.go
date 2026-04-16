package config

const (
    Version          = "1.0.0"
    ProxyAddr        = "127.0.0.1:8080"
    CACertFile       = "ca-cert.pem"
    CAKeyFile        = "ca-key.pem"
    BlocklistFile    = "blocklist.txt"
)

var (
    PrimaryBlocklistURL   = "https://cdn.jsdelivr.net/gh/Pns2051/Nov@latest/blocklist.txt"
    FallbackBlocklistURL  = "https://github.com/Pns2051/Nov/releases/latest/download/blocklist.txt"
    PrimaryVersionURL     = "https://cdn.jsdelivr.net/gh/Pns2051/Nov@latest/version.txt"
    FallbackVersionURL    = "https://github.com/Pns2051/Nov/releases/latest/download/version.txt"
    PrimaryBinaryURL      = "https://cdn.jsdelivr.net/gh/Pns2051/Nov@latest/adblock-proxy-%s-%s"
    FallbackBinaryURL     = "https://github.com/Pns2051/Nov/releases/latest/download/adblock-proxy-%s-%s"
)
