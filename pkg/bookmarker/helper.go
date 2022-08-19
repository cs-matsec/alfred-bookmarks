package bookmarker

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func parseURL(s string) (u *url.URL, err error) {
	u, err = url.Parse(s)
	// Ignore invalid URLs
	if err != nil {
		return
	}
	if u.Host == "" {
		return u, errors.New("hostname is empty")
	}
	return
}

func getHomeDir() (string, error) {
	return os.UserHomeDir()
}

// getLatestFile returns a path to latest files in dir
func getLatestFile(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files in the directory %s", dir)
	}

	latestIndex := 0
	for i, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return "", err
		}
		latestInfo, err := files[latestIndex].Info()
		if err != nil {
			return "", err
		}
		if time.Since(info.ModTime()) <= time.Since(latestInfo.ModTime()) {
			latestIndex = i
		}
	}

	return filepath.Join(dir, files[latestIndex].Name()), nil
}

// searchSuffixDir returns a directory name of suffix ignoring case-sensitive
func searchSuffixDir(dir, suffix string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if name := file.Name(); file.IsDir() &&
			strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffix)) {
			return name, nil
		}
	}

	return "", fmt.Errorf("not found a directory of suffix (%s) in %s directory", suffix, dir)
}

// hasReadCapability return nil if the filepath stats and has read permission
func hasReadCapability(path string) error {
	const (
		setuid uint32 = 1 << (12 - 1 - iota)
		setgid
		sticky
		userRead
		userWrite
		userExecute
		groupRead
		groupWrite
		groupExecute
		otherRead
		otherWrite
		otherExecute
	)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if perm := info.Mode().Perm(); perm&os.FileMode(userRead|groupRead|otherRead) == 0 {
		return fmt.Errorf("%s does not have read permission(%s)", filepath.Base(path), perm)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", filepath.Base(path), err)
	}
	f.Close()

	return nil
}
