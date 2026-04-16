package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Pns2051/Nov/proxy/internal/config"
	"github.com/Pns2051/Nov/proxy/internal/native"
	"github.com/Pns2051/Nov/proxy/internal/proxy"
	"github.com/Pns2051/Nov/proxy/internal/updater"
)

func main() {
	mode := flag.String("mode", "proxy", "Run mode: proxy or native")
	generateCA := flag.Bool("generate-ca", false, "Generate CA certificate and key, then exit")
	flag.Parse()

	if *generateCA {
		_, err := proxy.LoadOrCreateCA(config.CACertFile, config.CAKeyFile)
		if err != nil {
			log.Fatalf("Failed to generate CA: %v", err)
		}
		fmt.Println("CA certificate and key generated successfully.")
		os.Exit(0)
	}

	caCert, err := proxy.LoadOrCreateCA(config.CACertFile, config.CAKeyFile)
	if err != nil {
		log.Fatalf("Failed to load or create CA: %v", err)
	}

	adBlockerProxy := proxy.New(caCert)

	if err := adBlockerProxy.Blocklist.LoadFromFile(config.BlocklistFile); err != nil {
		log.Printf("Could not load blocklist from file: %v (starting empty)", err)
	} else {
		log.Printf("Loaded %d domains from %s", adBlockerProxy.Blocklist.Size(), config.BlocklistFile)
	}

	adBlockerProxy.Blocklist.StartBackgroundUpdater(24*time.Hour, []string{
		config.PrimaryBlocklistURL,
		config.FallbackBlocklistURL,
	})

	if *mode == "native" {
		native.RunNativeHost(adBlockerProxy)
	} else {
		go func() {
			err := adBlockerProxy.Start(config.ProxyAddr)
			if err != nil {
				log.Fatalf("Failed to start proxy server: %v", err)
			}
		}()

		go func() {
			// Initial delay so we don't block proxy startup
			time.Sleep(10 * time.Second)
			for {
				err := updater.CheckAndUpdate(config.Version)
				if err != nil {
					log.Printf("Update check failed: %v", err)
				}
				time.Sleep(24 * time.Hour)
			}
		}()

		select {}
	}
}
