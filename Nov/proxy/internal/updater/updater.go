package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CheckAndUpdate(currentVersion string) error {
	primaryVersionURL := "https://cdn.jsdelivr.net/gh/Pns2051/Nov@latest/version.txt"
	fallbackVersionURL := "https://github.com/Pns2051/Nov/releases/latest/download/version.txt"

	primaryBinaryURL := "https://cdn.jsdelivr.net/gh/Pns2051/Nov@latest/adblock-proxy-%s-%s"
	fallbackBinaryURL := "https://github.com/Pns2051/Nov/releases/latest/download/adblock-proxy-%s-%s"

	var version string
	for _, u := range []string{primaryVersionURL, fallbackVersionURL} {
		resp, err := http.Get(u)
		if err == nil && resp.StatusCode == 200 {
			b, _ := io.ReadAll(resp.Body)
			version = strings.TrimSpace(string(b))
			resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	if version == "" || version == currentVersion {
		return nil
	}

	osStr := runtime.GOOS
	archStr := runtime.GOARCH

	pUrl := fmt.Sprintf(primaryBinaryURL, osStr, archStr)
	fUrl := fmt.Sprintf(fallbackBinaryURL, osStr, archStr)

	var reader io.Reader
	for _, u := range []string{pUrl, fUrl} {
		resp, err := http.Get(u)
		if err == nil && resp.StatusCode == 200 {
			reader = resp.Body
			defer resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	if reader == nil {
		return fmt.Errorf("failed to download new binary")
	}

	tmpPath := os.TempDir() + "/adblock-proxy.new"
	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	
	_, err = io.Copy(out, reader)
	out.Close()
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	if osStr == "windows" {
		batPath := os.TempDir() + "\\update.bat"
		batContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak > NUL
move /Y "%s" "%s"
start "" "%s"
del "%%~f0"
`, tmpPath, exePath, exePath)
		os.WriteFile(batPath, []byte(batContent), 0755)
		cmd := exec.Command("cmd.exe", "/C", batPath)
		cmd.Start()
		os.Exit(0)
	} else {
		os.Chmod(tmpPath, 0755)
		err = os.Rename(tmpPath, exePath)
		if err != nil {
			return err
		}
		
		cmd := exec.Command(exePath)
		cmd.Start()
		os.Exit(0)
	}
	return nil
}
