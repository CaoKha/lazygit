package filetree

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/samber/lo"
)

// ParseFileOrder parses the output of the sort command, which must contain each
// of the given paths exactly once, one per line.
func ParseFileOrder(tr *i18n.TranslationSet, paths []string, output string) ([]string, error) {
	remaining := lo.SliceToMap(paths, func(path string) (string, bool) { return path, true })
	order := []string{}
	for line := range strings.SplitSeq(output, "\n") {
		path := strings.TrimSpace(line)
		if path == "" {
			continue
		}
		if !remaining[path] {
			if lo.Contains(order, path) {
				return nil, fmt.Errorf(tr.SortCommitFilesDuplicateFile, path)
			}
			return nil, fmt.Errorf(tr.SortCommitFilesUnknownFile, path)
		}
		delete(remaining, path)
		order = append(order, path)
	}
	if len(remaining) == 1 {
		return nil, errors.New(tr.SortCommitFilesMissingFile)
	}
	if len(remaining) > 1 {
		return nil, fmt.Errorf(tr.SortCommitFilesMissingFiles, len(remaining))
	}
	return order, nil
}
