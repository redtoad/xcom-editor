package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/redtoad/xcom-editor/savegame"
)

// maxSlots is the number of savegame slots of the game (GAME_1 ... GAME_10).
const maxSlots = 10

// watchInterval is how often the savegame folder is checked for new and
// changed savegames.
const watchInterval = 2 * time.Second

// settleTime is how long the files of a new savegame must be unchanged
// before it is loaded, so that a savegame is not loaded while the game is
// still writing it.
const settleTime = time.Second

// loadFailures holds the modification time of savegames (by folder) that
// could not be loaded, so that they are not loaded (and logged) again until
// they change. It is only used by scanGames.
var loadFailures = make(map[string]time.Time)

// newGameEntry returns an entry for a savegame that was just loaded.
func newGameEntry(sg *savegame.Savegame) *gameEntry {
	entry := &gameEntry{sg: sg}
	entry.stamp, _ = dirStamp(sg.Path)
	return entry
}

// markCurrent records that the savegame in memory matches the files on
// disk, e.g. after it was loaded or saved. The caller must hold entry.mu.
func (entry *gameEntry) markCurrent() {
	entry.stamp, _ = dirStamp(entry.sg.Path)
	entry.stale = false
}

// dirStamp returns the latest modification time of all files in a folder.
func dirStamp(dir string) (time.Time, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return time.Time{}, err
	}
	var latest time.Time
	for _, e := range entries {
		if !e.Type().IsRegular() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return time.Time{}, err
		}
		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
	}
	return latest, nil
}

// scanGames looks for savegames in the current savegame folder: it loads
// savegames that are new and marks loaded savegames as stale if their files
// were changed (e.g. because the game was saved in the game).
func scanGames() {
	stateMu.RLock()
	root := currentRoot
	loaded := make(map[string]*gameEntry, len(games))
	for slot, entry := range games {
		loaded[slot] = entry
	}
	stateMu.RUnlock()
	if root == "" {
		return // no folder loaded yet
	}

	for i := 1; i <= maxSlots; i++ {
		slot := fmt.Sprintf("GAME_%d", i)
		gamePath := filepath.Join(root, slot)

		if entry, ok := loaded[slot]; ok {
			stamp, err := dirStamp(gamePath)
			entry.mu.Lock()
			if !entry.stale && (err != nil || !stamp.Equal(entry.stamp)) {
				log.Printf("%s was changed on disk", slot)
				entry.stale = true
			}
			entry.mu.Unlock()
			continue
		}

		if _, err := os.Stat(filepath.Join(gamePath, "SAVEINFO.DAT")); err != nil {
			continue
		}
		stamp, err := dirStamp(gamePath)
		if err != nil || time.Since(stamp) < settleTime {
			continue // try again on the next scan
		}
		if failed, ok := loadFailures[gamePath]; ok && failed.Equal(stamp) {
			continue
		}
		sg, err := savegame.Load(gamePath)
		if err != nil {
			log.Printf("Warning: could not load %s: %v", slot, err)
			loadFailures[gamePath] = stamp
			continue
		}
		delete(loadFailures, gamePath)
		addGame(root, slot, newGameEntry(sg))
	}
}

// addGame adds a savegame that was found by scanGames, unless another
// folder was loaded in the meantime or the slot was loaded already.
func addGame(root, slot string, entry *gameEntry) {
	stateMu.Lock()
	defer stateMu.Unlock()
	if root != currentRoot || games[slot] != nil {
		return
	}
	games[slot] = entry
	log.Printf("Loaded %s: %s", slot, entry.sg.Title())
}

// watchGames scans the savegame folder for new and changed savegames in
// regular intervals.
func watchGames(interval time.Duration) {
	for range time.Tick(interval) {
		scanGames()
	}
}
