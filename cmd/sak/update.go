package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hezhizhen/sak/internal/log"
	"github.com/hezhizhen/sak/internal/version"
	"github.com/spf13/cobra"
)

const (
	updateRepo = "hezhizhen/sak"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update sak to the latest version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate()
		},
	}
	return cmd
}

func runUpdate() error {
	latestVersion, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	currentVersion := version.Version
	if currentVersion == latestVersion || "v"+currentVersion == latestVersion {
		log.Info("Already up to date (%s)", currentVersion)
		return nil
	}

	log.Info("Updating sak: %s -> %s", currentVersion, latestVersion)

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH
	fileVersion := strings.TrimPrefix(latestVersion, "v")
	filename := fmt.Sprintf("sak_%s_%s_%s.tar.gz", fileVersion, goos, goarch)
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", updateRepo, latestVersion, filename)

	log.Info("Downloading %s...", filename)

	tmpDir, err := os.MkdirTemp("", "sak-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, filename)
	if err := downloadFile(url, archivePath); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if err := extractTarGz(tmpDir, archivePath); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	newBinary := filepath.Join(tmpDir, "sak")
	if _, err := os.Stat(newBinary); err != nil {
		return fmt.Errorf("binary not found in archive")
	}

	if err := replaceBinary(execPath, newBinary); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	log.Info("Updated to %s", latestVersion)
	return nil
}

func getLatestVersion() (string, error) {
	if ghPath, err := exec.LookPath("gh"); err == nil {
		out, err := exec.Command(ghPath, "release", "view", "--repo", updateRepo, "--json", "tagName", "-q", ".tagName").Output()
		if err == nil {
			tag := strings.TrimSpace(string(out))
			if tag != "" {
				return tag, nil
			}
		}
	}

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", updateRepo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse GitHub API response: %w", err)
	}
	if release.TagName == "" {
		return "", fmt.Errorf("no tag_name in GitHub API response")
	}

	return release.TagName, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func extractTarGz(dir, archive string) error {
	cmd := exec.Command("tar", "-xzf", archive, "-C", dir)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func replaceBinary(target, source string) error {
	info, err := os.Stat(target)
	if err != nil {
		return err
	}

	input, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	return os.WriteFile(target, input, info.Mode())
}
