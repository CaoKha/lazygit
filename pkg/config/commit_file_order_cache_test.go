package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitFileOrderCache(t *testing.T) {
	stateDir := t.TempDir()
	t.Setenv("CONFIG_DIR", stateDir)
	paths := []string{"a", "dir/b", "c"}

	assert.Equal(t, "", LoadCommitFileOrder("from", "to", paths))

	assert.NoError(t, SaveCommitFileOrder("from", "to", paths, []string{"c", "a", "dir/b"}))
	// the order of the paths we're given doesn't matter for finding it
	assert.Equal(t, "c\na\ndir/b\n", LoadCommitFileOrder("from", "to", []string{"dir/b", "c", "a"}))

	// a different diff or file set doesn't get it
	assert.Equal(t, "", LoadCommitFileOrder("other", "to", paths))
	assert.Equal(t, "", LoadCommitFileOrder("from", "other", paths))
	assert.Equal(t, "", LoadCommitFileOrder("from", "to", []string{"a", "c"}))

	// saving again overwrites it, leaving no temp files behind
	assert.NoError(t, SaveCommitFileOrder("from", "to", paths, []string{"a", "c", "dir/b"}))
	assert.Equal(t, "a\nc\ndir/b\n", LoadCommitFileOrder("from", "to", paths))
	entries, err := os.ReadDir(filepath.Join(stateDir, "commit-file-orders"))
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
}
