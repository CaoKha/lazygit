package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SortCommitFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sort the files of a commit with the configured sort command",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.CommitFilesSortCommand = "sort -r"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("other", "")
		shell.Commit("other commit")
		shell.CreateFileAndAdd("dir/a", "")
		shell.CreateFileAndAdd("dir/b", "")
		shell.CreateFileAndAdd("c", "")
		shell.Commit("add files")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("add files").IsSelected(),
				Contains("other commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("  A c"),
				Equals("  ▼ dir"),
				Equals("    A a"),
				Equals("    A b"),
			).
			NavigateToLine(Contains("A a")).
			Press(keys.CommitFiles.SortFiles).
			// switches to the flat view, keeping the selection
			Lines(
				Equals("A dir/b"),
				Equals("A dir/a").IsSelected(),
				Equals("A c"),
			).
			// the tree view is sorted as usual
			Press(keys.Files.ToggleTreeView).
			Lines(
				Equals("▼ /"),
				Equals("  A c"),
				Equals("  ▼ dir"),
				Equals("    A a").IsSelected(),
				Equals("    A b"),
			).
			Press(keys.Files.ToggleTreeView).
			Lines(
				Equals("A dir/b"),
				Equals("A dir/a").IsSelected(),
				Equals("A c"),
			).
			// the order survives a refresh
			Press(keys.Universal.Refresh).
			Lines(
				Equals("A dir/b"),
				Equals("A dir/a").IsSelected(),
				Equals("A c"),
			).
			// filtering keeps the order
			FilterOrSearch("dir").
			Lines(
				Equals("A dir/b").IsSelected(),
				Equals("A dir/a"),
			).
			Press(keys.Universal.Return).
			Lines(
				Equals("A dir/b").IsSelected(),
				Equals("A dir/a"),
				Equals("A c"),
			).
			PressEscape()

		// the order doesn't apply to other commits
		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("other commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("A other").IsSelected(),
			).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("add files")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("A c").IsSelected(),
				Equals("A dir/a"),
				Equals("A dir/b"),
			)
	},
})
