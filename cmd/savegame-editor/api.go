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

func getGame(w http.ResponseWriter, r *http.Request) *savegame.Savegame {
	slot := mux.Vars(r)["slot"]
	sg, ok := games[slot]
	if !ok {
		writeError(w, http.StatusNotFound, fmt.Sprintf("game %s not found", slot))
		return nil
	}
	return sg
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

	var result []gameSummary
	for slot, sg := range games {
		result = append(result, gameSummary{
			Slot:         slot,
			Title:        sg.Title(),
			Time:         sg.Time().Format("2006-01-02 15:04"),
			SoldierCount: len(sg.Soldiers()),
			BaseCount:    len(sg.Bases()),
			CraftCount:   len(sg.Crafts()),
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Slot < result[j].Slot })
	writeJSON(w, result)
}

// GET /api/games/{slot}
func handleGetGame(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
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
	sg := getGame(w, r)
	if sg == nil {
		return
	}

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

	var result []soldierSummary
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
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soldier index")
		return
	}

	// Find soldier by index
	for _, s := range sg.Soldiers() {
		if s.Index() == idx {
			baseName := ""
			if b := s.Base(); b != nil {
				baseName = b.Name()
			}
			craftName := ""
			if c := s.Craft(); c != nil {
				craftName = c.Name()
			}
			d := s.Data()
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
				"armor":            armorToString[d.Armor],
				"gender":           d.Sex.String(),
				"appearance":       d.Appearance.String(),
				// Initial stats for editing
				"initialTimeUnits":        d.InitialTimeUnits,
				"initialHealth":           d.InitialHealth,
				"initialEnergy":           d.InitialEnergy,
				"initialReactions":        d.InitialReactions,
				"initialStrength":         d.InitialStrength,
				"initialFiringAccuracy":   d.InitialFiringAccuracy,
				"initialThrowingAccuracy": d.InitialThrowingAccuracy,
				"initialMeleeAccuracy":    d.InitialMeleeAccuracy,
				"initialPsionicStrength":  d.InitialPsionicStrength,
				"initialPsionicSkill":     d.InitialPsionicSkill,
				"initialBravery":          d.InitialBravery,
			})
			return
		}
	}
	writeError(w, http.StatusNotFound, "soldier not found")
}

// PUT /api/games/{slot}/soldiers/{idx}
func handleUpdateSoldier(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid soldier index")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	for _, s := range sg.Soldiers() {
		if s.Index() == idx {
			d := s.Data()
			if v, ok := updates["name"]; ok {
				s.SetName(v.(string))
			}
			if v, ok := updates["initialTimeUnits"]; ok {
				d.InitialTimeUnits = int(v.(float64))
			}
			if v, ok := updates["initialHealth"]; ok {
				d.InitialHealth = int(v.(float64))
			}
			if v, ok := updates["initialEnergy"]; ok {
				d.InitialEnergy = int(v.(float64))
			}
			if v, ok := updates["initialReactions"]; ok {
				d.InitialReactions = int(v.(float64))
			}
			if v, ok := updates["initialStrength"]; ok {
				d.InitialStrength = int(v.(float64))
			}
			if v, ok := updates["initialFiringAccuracy"]; ok {
				d.InitialFiringAccuracy = int(v.(float64))
			}
			if v, ok := updates["initialThrowingAccuracy"]; ok {
				d.InitialThrowingAccuracy = int(v.(float64))
			}
			if v, ok := updates["initialMeleeAccuracy"]; ok {
				d.InitialMeleeAccuracy = int(v.(float64))
			}
			if v, ok := updates["initialPsionicStrength"]; ok {
				d.InitialPsionicStrength = int(v.(float64))
			}
			if v, ok := updates["initialPsionicSkill"]; ok {
				d.InitialPsionicSkill = int(v.(float64))
			}
			if v, ok := updates["armor"]; ok {
				if s, isStr := v.(string); isStr {
					if a, found := armorFromString[s]; found {
						d.Armor = a
					}
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
	sg := getGame(w, r)
	if sg == nil {
		return
	}

	type baseSummary struct {
		Index      int    `json:"index"`
		Name       string `json:"name"`
		Active     bool   `json:"active"`
		Engineers  int    `json:"engineers"`
		Scientists int    `json:"scientists"`
		Coord      string `json:"coord"`
	}

	var result []baseSummary
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
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid base index")
		return
	}

	b := sg.Base(idx)
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

// PUT /api/games/{slot}/bases/{idx}
func handleUpdateBase(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	idx, err := strconv.Atoi(mux.Vars(r)["idx"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid base index")
		return
	}

	if idx < 0 || idx >= len(sg.BasesData.Bases) {
		writeError(w, http.StatusNotFound, "base not found")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	base := &sg.BasesData.Bases[idx]
	if v, ok := updates["engineers"]; ok {
		base.Engineers = int(v.(float64))
	}
	if v, ok := updates["scientists"]; ok {
		base.Scientists = int(v.(float64))
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

// GET /api/games/{slot}/craft
func handleListCraft(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}

	type craftSummary struct {
		Index  int    `json:"index"`
		Name   string `json:"name"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Damage int    `json:"damage"`
		Fuel   int    `json:"fuel"`
	}

	var result []craftSummary
	for _, c := range sg.Crafts() {
		result = append(result, craftSummary{
			Index:  c.Index(),
			Name:   c.Name(),
			Type:   c.Type().String(),
			Status: c.Status().String(),
			Damage: c.Damage(),
			Fuel:   c.Fuel(),
		})
	}
	writeJSON(w, result)
}

// GET /api/games/{slot}/transfers
func handleListTransfers(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}

	type transferSummary struct {
		Index       int `json:"index"`
		Origin      int `json:"origin"`
		Destination int `json:"destination"`
		HoursLeft   int `json:"hoursLeft"`
		Type        int `json:"type"`
		Quantity    int `json:"quantity"`
	}

	var result []transferSummary
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
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	writeJSON(w, map[string]interface{}{
		"currentBalance": sg.Financials.CurrentBalance,
		"expenditure":    sg.Financials.Expenditure,
		"maintenance":    sg.Financials.Maintenance,
		"balance":        sg.Financials.Balance,
	})
}

// PUT /api/games/{slot}/financials
func handleUpdateFinancials(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if v, ok := updates["currentBalance"]; ok {
		sg.Financials.CurrentBalance = int32(v.(float64))
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/actions/heal-all
func handleHealAll(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	sg.HealAllSoldiers()
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/actions/complete-constructions
func handleCompleteConstructions(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	sg.CompleteConstructions()
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/actions/speedup-deliveries
func handleSpeedupDeliveries(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	sg.SpeedupDelivery()
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/save
func handleSave(w http.ResponseWriter, r *http.Request) {
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	if err := sg.Save(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Saved game: %s", mux.Vars(r)["slot"])
	writeJSON(w, map[string]string{"status": "ok"})
}

// POST /api/games/{slot}/reload
func handleReload(w http.ResponseWriter, r *http.Request) {
	slot := mux.Vars(r)["slot"]
	sg := getGame(w, r)
	if sg == nil {
		return
	}
	if err := sg.Reload(); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Reloaded game: %s", slot)
	writeJSON(w, map[string]string{"status": "ok"})
}
