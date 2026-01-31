package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

// armorToString converts a geoscape.Armor value to a display string.
// NoArmor is mapped to "None" to match the frontend select options.
var armorToString = map[geoscape.Armor]string{
	geoscape.NoArmor:        "None",
	geoscape.PersonalArmour: "Personal Armour",
	geoscape.PowerSuit:      "Power Suit",
	geoscape.FlyingSuit:     "Flying Suit",
}

// armorFromString converts a display string to a geoscape.Armor value.
// Accepts both "None" and "NoArmor" for the no-armor case.
var armorFromString = map[string]geoscape.Armor{
	"None":            geoscape.NoArmor,
	"NoArmor":         geoscape.NoArmor,
	"Personal Armour": geoscape.PersonalArmour,
	"Power Suit":      geoscape.PowerSuit,
	"Flying Suit":     geoscape.FlyingSuit,
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func getGameEntry(w http.ResponseWriter, r *http.Request) *gameEntry {
	slot := mux.Vars(r)["slot"]
	entry, ok := games[slot]
	if !ok {
		writeError(w, http.StatusNotFound, fmt.Sprintf("game %s not found", slot))
		return nil
	}
	return entry
}

// GET /api/games
func handleListGames(w http.ResponseWriter, r *http.Request) {
	type gameSummary struct {
		Slot         string `json:"slot"`
		Title        string `json:"title"`
		Time         string `json:"time"`
		SoldierCount int    `json:"soldierCount"`
		BaseCount    int    `json:"baseCount"`
		CraftCount   int    `json:"craftCount"`
	}

	result := make([]gameSummary, 0)
	for slot, entry := range games {
		entry.mu.RLock()
		result = append(result, gameSummary{
			Slot:         slot,
			Title:        entry.sg.Title(),
			Time:         entry.sg.Time().Format("2006-01-02 15:04"),
			SoldierCount: len(entry.sg.Soldiers()),
			BaseCount:    len(entry.sg.Bases()),
			CraftCount:   len(entry.sg.Crafts()),
		})
		entry.mu.RUnlock()
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Slot < result[j].Slot })
	writeJSON(w, result)
}

// GET /api/games/{slot}
func handleGetGame(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg
	writeJSON(w, map[string]interface{}{
		"slot":         mux.Vars(r)["slot"],
		"title":        sg.Title(),
		"time":         sg.Time().Format("2006-01-02 15:04"),
		"soldierCount": len(sg.Soldiers()),
		"baseCount":    len(sg.Bases()),
		"craftCount":   len(sg.Crafts()),
		"balance":      sg.Financials.CurrentBalance,
	})
}

// GET /api/games/{slot}/soldiers
func handleListSoldiers(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg

	type soldierSummary struct {
		Index     int    `json:"index"`
		Name      string `json:"name"`
		Rank      string `json:"rank"`
		BaseName  string `json:"baseName"`
		CraftName string `json:"craftName"`
		IsDead    bool   `json:"isDead"`
		IsWounded bool   `json:"isWounded"`
		Missions  int    `json:"missions"`
		Kills     int    `json:"kills"`
	}

	result := make([]soldierSummary, 0)
	for _, s := range sg.Soldiers() {
		baseName := ""
		if b := s.Base(); b != nil {
			baseName = b.Name()
		}
		craftName := ""
		if c := s.Craft(); c != nil {
			craftName = c.Name()
		}
		result = append(result, soldierSummary{
			Index:     s.Index(),
			Name:      s.Name(),
			Rank:      s.Rank().String(),
			BaseName:  baseName,
			CraftName: craftName,
			IsDead:    s.IsDead(),
			IsWounded: s.IsWounded(),
			Missions:  s.Missions(),
			Kills:     s.Kills(),
		})
	}
	writeJSON(w, result)
}

// GET /api/games/{slot}/soldiers/{idx}
func handleGetSoldier(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soldier index")
		return
	}

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	// Find soldier by index
	for _, s := range entry.sg.Soldiers() {
		if s.Index() == idx {
			baseName := ""
			if b := s.Base(); b != nil {
				baseName = b.Name()
			}
			craftName := ""
			if c := s.Craft(); c != nil {
				craftName = c.Name()
			}
			writeJSON(w, map[string]interface{}{
				"index":            s.Index(),
				"name":             s.Name(),
				"rank":             s.Rank().String(),
				"baseName":         baseName,
				"craftName":        craftName,
				"isDead":           s.IsDead(),
				"isWounded":        s.IsWounded(),
				"missions":         s.Missions(),
				"kills":            s.Kills(),
				"recoveryDays":     s.RecoveryDays(),
				"timeUnits":        s.TimeUnits(),
				"health":           s.Health(),
				"energy":           s.Energy(),
				"reactions":        s.Reactions(),
				"strength":         s.Strength(),
				"firingAccuracy":   s.FiringAccuracy(),
				"throwingAccuracy": s.ThrowingAccuracy(),
				"meleeAccuracy":    s.MeleeAccuracy(),
				"psionicStrength":  s.PsionicStrength(),
				"psionicSkill":     s.PsionicSkill(),
				"bravery":          s.Bravery(),
				"armor":            armorToString[s.Armor()],
				"gender":           s.Gender(),
				"appearance":       s.Appearance(),
				// Initial stats for editing
				"initialTimeUnits":        s.InitialTimeUnits(),
				"initialHealth":           s.InitialHealth(),
				"initialEnergy":           s.InitialEnergy(),
				"initialReactions":        s.InitialReactions(),
				"initialStrength":         s.InitialStrength(),
				"initialFiringAccuracy":   s.InitialFiringAccuracy(),
				"initialThrowingAccuracy": s.InitialThrowingAccuracy(),
				"initialMeleeAccuracy":    s.InitialMeleeAccuracy(),
				"initialPsionicStrength":  s.InitialPsionicStrength(),
				"initialPsionicSkill":     s.InitialPsionicSkill(),
				"initialBravery":          s.InitialBravery(),
			})
			return
		}
	}
	writeError(w, http.StatusNotFound, "soldier not found")
}

type soldierUpdateRequest struct {
	Name                    *string `json:"name"`
	InitialTimeUnits        *int    `json:"initialTimeUnits"`
	InitialHealth           *int    `json:"initialHealth"`
	InitialEnergy           *int    `json:"initialEnergy"`
	InitialReactions        *int    `json:"initialReactions"`
	InitialStrength         *int    `json:"initialStrength"`
	InitialFiringAccuracy   *int    `json:"initialFiringAccuracy"`
	InitialThrowingAccuracy *int    `json:"initialThrowingAccuracy"`
	InitialMeleeAccuracy    *int    `json:"initialMeleeAccuracy"`
	InitialPsionicStrength  *int    `json:"initialPsionicStrength"`
	InitialPsionicSkill     *int    `json:"initialPsionicSkill"`
	Armor                   *string `json:"armor"`
}

// PUT /api/games/{slot}/soldiers/{idx}
func handleUpdateSoldier(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soldier index")
		return
	}

	var req soldierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	for _, s := range entry.sg.Soldiers() {
		if s.Index() == idx {
			if req.Name != nil {
				s.SetName(*req.Name)
			}
			if req.InitialTimeUnits != nil {
				s.SetInitialTimeUnits(*req.InitialTimeUnits)
			}
			if req.InitialHealth != nil {
				s.SetInitialHealth(*req.InitialHealth)
			}
			if req.InitialEnergy != nil {
				s.SetInitialEnergy(*req.InitialEnergy)
			}
			if req.InitialReactions != nil {
				s.SetInitialReactions(*req.InitialReactions)
			}
			if req.InitialStrength != nil {
				s.SetInitialStrength(*req.InitialStrength)
			}
			if req.InitialFiringAccuracy != nil {
				s.SetInitialFiringAccuracy(*req.InitialFiringAccuracy)
			}
			if req.InitialThrowingAccuracy != nil {
				s.SetInitialThrowingAccuracy(*req.InitialThrowingAccuracy)
			}
			if req.InitialMeleeAccuracy != nil {
				s.SetInitialMeleeAccuracy(*req.InitialMeleeAccuracy)
			}
			if req.InitialPsionicStrength != nil {
				s.SetInitialPsionicStrength(*req.InitialPsionicStrength)
			}
			if req.InitialPsionicSkill != nil {
				s.SetInitialPsionicSkill(*req.InitialPsionicSkill)
			}
			if req.Armor != nil {
				if a, found := armorFromString[*req.Armor]; found {
					s.SetArmor(a)
				}
			}
			writeJSON(w, map[string]string{"status": "ok"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "soldier not found")
}

// GET /api/games/{slot}/bases
func handleListBases(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg

	type baseSummary struct {
		Index      int    `json:"index"`
		Name       string `json:"name"`
		Active     bool   `json:"active"`
		Engineers  int    `json:"engineers"`
		Scientists int    `json:"scientists"`
		Coord      string `json:"coord"`
	}

	result := make([]baseSummary, 0)
	for _, b := range sg.Bases() {
		result = append(result, baseSummary{
			Index:      b.Index(),
			Name:       b.Name(),
			Active:     b.Active(),
			Engineers:  b.Engineers(),
			Scientists: b.Scientists(),
			Coord:      b.Coord().String(),
		})
	}
	writeJSON(w, result)
}

// GET /api/games/{slot}/bases/{idx}
func handleGetBase(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid base index")
		return
	}

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	b := entry.sg.Base(idx)
	if b == nil {
		writeError(w, http.StatusNotFound, "base not found")
		return
	}

	type tileInfo struct {
		Type             string `json:"type"`
		DaysToCompletion int    `json:"daysToCompletion"`
	}

	tiles := make([]tileInfo, 36)
	for i, t := range b.Tiles() {
		tiles[i] = tileInfo{
			Type:             t.Type.String(),
			DaysToCompletion: t.DaysToCompletion,
		}
	}

	inventory := b.Inventory()
	inventoryMap := make(map[string]int)
	for i := 0; i < 96; i++ {
		if inventory[i] > 0 {
			inventoryMap[geoscape.Inventory(i).String()] = inventory[i]
		}
	}

	writeJSON(w, map[string]interface{}{
		"index":      b.Index(),
		"name":       b.Name(),
		"active":     b.Active(),
		"engineers":  b.Engineers(),
		"scientists": b.Scientists(),
		"coord":      b.Coord().String(),
		"tiles":      tiles,
		"inventory":  inventoryMap,
	})
}

type baseUpdateRequest struct {
	Engineers  *int `json:"engineers"`
	Scientists *int `json:"scientists"`
}

// PUT /api/games/{slot}/bases/{idx}
func handleUpdateBase(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid base index")
		return
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if idx < 0 || idx >= len(entry.sg.BasesData.Bases) {
		writeError(w, http.StatusNotFound, "base not found")
		return
	}

	var req baseUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	base := &entry.sg.BasesData.Bases[idx]
	if req.Engineers != nil {
		base.Engineers = *req.Engineers
	}
	if req.Scientists != nil {
		base.Scientists = *req.Scientists
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

// GET /api/games/{slot}/craft
func handleListCraft(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg

	type craftSummary struct {
		Index    int    `json:"index"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Status   string `json:"status"`
		Damage   int    `json:"damage"`
		Fuel     int    `json:"fuel"`
		BaseName string `json:"baseName"`
	}

	result := make([]craftSummary, 0)
	for _, c := range sg.Crafts() {
		baseName := ""
		if b := c.Base(); b != nil {
			baseName = b.Name()
		}
		result = append(result, craftSummary{
			Index:    c.Index(),
			Name:     c.Name(),
			Type:     c.Type().String(),
			Status:   c.Status().String(),
			Damage:   c.Damage(),
			Fuel:     c.Fuel(),
			BaseName: baseName,
		})
	}
	writeJSON(w, result)
}

// GET /api/games/{slot}/transfers
func handleListTransfers(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg

	type transferSummary struct {
		Index       int `json:"index"`
		Origin      int `json:"origin"`
		Destination int `json:"destination"`
		HoursLeft   int `json:"hoursLeft"`
		Type        int `json:"type"`
		Quantity    int `json:"quantity"`
	}

	result := make([]transferSummary, 0)
	for _, t := range sg.Transfers() {
		result = append(result, transferSummary{
			Index:       t.Index(),
			Origin:      t.Origin(),
			Destination: t.Destination(),
			HoursLeft:   t.HoursLeft(),
			Type:        t.Type(),
			Quantity:    t.Quantity(),
		})
	}
	writeJSON(w, result)
}

// GET /api/games/{slot}/financials
func handleGetFinancials(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg
	writeJSON(w, map[string]interface{}{
		"currentBalance": sg.Financials.CurrentBalance,
		"expenditure":    sg.Financials.Expenditure,
		"maintenance":    sg.Financials.Maintenance,
		"balance":        sg.Financials.Balance,
	})
}

type financialsUpdateRequest struct {
	CurrentBalance *int32 `json:"currentBalance"`
}

// PUT /api/games/{slot}/financials
func handleUpdateFinancials(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	var req financialsUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if req.CurrentBalance != nil {
		entry.sg.Financials.CurrentBalance = *req.CurrentBalance
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/actions/heal-all
func handleHealAll(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.sg.HealAllSoldiers()
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/actions/complete-constructions
func handleCompleteConstructions(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.sg.CompleteConstructions()
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/actions/speedup-deliveries
func handleSpeedupDeliveries(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.sg.SpeedupDelivery()
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/save
func handleSave(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	if err := entry.sg.Save(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Saved game: %s", mux.Vars(r)["slot"])
	writeJSON(w, map[string]string{"status": "ok"})
}

// GET /api/games/{slot}/locations
func handleListLocations(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	sg := entry.sg

	type coordJS struct {
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
	}
	type locationSummary struct {
		Type     string  `json:"type"`
		TypeCode int     `json:"typeCode"`
		Name     string  `json:"name"`
		Coord    coordJS `json:"coord"`
	}

	result := make([]locationSummary, 0)
	for _, loc := range sg.Locations() {
		var name string
		switch loc.Type {
		case geoscape.XCOMBase:
			if b := sg.Base(loc.TableReference); b != nil {
				name = b.Name()
			} else {
				name = "X-COM Base"
			}
		case geoscape.XCOMShip:
			if c := sg.Craft(loc.TableReference); c != nil {
				name = c.Name()
			} else {
				name = fmt.Sprintf("CRAFT-%d", loc.CountSuffix)
			}
		case geoscape.AlienShip:
			craftDetail := ""
			if c := sg.Craft(loc.TableReference); c != nil {
				craftDetail = c.Type().String()
			}
			if craftDetail != "" {
				name = fmt.Sprintf("UFO-%d (%s)", loc.CountSuffix, craftDetail)
			} else {
				name = fmt.Sprintf("UFO-%d", loc.CountSuffix)
			}
		case geoscape.AlienBase:
			name = "Alien Base"
		case geoscape.CrashSite:
			name = fmt.Sprintf("Crash Site-%d", loc.CountSuffix)
		case geoscape.LandedUFO:
			name = fmt.Sprintf("Landed UFO-%d", loc.CountSuffix)
		case geoscape.TerrorSite:
			name = fmt.Sprintf("Terror Site-%d", loc.CountSuffix)
		default:
			name = loc.Type.String()
		}

		c := loc.Coord()
		result = append(result, locationSummary{
			Type:     loc.Type.String(),
			TypeCode: int(loc.Type),
			Name:     name,
			Coord:    coordJS{Lat: c.Lat, Lon: c.Lon},
		})
	}
	writeJSON(w, result)
}

// POST /api/games/{slot}/reload
func handleReload(w http.ResponseWriter, r *http.Request) {
	slot := mux.Vars(r)["slot"]
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	if err := entry.sg.Reload(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Reloaded game: %s", slot)
	writeJSON(w, map[string]string{"status": "ok"})
}
