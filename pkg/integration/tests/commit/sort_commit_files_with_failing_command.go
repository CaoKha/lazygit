package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SortCommitFilesWithFailingCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sorting the files of a commit leaves them unchanged when the sort command fails",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.CommitFilesSortCommand = "false"
	},
	SetupRepo: func(shell *Shell) {
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
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			NavigateToLine(Contains("A b")).
			Press(keys.CommitFiles.SortFiles).
			Tap(func() {
				t.ExpectToast(Equals("AI sort failed: exit status 1"))
			}).
			Lines(
				Equals("▼ /"),
				Equals("  A c"),
				Equals("  ▼ dir"),
				Equals("    A a"),
				Equals("    A b").IsSelected(),
			).
			IsFocused()
	},
})
