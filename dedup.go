// dedup.go - A simple duplicate file finder using SHA-256 hashing.
// Usage: go run dedup.go [directory]

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// FileInfo holds metadata for each file
type FileInfo struct {
	Path       string
	Size       int64
	ModTime    time.Time
	ChangeTime time.Time
}

// hashFile computes the SHA-256 hash of the file at the given path
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// getChangeTime retrieves the change (ctime) from syscall.Stat_t (Linux)
func getChangeTime(info os.FileInfo) time.Time {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return info.ModTime()
	}
	// Seconds and nanoseconds
	sec := stat.Ctim.Sec
	nsec := stat.Ctim.Nsec
	return time.Unix(sec, nsec)
}

func main() {
	// Parse directory from command-line flags
	dir := flag.String("dir", ".", "Directory to scan for duplicates")
	flag.Parse()

	hashMap := make(map[string][]FileInfo)

	// Walk through the directory tree
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot access %s: %v\n", path, err)
			return nil
		}

		// Skip non-regular files
		if !info.Mode().IsRegular() {
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error hashing %s: %v\n", path, err)
			return nil
		}

		// Collect file info
		hashMap[hash] = append(hashMap[hash], FileInfo{
			Path:       path,
			Size:       info.Size(),
			ModTime:    info.ModTime(),
			ChangeTime: getChangeTime(info),
		})

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the path %q: %v\n", *dir, err)
		os.Exit(1)
	}

	// Print duplicates
	for hash, files := range hashMap {
		if len(files) <= 1 {
			continue
		}

		fmt.Printf("\nDuplicate Hash: %s\n", hash)
		for _, fi := range files {
			fmt.Printf("  - %s\n", fi.Path)
			fmt.Printf("      Size: %d bytes\n", fi.Size)
			fmt.Printf("      Modified: %s\n", fi.ModTime.Format(time.RFC3339))
			fmt.Printf("      Changed:  %s\n", fi.ChangeTime.Format(time.RFC3339))
		}
	}
}
