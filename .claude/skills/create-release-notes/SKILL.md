---
name: create-release-notes
description: Use when writing or updating a changelog for an existing GitHub release. Triggers on "write release notes", "update changelog", "release notes for release", or when asked to document what changed in a release.
---

# Release Changelog

## Overview

Writes or updates the changelog of an **existing** GitHub release via `gh release edit`. No raw Git log — understandable sentences from PR titles.

## Step 1: Determine Release

If no release is specified:

```bash
gh release list --limit 5 --json tagName,publishedAt \
  --jq '.[] | "\(.tagName)  (vor \( (now - (.publishedAt | fromdateiso8601) ) / 86400 | floor ) Tagen)"'
```

Show the list and suggest the latest:
> "No release specified. The latest is **0.18.0** (3 days ago) — should I use this one?"

Wait for confirmation or selection, then continue.

## Step 2: Determine PRs for the Release

**IMPORTANT:** `gh pr list` returns only 30 PRs by default and may miss older PRs. Therefore, ALWAYS use both sources and reconcile them.

### 2a) Determine PR numbers from Git history (complete, no gaps)

```bash
PREV_TAG=$(gh release list --limit 10 --json tagName --jq '.[].tagName' \
  | grep -A1 "TARGET_TAG" | tail -1)

# All merge commits between the releases — returns ALL PR numbers
git log "$PREV_TAG..TARGET_TAG" --oneline --merges
```

Extract all PR numbers from the output (format: `Merge pull request #NR`).

**Assign non-merge commits correctly — mandatory!**

For non-squash merges, the history contains individual commits from the PR branch *without* their own PR number in the subject. These belong to the enclosing merge commit and must be assigned to the respective PR — not treated as standalone entries.

Check each commit without `#NR` in the subject for a reference in the commit body:

```bash
git log "$PREV_TAG..TARGET_TAG" --format="%H %s%n%b"
```

Commits between two merge commits (`Merge pull request #X` … `Merge pull request #Y`) belong to PR #X. Commits after the last merge commit until the end of the range belong to the last PR.

Also: Commit bodies may contain `Closes #NR` or `Fixes #NR` — extract these references as well.

### 2b) Load PR details via GitHub API

```bash
# For each found PR number, fetch details:
gh pr view NR --json number,title,body,labels,mergedAt
```

Or as a batch for all found numbers at once:

```bash
# Time range as fallback/verification
SINCE=$(git log "$PREV_TAG" -1 --format=%cI)
UNTIL=$(git log "TARGET_TAG" -1 --format=%cI)

gh pr list --state merged --base main --limit 200 \
  --json number,title,body,labels,mergedAt \
  --jq "[.[] | select(.mergedAt > \"$SINCE\" and .mergedAt <= \"$UNTIL\")] | sort_by(.mergedAt)"
```

### 2c) Reconciliation — mandatory!

Compare the PR numbers from `git log` (2a) with the results from `gh pr list` (2b).  
If a number is missing in 2b, fetch it explicitly via `gh pr view NR`.  
Only when both sources match, proceed to step 3.

### 2d) Query Dependabot security alerts — mandatory!

`gh pr list` and `git log` do not capture all CVE fixes, as Dependabot alerts can also be closed automatically without their own merge commit. Therefore, ALWAYS query the Dependabot alerts in addition:

```bash
SINCE=$(git log "$PREV_TAG" -1 --format=%cI)
UNTIL=$(git log "TARGET_TAG" -1 --format=%cI)

gh api repos/OWNER/REPO/dependabot/alerts \
  --jq "[.[] | select(.state == \"fixed\" and (.fixed_at // .auto_dismissed_at) > \"$SINCE\" and (.fixed_at // .auto_dismissed_at) <= \"$UNTIL\")] | .[] | {cve: .security_advisory.cve_id, ghsa: .security_advisory.ghsa_id, summary: .security_advisory.summary, fixed_at: (.fixed_at // .auto_dismissed_at), alert_number: .number}"
```

All found alerts go into the **Security** section of the changelog — even if there is no associated PR.  
Identifier hierarchy: CVE-ID > GHSA-ID > Advisory-Summary + Dependabot-Alert-Number.  
Alerts without an associated PR are listed without a `#NR` reference.

## Step 3: Categorize PRs

Only include populated categories:

| Category | Contains |
|-----------|---------|
| **New Features** | New user features, UI changes, new endpoints |
| **Technical Improvements** | Refactoring, performance, internal architecture |
| **Maintenance & Dependencies** | Dependency bumps, CI/CD, configuration |
| **Security** | CVE fixes, security patches |

## Step 4: Write changelog text

**Human-readable sentences**, not Git commit messages. Summarize related PRs together.

```markdown
## New Features
* [Feature Description] #123 #124

## Technical Improvements
* [Improvement] #125

## Maintenance & Dependencies
* Updated dependencies #126 #127

## Security
* Fixed CVE-XXXX-XXXXX #130

**Full Changelog**: https://github.com/OWNER/REPO/compare/PREV_TAG...TARGET_TAG
```

**Rules:**
- PR references as `#NR`, not as URL (except in the Full Changelog link)
- Multiple related PRs in one line: `#123, #124`
- English, active, concise: "Added document rules" instead of "Document rules were added"
- No technical details that end users won't understand
- Security entries from Dependabot (Step 2d): Use identifier hierarchy — prefer CVE-ID, otherwise GHSA-ID, otherwise Advisory-Summary + Dependabot-Alert-Number (e.g., "Alert #53"). CVE-IDs are not always present — never leave empty, fall back to the next level. Manual security fixes without a Dependabot origin can be entered without an ID if the fix is clearly identifiable as a security fix.
- New features in, or visible changes to, the savegame editor get a screenshot (Step 4b)

## Step 4b: Screenshots of UI changes — mandatory for visible changes!

Take screenshots automatically whenever a PR adds or visibly changes something in the savegame editor (`cmd/savegame-editor/frontend/`, or backend changes that alter what it shows, e.g. sprites or map data). Bugfixes and internal changes get no screenshot. At most one screenshot per feature, at most 3 per release.

### Which screenshots

Check which PRs changed the UI:

```bash
git diff --stat "$PREV_TAG..TARGET_TAG" -- cmd/savegame-editor/
```

For each such feature, decide from the PR description and the diff which **savegame**, **tab** and **clicks** show it best:

| Savegame (`testdata/`) | Shows |
|-----------|---------|
| `GAME_1` … `GAME_3` | Geoscape: soldiers, bases, crafts, transfers, financials, map, globe |
| `GAME_4`, `GAME_5`, `GAME_6` | Mission tab (saved during a mission); `GAME_6` has a tank |

Tabs: Mission, Soldiers, Bases, Crafts, Transfers, Financials, Map, Globe.

### Start the editor at TARGET_TAG

The screenshots must show the released state, not the current working tree. Run everything in a scratch directory (scratchpad or `mktemp -d`), never in `testdata/` — the editor watches and can write the savegame folder.

```bash
WORK=$(mktemp -d)
git worktree add --detach "$WORK/src" TARGET_TAG
mkdir -p "$WORK/game"
cp -R "$WORK"/src/testdata/GAME_* "$WORK/game/"
```

If the user has the original game folder (with `GEODATA`, `GEOGRAPH`, `TERRAIN`, `UNITS`), ask for its path and copy it to `$WORK/game` instead, so the maps and units are drawn with the game graphics. Never point the editor at the original folder.

Start the editor in the background (the frontend build in `frontend/dist` is committed and embedded; `open` of the browser fails harmlessly):

```bash
cd "$WORK/src" && go run ./cmd/savegame-editor -port 8091 "$WORK/game"
```

### Take the screenshots

Install Playwright once (`node_modules` is not committed). The script uses an installed Google Chrome; only without Chrome install Playwright's Chromium as well (`npx playwright install chromium`, a large download — run it in the background):

```bash
SCRIPTS=.claude/skills/create-release-notes/scripts
(cd $SCRIPTS && npm install)
```

Then one call per screenshot:

```bash
node $SCRIPTS/screenshot.mjs --game GAME_4 --tab Mission --click "Isometric" \
  --out "$WORK/TARGET_TAG-mission-map.png"
```

Options: `--game` (savegame in the sidebar), `--tab`, `--click "Text"` (repeatable, clicks by visible text), `--element CSS` (only this element), `--wait MS` (default 1500; use 3000+ for the globe and isometric maps), `--width`/`--height`, `--scale` (default 0.4 of the 2x screenshot, i.e. 1152px wide; keep it, so all release screenshots have the same size).

**Look at every screenshot** (Read the PNG) before using it: does it show the feature, is it fully loaded, no error message? Otherwise adjust the clicks/selector and repeat. Some features only show up in a certain state (e.g. a "Heal" button only for wounded soldiers): try the other savegames; if none of them shows the feature, leave that screenshot out and say so in Step 5 — never edit the savegames to fake the state. Name files `TARGET_TAG-<feature>.png` (lowercase, hyphens).

Afterwards stop the editor and clean up:

```bash
git worktree remove --force "$WORK/src"
```

### In the changelog

Put the screenshots directly below the feature entry they belong to, with a short alt text:

```markdown
## New Features
* Added a Mission tab that shows a map of the ongoing mission #20

  ![Mission tab with the isometric map](https://github.com/OWNER/REPO/releases/download/TARGET_TAG/TARGET_TAG-mission-map.png)
```

The URL points to the release asset uploaded in Step 6.

## Step 5: Show changelog for review

Show the full changelog text in the chat, show the screenshots to the user (send the PNG files) and ask:
> "Should I save this changelog for **TARGET_TAG**?"

Wait for explicit confirmation. Only then proceed to Step 6.

## Step 6: Update release

Upload the screenshots first (as release assets; `--clobber` replaces a screenshot of an earlier run), then update the notes:

```bash
gh release upload TARGET_TAG "$WORK"/TARGET_TAG-*.png --clobber

gh release edit TARGET_TAG --notes "$(cat <<'EOF'

...changelog text...

**Full Changelog**: https://github.com/OWNER/REPO/compare/PREV_TAG...TARGET_TAG
EOF
)"
```

Open the release page afterwards and check that the screenshots are shown.

## Step 7: Store security fixes in Memory

If the **Security** section was not empty, record the fixes in Memory under `project_security_fixes.md`:

```markdown
## Release TARGET_TAG (DATE)

* **[CVE-ID / GHSA-ID / Advisory-Summary]**
  [Short description]. Fixed via [Package] [old] → [new] (#PR). Dependabot Alert #NR.
```

This way, CVE IDs and alert numbers remain traceable even if GitHub later updates the advisory or a CVE ID is assigned retroactively.

## Common Mistakes

- **Only use `gh pr list`** → Returns a maximum of 30 PRs, older ones are missing. Always cross-check with `git log --merges` (Step 2c)
- **Treat commits without PR number as "direct commits"** → In non-squash merges, the history contains individual branch commits without their own PR number. These belong to the enclosing merge commit (Step 2a) and must be associated with the respective PR. Check commit bodies for `Closes #NR` / `Fixes #NR`.
- **Skip cross-checking** → Step 2c is mandatory, not optional
- **Ignore Dependabot alerts** → CVEs are often closed without their own PR. Step 2d (Dependabot API) is mandatory — otherwise CVEs will be missing in the Security section
- **Copy raw commit messages** → Always use the PR title as a basis, then rephrase
- **List all PRs individually** → Group related PRs together
- **English texts** → Everything in English
- **Create new release** → This skill only edits existing releases (`gh release edit`, not `gh release create`)
- **No technical details** → End users don't care about internal implementation details. Only include what is relevant to them.
- **Skip screenshots or ask the user for them** → Step 4b takes them automatically; only the original game folder is asked for (optional)
- **Screenshot of the working tree instead of TARGET_TAG** → Always start the editor from a worktree at TARGET_TAG
- **Run the editor on `testdata/` or the original game folder** → Always on a copy in a scratch directory
- **Use a screenshot without looking at it** → Half-loaded maps or error messages are worse than no screenshot
