package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal/battlescape"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
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

// missionJS summarizes an ongoing mission.
type missionJS struct {
	Turn        int    `json:"turn"`
	CraftName   string `json:"craftName"`
	XcomAlive   int    `json:"xcomAlive"`
	XcomTotal   int    `json:"xcomTotal"`
	AliensAlive int    `json:"aliensAlive"`
	AliensTotal int    `json:"aliensTotal"`
	Score       int    `json:"score"`
}

// missionInfo returns the mission summary of a savegame, or nil if it was not
// saved during a mission.
func missionInfo(sg *savegame.Savegame) *missionJS {
	status := sg.MissionStatus()
	if status == nil {
		return nil
	}
	info := &missionJS{
		Turn:        status.Turn,
		XcomAlive:   status.XcomAlive,
		XcomTotal:   status.XcomTotal,
		AliensAlive: status.AliensAlive,
		AliensTotal: status.AliensTotal,
		Score:       status.Score,
	}
	if status.Craft != nil {
		info.CraftName = status.Craft.Name()
	}
	return info
}

// GET /api/games
func handleListGames(w http.ResponseWriter, r *http.Request) {
	type gameSummary struct {
		Slot         string     `json:"slot"`
		Title        string     `json:"title"`
		Time         string     `json:"time"`
		SoldierCount int        `json:"soldierCount"`
		BaseCount    int        `json:"baseCount"`
		CraftCount   int        `json:"craftCount"`
		Mission      *missionJS `json:"mission"`
		// Stale is true if the savegame was changed on disk since it was
		// loaded and needs to be reloaded.
		Stale bool `json:"stale"`
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
			Mission:      missionInfo(entry.sg),
			Stale:        entry.stale,
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
		"mission":      missionInfo(sg),
		"stale":        entry.stale,
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
		// LastRank is the (estimated) rank before being killed in action.
		LastRank  string `json:"lastRank"`
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
			LastRank:  s.LastRank().String(),
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
				"lastRank":         s.LastRank().String(),
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
				"gender":           s.Gender().String(),
				"appearance":       s.Appearance().String(),
				"images":           soldierImagesFor(s.Gender(), s.Appearance()),
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

	result := map[string]interface{}{
		"index":      b.Index(),
		"name":       b.Name(),
		"active":     b.Active(),
		"engineers":  b.Engineers(),
		"scientists": b.Scientists(),
		"coord":      b.Coord().String(),
		"tiles":      tiles,
		"inventory":  inventoryMap,
	}
	// sprites to draw the base with the graphics of the game
	if baseBits != nil {
		var crafts []hangarCraft
		for _, c := range entry.sg.Crafts() {
			if base := c.Base(); base != nil && base.Index() == b.Index() && c.Type() < geoscape.SmallScout {
				crafts = append(crafts, hangarCraft{Type: c.Type(), Out: c.Status() == geoscape.Out})
			}
		}
		result["layers"] = baseLayers(b.Tiles(), crafts)
	}
	writeJSON(w, result)
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
		Damaged  bool   `json:"damaged"`
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
			Damaged:  c.Damaged(),
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

// POST /api/games/{slot}/soldiers/{idx}/resurrect
func handleResurrectSoldier(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soldier index")
		return
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	for _, s := range entry.sg.Soldiers() {
		if s.Index() == idx {
			if !s.Resurrect() {
				writeError(w, http.StatusBadRequest, "soldier is not dead")
				return
			}
			writeJSON(w, map[string]string{"status": "ok"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "soldier not found")
}

// POST /api/games/{slot}/soldiers/{idx}/heal
func handleHealSoldier(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soldier index")
		return
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	for _, s := range entry.sg.Soldiers() {
		if s.Index() == idx {
			if !s.HealWounds() {
				writeError(w, http.StatusBadRequest, "soldier is not wounded")
				return
			}
			writeJSON(w, map[string]string{"status": "ok"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "soldier not found")
}

// POST /api/games/{slot}/actions/personal-armour
func handlePersonalArmour(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	n := entry.sg.EquipPersonalArmour()
	writeJSON(w, map[string]interface{}{"status": "ok", "equipped": n})
}

// POST /api/games/{slot}/actions/repair-refuel-crafts
func handleRepairRefuelCrafts(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	n := entry.sg.RepairAndRefuelAllCrafts()
	writeJSON(w, map[string]interface{}{"status": "ok", "crafts": n})
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
	entry.markSaved()
	log.Printf("Saved game: %s", mux.Vars(r)["slot"])
	writeJSON(w, map[string]string{"status": "ok"})
}

// craftDetails holds details about a craft shown on the globe.
type craftDetails struct {
	Type        string `json:"type"`
	Status      string `json:"status"`
	Altitude    string `json:"altitude"`
	Speed       int    `json:"speed"`
	Fuel        int    `json:"fuel"`
	Damage      int    `json:"damage"`
	BaseName    string `json:"baseName,omitempty"`
	Mission     string `json:"mission,omitempty"`
	Zone        string `json:"zone,omitempty"`
	Race        string `json:"race,omitempty"`
	Destination string `json:"destination,omitempty"`
	// Damaged is true if the entry in CRAFT.DAT is (nearly) empty.
	Damaged bool `json:"damaged,omitempty"`
}

// radarJS is the detection capability of a radar facility of a base.
type radarJS struct {
	Facility string `json:"facility"`
	Range    int    `json:"range"`  // nautical miles
	Chance   int    `json:"chance"` // percent
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
		Type      string        `json:"type"`
		TypeCode  int           `json:"typeCode"`
		Name      string        `json:"name"`
		Coord     coordJS       `json:"coord"`
		DestCoord *coordJS      `json:"destCoord,omitempty"`
		Craft     *craftDetails `json:"craft,omitempty"`
		Radars    []radarJS     `json:"radars,omitempty"`
	}

	result := make([]locationSummary, 0)
	for _, loc := range sg.Locations() {
		name := locationName(sg, loc)

		c := loc.Coord()
		summary := locationSummary{
			Type:     loc.Type.String(),
			TypeCode: int(loc.Type),
			Name:     name,
			Coord:    coordJS{Lat: c.Lat, Lon: c.Lon},
		}

		if loc.Type == geoscape.XCOMBase {
			if b := sg.Base(loc.TableReference); b != nil {
				for _, radar := range b.Radars() {
					summary.Radars = append(summary.Radars, radarJS{
						Facility: radar.Facility.String(),
						Range:    radar.Range,
						Chance:   radar.Chance,
					})
				}
			}
		}

		if loc.Type == geoscape.XCOMShip || loc.Type == geoscape.AlienShip {
			if craft := sg.Craft(loc.TableReference); craft != nil {
				if dest := craft.Destination(); dest != nil {
					summary.DestCoord = &coordJS{Lat: dest.Lat, Lon: dest.Lon}
				} else if wp := craft.NextWaypoint(); wp != nil {
					summary.DestCoord = &coordJS{Lat: wp.Lat, Lon: wp.Lon}
				}
				details := &craftDetails{
					Type:     craft.Type().String(),
					Status:   craft.Status().String(),
					Altitude: craft.Altitude().String(),
					Speed:    craft.Speed(),
					Fuel:     craft.Fuel(),
					Damage:   craft.Damage(),
					Damaged:  craft.Damaged(),
				}
				if loc.Type == geoscape.AlienShip {
					details.Mission = craft.MissionType()
					details.Zone = craft.MissionZone()
					details.Race = craft.AlienRace()
				} else if b := craft.Base(); b != nil {
					details.BaseName = b.Name()
				}
				if dest := craft.DestinationLocation(); dest != nil {
					details.Destination = locationName(sg, *dest)
				}
				summary.Craft = details
			}
		}

		result = append(result, summary)
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
	if err := entry.reload(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Reloaded game: %s", slot)
	writeJSON(w, map[string]string{"status": "ok"})
}

// locationName returns a display name for a location.
func locationName(sg *savegame.Savegame, loc savegame.Location) string {
	switch loc.Type {
	case geoscape.XCOMBase:
		if b := sg.Base(loc.TableReference); b != nil {
			return b.Name()
		}
		return "X-COM Base"
	case geoscape.XCOMShip:
		if c := sg.Craft(loc.TableReference); c != nil {
			return c.Name()
		}
		return fmt.Sprintf("CRAFT-%d", loc.CountSuffix)
	case geoscape.AlienShip:
		if c := sg.Craft(loc.TableReference); c != nil {
			return fmt.Sprintf("UFO-%d (%s, %s)", loc.CountSuffix, c.Type(), c.MissionType())
		}
		return fmt.Sprintf("UFO-%d", loc.CountSuffix)
	case geoscape.AlienBase:
		return "Alien Base"
	case geoscape.CrashSite:
		return fmt.Sprintf("Crash Site-%d", loc.CountSuffix)
	case geoscape.LandedUFO:
		return fmt.Sprintf("Landed UFO-%d", loc.CountSuffix)
	case geoscape.TerrorSite:
		return fmt.Sprintf("Terror Site-%d", loc.CountSuffix)
	default:
		return loc.Type.String()
	}
}

// Flags of the tiles in battleMapJS.Tiles.
const (
	tileGround    = 1 << 0
	tileWestWall  = 1 << 1
	tileNorthWall = 1 << 2
	tileObject    = 1 << 3
	tileSeen      = 1 << 4
)

// battleUnitJS is a unit on the Battlescape map.
type battleUnitJS struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Z           int    `json:"z"`
	Large       bool   `json:"large"`
	Alive       bool   `json:"alive"`
	Unconscious bool   `json:"unconscious"`
	Visible     bool   `json:"visible"`
	Health      int    `json:"health"`
	MaxHealth   int    `json:"maxHealth"`
	TimeUnits   int    `json:"timeUnits"`
	MaxTU       int    `json:"maxTimeUnits"`
	Morale      int    `json:"morale"`
	Tank        bool   `json:"tank"`
}

// battleCraftJS is the area covered by a craft on the Battlescape map.
type battleCraftJS struct {
	Alien  bool   `json:"alien"`
	Name   string `json:"name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// battleMapJS is the Battlescape map of an ongoing mission.
type battleMapJS struct {
	Turn  int `json:"turn"`
	SizeX int `json:"sizeX"`
	SizeY int `json:"sizeY"`
	SizeZ int `json:"sizeZ"`
	// Tiles holds the flags (see tileGround etc.) of all tiles, ordered
	// along X, then Y, then Z (the same as the Z of the units, 0 = top level).
	Tiles  []int           `json:"tiles"`
	Units  []battleUnitJS  `json:"units"`
	Crafts []battleCraftJS `json:"crafts"`
}

// ownerName returns the name of a unit owner in UNITPOS.DAT.
func ownerName(owner uint8) string {
	switch owner {
	case battlescape.OwnerXcom:
		return "xcom"
	case battlescape.OwnerAlien:
		return "alien"
	case battlescape.OwnerCivilian:
		return "civilian"
	default:
		return "left-behind"
	}
}

// battleMap returns the Battlescape map of a savegame, or nil if it was not
// saved during a mission.
func battleMap(sg *savegame.Savegame) *battleMapJS {
	b := sg.Battle
	if b == nil {
		return nil
	}
	sx, sy, sz := b.Size()
	result := &battleMapJS{
		Turn: int(b.Globals.Turn), SizeX: sx, SizeY: sy, SizeZ: sz,
		Tiles:  make([]int, 0, sx*sy*sz),
		Units:  []battleUnitJS{},
		Crafts: []battleCraftJS{},
	}
	for z := 0; z < sz; z++ {
		for y := 0; y < sy; y++ {
			for x := 0; x < sx; x++ {
				tile, seen := b.Tile(x, y, z)
				flags := 0
				if tile.Ground != 0 {
					flags |= tileGround
				}
				if tile.WestWall != 0 {
					flags |= tileWestWall
				}
				if tile.NorthWall != 0 {
					flags |= tileNorthWall
				}
				if tile.Object != 0 {
					flags |= tileObject
				}
				if seen {
					flags |= tileSeen
				}
				result.Tiles = append(result.Tiles, flags)
			}
		}
	}
	for _, u := range b.MapUnits() {
		result.Units = append(result.Units, battleUnitJS{
			Index:       u.Index,
			Name:        strings.Join(strings.Fields(u.Unit.Name.String()), " "),
			Owner:       ownerName(u.Owner),
			X:           u.X,
			Y:           u.Y,
			Z:           u.Z,
			Large:       u.Large,
			Alive:       u.Alive(),
			Unconscious: u.Unconscious(),
			Visible:     u.Visible,
			Health:      int(u.Unit.Health),
			MaxHealth:   int(u.Unit.MaxHealth),
			TimeUnits:   int(u.Unit.TimeUnits),
			MaxTU:       int(u.Unit.MaxTimeUnits),
			Morale:      int(u.Unit.Morale),
			Tank:        u.Unit.Tank != 0,
		})
	}
	for _, c := range b.CraftAreas() {
		result.Crafts = append(result.Crafts, battleCraftJS{
			Alien: c.Alien, Name: c.Name(),
			X: c.X, Y: c.Y, Width: c.Width, Height: c.Height,
		})
	}
	return result
}

// GET /api/games/{slot}/battle
func handleGetBattleMap(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	m := battleMap(entry.sg)
	if m == nil {
		writeError(w, http.StatusNotFound, "savegame was not saved during a mission")
		return
	}
	writeJSON(w, m)
}
