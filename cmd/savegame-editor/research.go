package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// researchJS is a topic being researched at a base.
type researchJS struct {
	Base          int    `json:"base"`
	BaseName      string `json:"baseName"`
	Topic         int    `json:"topic"`
	Name          string `json:"name"`
	DaysRemaining int    `json:"daysRemaining"`
	Scientists    int    `json:"scientists"`
}

// productionJS is an item being manufactured at a base.
type productionJS struct {
	Base         int    `json:"base"`
	BaseName     string `json:"baseName"`
	Product      int    `json:"product"`
	Name         string `json:"name"`
	HoursPerUnit int    `json:"hoursPerUnit"`
	HoursDone    int    `json:"hoursDone"`
	Engineers    int    `json:"engineers"`
	Total        int    `json:"total"`
	Produced     int    `json:"produced"`
}

// projectVars parses the base and project index from the URL. It writes an
// error and returns false if one of them is not a number.
func projectVars(w http.ResponseWriter, r *http.Request, project string) (int, int, bool) {
	base, err := strconv.Atoi(mux.Vars(r)["base"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid base index")
		return 0, 0, false
	}
	idx, err := strconv.Atoi(mux.Vars(r)[project])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid "+project+" index")
		return 0, 0, false
	}
	return base, idx, true
}

// GET /api/games/{slot}/research
func handleListResearch(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()

	result := make([]researchJS, 0)
	for _, p := range entry.sg.ResearchProjects() {
		result = append(result, researchJS{
			Base:          p.Base().Index(),
			BaseName:      p.Base().Name(),
			Topic:         p.Topic(),
			Name:          p.Name(),
			DaysRemaining: p.DaysRemaining(),
			Scientists:    p.Scientists(),
		})
	}
	writeJSON(w, result)
}

// POST /api/games/{slot}/research/{base}/{topic}/speedup
func handleSpeedupResearchProject(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	base, topic, ok := projectVars(w, r, "topic")
	if !ok {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	for _, p := range entry.sg.ResearchProjects() {
		if p.Base().Index() == base && p.Topic() == topic {
			p.SpeedUp()
			writeJSON(w, map[string]string{"status": "ok"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "research project not found")
}

// POST /api/games/{slot}/actions/speedup-research
func handleSpeedupResearch(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	n := entry.sg.SpeedupResearch()
	writeJSON(w, map[string]interface{}{"status": "ok", "projects": n})
}

// GET /api/games/{slot}/production
func handleListProduction(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()

	result := make([]productionJS, 0)
	for _, p := range entry.sg.ProductionProjects() {
		result = append(result, productionJS{
			Base:         p.Base().Index(),
			BaseName:     p.Base().Name(),
			Product:      p.Product(),
			Name:         p.Name(),
			HoursPerUnit: p.HoursPerUnit(),
			HoursDone:    p.HoursDone(),
			Engineers:    p.Engineers(),
			Total:        p.Total(),
			Produced:     p.Produced(),
		})
	}
	writeJSON(w, result)
}

// POST /api/games/{slot}/production/{base}/{product}/speedup
func handleSpeedupProductionProject(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	base, product, ok := projectVars(w, r, "product")
	if !ok {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	for _, p := range entry.sg.ProductionProjects() {
		if p.Base().Index() == base && p.Product() == product {
			p.SpeedUp()
			writeJSON(w, map[string]string{"status": "ok"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "production project not found")
}

// POST /api/games/{slot}/actions/speedup-production
func handleSpeedupProduction(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	n := entry.sg.SpeedupProduction()
	writeJSON(w, map[string]interface{}{"status": "ok", "projects": n})
}
