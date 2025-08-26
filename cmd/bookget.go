package main

import (
	"bookget/app"
	"bookget/config"
	"bookget/pkg/queue"
	"bookget/pkg/version"
	"bookget/router"
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	wg             sync.WaitGroup
	versionChecker = version.NewChecker(
		config.Version,
		"deweizhu", // GitHub repository owner
		"bookget",  // GitHub repository name
	)
)

func main() {
	ctx := context.Background()

	// Initialize configuration
	if !initializeConfig(ctx) {
		return
	}

	// Check for updates
	checkForUpdates()

	// Execute based on run mode
	executeByRunMode(ctx)
}

// initializeConfig handles configuration initialization
func initializeConfig(ctx context.Context) bool {
	if !config.Init(ctx) {
		return false
	}
	return true
}

// executeByRunMode executes based on the run mode
func executeByRunMode(ctx context.Context) {
	switch determineRunMode() {
	case RunModeSingleURL:
		executeSingleURL(ctx, config.Conf.DUrl)
	case RunModeBatchURLs:
		executeBatchURLs()
	case RunModeInteractive:
		runInteractiveMode(ctx)
	case RunModeInteractiveImage:
		runInteractiveModeImage(ctx)
	}

	log.Println("Download complete.")
}

type RunMode int

const (
	RunModeSingleURL RunMode = iota
	RunModeBatchURLs
	RunModeInteractive
	RunModeInteractiveImage
)

// determineRunMode determines the run mode
func determineRunMode() RunMode {
	if config.Conf.DownloaderMode == 1 {
		return RunModeInteractiveImage
	}
	if config.Conf.DUrl != "" {
		return RunModeSingleURL
	}
	if hasValidURLsFile() {
		return RunModeBatchURLs
	}
	return RunModeInteractive
}

// hasValidURLsFile checks if there is a valid URLs file
func hasValidURLsFile() bool {
	f, err := os.Stat(config.Conf.UrlsFile)
	return err == nil && f.Size() > 0
}

// executeSingleURL handles single URL mode
func executeSingleURL(ctx context.Context, rawUrl string) {
	if err := processURL(ctx, rawUrl); err != nil {
		log.Println(err)
	}
}

// executeBatchURLs handles batch URLs mode
func executeBatchURLs() {
	allUrls, err := loadAndFilterURLs(config.Conf.UrlsFile)
	if err != nil {
		log.Println(err)
		return
	}

	q := queue.NewConcurrentQueue(int(config.Conf.Threads))
	if config.Conf.DownloaderMode == 1 {
		processURLsDownloaderMode(q, allUrls)
	} else {
		processURLsManual(q, allUrls)
	}
	wg.Wait()
}

// runInteractiveMode runs interactive mode
func runInteractiveMode(ctx context.Context) {
	//cleanupCookieFile()
	for {
		rawUrl, err := readURLFromInput()
		if err != nil {
			break
		}

		if err = processURL(ctx, rawUrl); err != nil {
			log.Println(err)
		}
	}
}

// runInteractiveModeImage runs interactive mode: image download
func runInteractiveModeImage(ctx context.Context) {
	//cleanupCookieFile()
	app.NewImageDownloader().Run("")
}

// loadAndFilterURLs loads and filters URLs
func loadAndFilterURLs(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read URL file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var urls []string
	for _, line := range lines {
		sUrl := strings.TrimSpace(strings.Trim(line, "\r"))
		if isValidURL(sUrl) {
			urls = append(urls, sUrl)
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("no valid URLs found in URL file")
	}

	return urls, nil
}

// isValidURL validates if URL is valid
func isValidURL(url string) bool {
	return url != "" && strings.HasPrefix(url, "http")
}

// processURLsDownloaderMode handles URLs in auto-detection mode
func processURLsDownloaderMode(q *queue.ConcurrentQueue, allUrls []string) {
	for _, v := range allUrls {
		wg.Add(1)
		rawURL := v // Create local variable for closure use
		q.Go(func() {
			defer wg.Done()
			processURLSet("bookget", rawURL)
		})
	}
}

// processURLsManual handles URLs in manual mode
func processURLsManual(q *queue.ConcurrentQueue, allUrls []string) {
	for _, v := range allUrls {
		u, err := url.Parse(v)
		if err != nil {
			log.Printf("URL parsing failed: %s, error: %v\n", v, err)
			continue
		}

		wg.Add(1)
		rawURL := v // Create local variable for closure use
		q.Go(func() {
			defer wg.Done()
			processURLSet(u.Host, rawURL)
		})
	}
}

// processURLSet processes a group of URLs
func processURLSet(siteID string, rawUrl string) {
	result, err := router.FactoryRouter(siteID, rawUrl)
	if err != nil {
		log.Println(err)
		return
	}
	// Use result
	_ = result
}

// readURLFromInput reads URL from user input
func readURLFromInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()
	fmt.Println("Enter an URL:")
	fmt.Print("-> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(input), nil
}

// processURL processes a single URL
func processURL(ctx context.Context, rawUrl string) error {
	rawURL := strings.TrimSpace(rawUrl)
	if !isValidURL(rawURL) {
		return fmt.Errorf("invalid URL: %s", rawUrl)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("URL parsing failed: %w", err)
	}

	result, err := router.FactoryRouter(u.Host, rawURL)
	if err != nil {
		log.Println(err)
		return err
	}
	// Use result
	_ = result

	return nil
}

// cleanupCookieFile cleans up cookie file
func cleanupCookieFile() {
	if err := os.Remove(config.Conf.CookieFile); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to clean up cookie file: %v\n", err)
	}
}

// checkForUpdates checks for version updates
func checkForUpdates() {
	// Delete old redundant directories (iterate through a few version numbers, this line can be removed)
	_ = os.RemoveAll(config.BookgetHomeDir())
	latestVersion, updateAvailable, err := versionChecker.CheckForUpdate()
	if err != nil {
		log.Printf("Version check failed: %v\n", err)
		return
	}

	if updateAvailable {
		fmt.Printf("\nNew version available: %s (current version: %s)\n", latestVersion, versionChecker.CurrentVersion)
		fmt.Printf("Please visit https://github.com/deweizhu/bookget/releases/latest to upgrade.\n\n")
	} else if latestVersion != "" {
		fmt.Printf("Current version is already the latest: %s\n", versionChecker.CurrentVersion)
	}
}
