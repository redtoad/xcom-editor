package geoscape

// UiglobFile stores counters used for naming new crafts, UFOs, missions and
// waypoints, the current year and the monthly scores.
//
// The file is 48 bytes long.
//
// https://www.ufopaedia.org/index.php/UIGLOB.DAT
type UiglobFile struct {
	// 0: Number of the next spotted UFO.
	NextUFO uint16
	// 2: Number of the next terror mission.
	NextTerrorSite uint16
	// 4: Number of the next alien base.
	NextAlienBase uint16
	// 6: Number of minimized interception windows.
	MinimizedInterceptionWindows uint16
	// 8: Months elapsed since the start of the game. This controls the alien
	// stages and which (early/mid/late) weapons the aliens carry.
	MonthsElapsed uint16
	// 10: Number of the next waypoint.
	NextWaypoint uint16
	// 12: Number of the next Skyranger.
	NextSkyranger uint16
	// 14: Number of the next Lightning.
	NextLightning uint16
	// 16: Number of the next Avenger.
	NextAvenger uint16
	// 18: Number of the next Interceptor.
	NextInterceptor uint16
	// 20: Number of the next Firestorm.
	NextFirestorm uint16
	// 22: Current year (e.g. 1999).
	Year uint16
	// 24-47: Score for each month since the start of the game (index 0 is the
	// first month). Scores can become negative.
	MonthlyScores [12]int16
}
