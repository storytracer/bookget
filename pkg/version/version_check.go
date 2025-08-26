package version

import (
	"bookget/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	CacheFileName        = ".bookget_version_cache"
	DefaultCheckInterval = 24 * time.Hour
)

type Checker struct {
	CurrentVersion string
	RepoOwner      string
	RepoName       string
	CachePath      string
	LastChecked    time.Time
}

type cache struct {
	Version     string    `json:"version"`
	LastChecked time.Time `json:"last_checked"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func NewChecker(currentVersion, repoOwner, repoName string) *Checker {
	cachePath, _ := getCachePath()
	return &Checker{
		CurrentVersion: currentVersion,
		RepoOwner:      repoOwner,
		RepoName:       repoName,
		CachePath:      cachePath,
	}
}

func (c *Checker) CheckForUpdate() (string, bool, error) {
	// Check if this check should be skipped
	if time.Since(c.LastChecked) < DefaultCheckInterval {
		return "", false, nil
	}

	// Get latest version
	latestVersion, err := c.getLatestVersion()
	if err != nil {
		return "", false, fmt.Errorf("failed to get latest version: %w", err)
	}

	c.LastChecked = time.Now()

	// Compare versions
	if !c.compareVersions(latestVersion) {
		return latestVersion, true, nil
	}

	return latestVersion, false, nil
}

func (c *Checker) getLatestVersion() (string, error) {
	// First try to read from cache
	if cached, err := c.readCache(); err == nil && cached != nil {
		return cached.Version, nil
	}

	// Get latest version from GitHub API
	version, err := c.fetchFromGitHub()
	if err != nil {
		// If API fails but cache exists, return cached version
		if cached, err := c.readCache(); err == nil && cached != nil {
			return cached.Version, nil
		}
		return "", err
	}

	// Update cache
	if err := c.writeCache(version); err != nil {
		return "", fmt.Errorf("failed to update cache: %w", err)
	}

	return version, nil
}

func (c *Checker) fetchFromGitHub() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", c.RepoOwner, c.RepoName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GitHub API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}

func (c *Checker) compareVersions(latest string) bool {
	return c.CurrentVersion == latest
}

func (c *Checker) readCache() (*cache, error) {
	if c.CachePath == "" {
		return nil, nil
	}

	file, err := os.ReadFile(c.CachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var data cache
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Checker) writeCache(version string) error {
	if c.CachePath == "" {
		return nil
	}

	data := cache{
		Version:     version,
		LastChecked: time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(c.CachePath, jsonData, 0644)
}

func getCachePath() (string, error) {
	homeDir := config.UserHomeDir()
	return filepath.Join(homeDir, CacheFileName), nil
}
