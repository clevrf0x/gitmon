package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/clevrf0x/gitmon/internal/fslock"
)

var (
	repos       []string
	mu          sync.Mutex
	logger      *log.Logger
	appInstance *fslock.Lock
)

func init() {
	// Setup Logger
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Check if app is already running or not
	appInstance = fslock.New("", "gomon.pid")
	if err := appInstance.Lock(); err != nil {
		logger.Fatalf("Error: %v", err)
	}

	// Parse CLI Args
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [repo1] [repo2] ... [repoN]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	repos = flag.Args()
	if len(repos) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Printf("Received signal: %v", sig)
		cancel() // Cancel the context to signal shutdown
	}()

	// Validate repositories and start monitoring
	var wg sync.WaitGroup
	for _, repo := range repos {
		repoPath, err := filepath.Abs(repo)
		if err != nil {
			logger.Printf("Invalid path: %s\n", repo)
			continue
		}
		if isValidGitRepo(repoPath) {
			logger.Printf("Monitoring changes in %s\n", repoPath)
			wg.Add(1)
			go monitorChanges(ctx, &wg, repoPath)
		} else {
			logger.Printf("Invalid Git repository: %s\n", repoPath)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	appInstance.Unlock()
	logger.Println("Shutting down...")
}

func monitorChanges(ctx context.Context, wg *sync.WaitGroup, repoPath string) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Printf("Stopping monitoring for %s...\n", repoPath)
			return
		case <-ticker.C:
			if hasChanges(repoPath) {
				go commitChanges(repoPath)
			}
		}
	}
}

func commitChanges(repoPath string) {
	mu.Lock()
	defer mu.Unlock()
	if hasChanges(repoPath) {
		if commitAndPush(repoPath) {
			logger.Printf("Changes committed and pushed for %s\n", repoPath)
		} else {
			logger.Printf("Failed to commit and push changes for %s\n", repoPath)
		}
	}
}

func commitAndPush(repoPath string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	commitMsg := fmt.Sprintf("Auto Commit %s", time.Now().Format(time.RFC822))
	if _, err := runGitCommand(ctx, repoPath, "commit", "-am", commitMsg); err != nil {
		logger.Printf("Error committing changes for %s: %v\n", repoPath, err)
		return false
	}

	if _, err := runGitCommand(ctx, repoPath, "push", "origin", "main"); err != nil {
		logger.Printf("Error pushing changes for %s: %v\n", repoPath, err)
		return false
	}

	return true
}

func isValidGitRepo(repoPath string) bool {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return false
	}

	if _, err := runGitCommand(context.Background(), repoPath, "status"); err != nil {
		return false
	}

	return true
}

func hasChanges(repoPath string) bool {
	output, err := runGitCommand(context.Background(), repoPath, "status", "--porcelain")
	if err != nil {
		logger.Printf("Error checking status for %s: %v\n", repoPath, err)
		return false
	}
	return len(output) > 0
}

func runGitCommand(ctx context.Context, repoPath string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", repoPath}, args...)...)
	output, err := cmd.Output()
	return output, err
}
