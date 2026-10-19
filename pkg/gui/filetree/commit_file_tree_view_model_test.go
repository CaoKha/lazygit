package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// When the tree shrinks under the selection - e.g. moving a patch out into the
// index removes a file - SetTree must keep the selection in range. Otherwise
// GetSelectedItems returns a nil node, which crashes callers such as
// canEditFiles when the options map is rendered during layout.
func TestCommitFileTreeViewModelSetTreeClampsSelectionOnShrink(t *testing.T) {
	files := []*models.CommitFile{
		{Path: "file1"},
		{Path: "file2"},
		{Path: "file3"},
	}
	viewModel := NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return files },
		common.NewDummyCommon(),
		false, // flat list
	)
	viewModel.SetTree()
	viewModel.SetSelectedLineIdx(viewModel.Len() - 1)

	// The file under the cursor goes away and the tree shrinks.
	files = []*models.CommitFile{{Path: "file1"}}
	viewModel.SetTree()

	assert.Less(t, viewModel.GetSelectedLineIdx(), viewModel.Len())
	assert.NotNil(t, viewModel.GetSelected())
	items, _, _ := viewModel.GetSelectedItems()
	assert.NotEmpty(t, items)
	for _, item := range items {
		assert.NotNil(t, item)
	}
}

func TestCommitFileTreeViewModelCustomOrder(t *testing.T) {
	files := []*models.CommitFile{
		{Path: "b/x"},
		{Path: "a/y"},
		{Path: "c"},
		{Path: "a/z"},
	}
	viewModel := NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return files },
		common.NewDummyCommon(),
		false, // flat list
	)
	paths := func() []string {
		return lo.Map(viewModel.GetAllItems(), func(node *CommitFileNode, _ int) string {
			return node.GetPath()
		})
	}

	viewModel.SetTree()
	assert.Equal(t, []string{"a/y", "a/z", "b/x", "c"}, paths())

	viewModel.SetCustomOrder([]string{"c", "a/z", "b/x", "a/y"})
	viewModel.SetTree()
	assert.Equal(t, []string{"c", "a/z", "b/x", "a/y"}, paths())

	// filtering keeps the relative order
	viewModel.SetFilter("a/", false)
	assert.Equal(t, []string{"a/z", "a/y"}, paths())
	viewModel.ClearFilter()
	assert.Equal(t, []string{"c", "a/z", "b/x", "a/y"}, paths())

	// the tree view is sorted as usual
	viewModel.ToggleShowTree()
	assert.Equal(t, []string{".", "a", "a/y", "a/z", "b", "b/x", "c"}, paths())
	viewModel.ToggleShowTree()
	assert.Equal(t, []string{"c", "a/z", "b/x", "a/y"}, paths())

	// files that aren't in the order go last, sorted by path
	files = append(files, &models.CommitFile{Path: "e"}, &models.CommitFile{Path: "d"})
	viewModel.SetTree()
	assert.Equal(t, []string{"c", "a/z", "b/x", "a/y", "d", "e"}, paths())

	viewModel.SetCustomOrder(nil)
	viewModel.SetTree()
	assert.Equal(t, []string{"a/y", "a/z", "b/x", "c", "d", "e"}, paths())
}
