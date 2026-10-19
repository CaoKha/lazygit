# CaoKha/lazygit fork

This fork is upstream [jesseduffield/lazygit](https://github.com/jesseduffield/lazygit) at its
latest release plus:

- **AI sort for commit files**: `<c-a>` in the commit files panel ("AI sort files") pipes the
  commit's paths, one per line, to `git.commitFilesSortCommand` (by default `claude -p …`) and
  shows the files in the order it prints, in flat view. If the command fails or does not print
  exactly the same paths, a toast says why and the list stays as it was.

The fork's branch is `main`, which tracks upstream's release tags (`vX.Y.Z`), not `master`.
`master` is upstream's and is left alone.

## How updates flow

1. **You push to `main`** → [Fork release](.github/workflows/fork-release.yml) runs the unit
   tests, cross-compiles `darwin-arm64` and `linux-amd64`, publishes them as a release
   (`v<upstream>-fork.<date>.<10000 + run>`, e.g. `v0.66.2-fork.20261019.10001`, with
   `SHA256SUMS`), then uploads `latest.json` (version, tag, commit, and each asset's URL and
   SHA-256) to the fixed **`channel`** release. `lazygit --version` shows that version.
   Markdown-only pushes build nothing; a build that `main` has moved past is not published.
   The newest 10 builds are kept.
2. **Upstream publishes a release** → [Fork upstream sync](.github/workflows/fork-upstream-sync.yml)
   (every 6 hours, on every push to `main`, or by hand) merges the newest upstream `vX.Y.Z` tag
   into `main` and pushes it, which starts a release. It writes `upstream-status.json` to the
   `channel` release and keeps an "Upstream sync needs attention" issue open while it is stuck:
   - `conflict`: the merge conflicts. Resolve it locally:

     ```sh
     git fetch upstream --tags
     git merge <tag>   # fix conflicts
     git commit && git push
     ```

   - `needs-review`: upstream adds workflow files, which would start running here once pushed.
     Review them, then merge and push by hand as above (disable unwanted ones in Actions).
   - `push-failed`: the merge was clean but the push was rejected; check the run and the key.

   Your push reruns the sync, which then reports `in-sync` and closes the issue.

3. **Local install (LazyVim only)**: `~/.config/nvim/lua/config/lazygit_fork.lua` reads
   `https://github.com/CaoKha/lazygit/releases/download/channel/latest.json` (download URLs are
   not subject to the GitHub API rate limit) at most every 6 hours, downloads the
   `darwin-arm64` archive, verifies its SHA-256 and installs it under
   `~/.local/share/lazygit-fork/versions/<version>`, repointing `current` (the newest 3 versions
   are kept). LazyVim's lazygit terminal puts `~/.local/share/lazygit-fork/bin` first on its
   `PATH`; everywhere else `lazygit` is still the Homebrew one.

## One-time setup

- **`fork-sync` environment** (Settings → Environments) holding the `FORK_SYNC_SSH_KEY`
  secret, with no required reviewers (the sync runs unattended) and deployments limited to
  `main`. The key pair is dedicated to this repository and matches the
  "fork-upstream-sync (Actions)" deploy key (write): a deploy key may push upstream's workflow
  changes and its pushes trigger the release. Rotate it under Settings → Deploy keys.
- `main` is the default branch, so the schedule runs from it.
- Upstream's own workflows are disabled in this fork's Actions settings; only the two fork
  workflows run.
