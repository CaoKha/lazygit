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
		// fails when run a second time, so we can tell that the saved order is used instead
		config.GetUserConfig().Git.CommitFilesSortCommand = "if [ -e .git/sorted ]; then exit 1; fi; touch .git/sorted; sort -r"
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
			// leave in the tree view
			Press(keys.Files.ToggleTreeView).
			PressEscape()

		// entering the commit again shows the files as usual, and sorting them
		// again applies the saved order without running the command
		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("add files")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  A c"),
				Equals("  ▼ dir"),
				Equals("    A a"),
				Equals("    A b"),
			).
			Press(keys.CommitFiles.SortFiles).
			Lines(
				Equals("A dir/b"),
				Equals("A dir/a"),
				Equals("A c").IsSelected(),
			)
	},
})
