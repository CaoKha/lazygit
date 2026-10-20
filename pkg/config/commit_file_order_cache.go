package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// The custom file orders of commit diffs are stored in the state dir, one file
// per diff, holding the paths one per line. The file name is a hash of the
// diff's (full) from and to hashes and of its paths, so that a diff with
// different files never gets the order.
func commitFileOrderPath(from, to string, paths []string) (string, error) {
	stateFile, err := stateFilePath(stateFileName)
	if err != nil {
		return "", err
	}
	key := strings.Join(append([]string{from, to}, slices.Sorted(slices.Values(paths))...), "\x00")
	hash := sha256.Sum256([]byte(key))
	return filepath.Join(filepath.Dir(stateFile), "commit-file-orders", hex.EncodeToString(hash[:])), nil
}

// LoadCommitFileOrder returns the saved order of a diff's files, one path per
// line, or an empty string if there is none.
func LoadCommitFileOrder(from, to string, paths []string) string {
	path, err := commitFileOrderPath(from, to, paths)
	if err != nil {
		return ""
	}
	content, _ := os.ReadFile(path)
	return string(content)
}

func SaveCommitFileOrder(from, to string, paths []string, order []string) error {
	path, err := commitFileOrderPath(from, to, paths)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// write it atomically, so that we never load half of it
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, []byte(strings.Join(order, "\n")+"\n"), 0o644); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}
