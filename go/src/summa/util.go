package summa

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// FileExists returns true if the given path
// exists, or false otherwise
func FileExists(path string) (bool, os.FileInfo, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return true, stat, nil
	}
	if os.IsNotExist(err) {
		return false, nil, nil
	}
	return false, nil, err
}

// IsDir returns true if the given path exists and is a
// directory, or false otherwise
func IsDir(path string) bool {
	exists, stat, _ := FileExists(path)
	return exists && stat.IsDir()
}

// IsFile returns true if the given path exists and is a
// regular file, or false otherwise
func IsFile(path string) bool {
	exists, stat, _ := FileExists(path)
	return exists && stat.Mode().IsRegular()
}

// ResolvePath returns an absolute path generated from
// resolving symbolic links and relative path parts
func ResolvePath(path string, basePath string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(basePath, path)
	}

	absPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	}

	return absPath, nil
}

// UnixMilliseconds returns the current Unix timestamp in
// milliseconds since the epoch
func UnixMilliseconds() int64 {
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	return (int64(tv.Sec)*1e3 + int64(tv.Usec)/1e3)
}

// Reverse reverses a string
func Reverse(s string) string {
	n := 0
	rune := make([]rune, len(s))
	for _, r := range s {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	// Convert back to UTF-8.
	return string(rune)
}

// FromBase36 convers a base36 string to an integer
func FromBase36(s string) (int64, error) {
	return strconv.ParseInt(s, 36, 64)
}

// ToBase36 converts an integer into a base 36 string
func ToBase36(i int64) string {
	return strings.ToUpper(strconv.FormatInt(i, 36))
}
