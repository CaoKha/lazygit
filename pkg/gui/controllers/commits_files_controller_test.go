package controllers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestParseFileOrder(t *testing.T) {
	paths := []string{"a", "b/c", "d"}
	scenarios := []struct {
		name          string
		output        string
		expectedOrder []string
		expectedErr   string
	}{
		{
			name:          "permutation",
			output:        "d\na\nb/c\n",
			expectedOrder: []string{"d", "a", "b/c"},
		},
		{
			name:          "whitespace and blank lines are ignored",
			output:        "\n  d \n\n\ta\r\nb/c",
			expectedOrder: []string{"d", "a", "b/c"},
		},
		{
			name:        "missing",
			output:      "d\n",
			expectedErr: "missing 2 files",
		},
		{
			name:        "missing one",
			output:      "d\na\n",
			expectedErr: "missing 1 file",
		},
		{
			name:        "unknown",
			output:      "d\na\nb/c\ne\n",
			expectedErr: "unknown file e",
		},
		{
			name:        "duplicate",
			output:      "d\na\nd\nb/c\n",
			expectedErr: "duplicate file d",
		},
		{
			name:        "empty",
			output:      "",
			expectedErr: "missing 3 files",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			order, err := parseFileOrder(i18n.EnglishTranslationSet(), paths, s.output)
			if s.expectedErr != "" {
				assert.EqualError(t, err, s.expectedErr)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, s.expectedOrder, order)
			}
		})
	}
}
