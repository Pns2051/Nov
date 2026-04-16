package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Blocklist struct {
	mu      sync.RWMutex
	domains map[string]bool
}

func NewBlocklist() *Blocklist {
	return &Blocklist{
		domains: make(map[string]bool),
	}
}

func (b *Blocklist) LoadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	newMap := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		newMap[line] = true
	}

	b.mu.Lock()
	b.domains = newMap
	b.mu.Unlock()

	return scanner.Err()
}

func (b *Blocklist) SaveToFile(path string) error {
	b.mu.RLock()
	domains := make([]string, 0, len(b.domains))
	for d := range b.domains {
		domains = append(domains, d)
	}
	b.mu.RUnlock()

	sort.Strings(domains)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, d := range domains {
		if _, err := w.WriteString(d + "\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

func (b *Blocklist) Contains(domain string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	// Check exact match
	if b.domains[domain] {
		return true
	}

	// Check subdomains
	parts := strings.Split(domain, ".")
	for i := 1; i < len(parts)-1; i++ {
		sub := strings.Join(parts[i:], ".")
		if b.domains[sub] {
			return true
		}
	}
	return false
}

func (b *Blocklist) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.domains)
}

func (b *Blocklist) UpdateFromURLs(urls []string) error {
	var body io.Reader
	var err error

	for _, u := range urls {
		resp, errGet := http.Get(u)
		if errGet == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			body = resp.Body
			err = nil
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		err = fmt.Errorf("last URL request failed: %v", errGet)
	}

	if body == nil {
		return fmt.Errorf("failed to download blocklist: %v", err)
	}

	newMap := make(map[string]bool)
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		line = strings.TrimPrefix(line, "||")
		line = strings.TrimSuffix(line, "^")
		if strings.HasPrefix(line, "0.0.0.0 ") {
			line = strings.TrimPrefix(line, "0.0.0.0 ")
		}
		line = strings.TrimSpace(line)
		
		if line != "" {
			newMap[line] = true
		}
	}

	b.mu.Lock()
	b.domains = newMap
	b.mu.Unlock()

	return b.SaveToFile("blocklist.txt")
}

func (b *Blocklist) StartBackgroundUpdater(interval time.Duration, urls []string) {
	go func() {
		err := b.UpdateFromURLs(urls)
		if err != nil {
			fmt.Printf("Initial blocklist update failed: %v\n", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			err := b.UpdateFromURLs(urls)
			if err != nil {
				fmt.Printf("Periodic blocklist update failed: %v\n", err)
			}
		}
	}()
}
