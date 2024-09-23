# X-COM Enemy Unknown 
# savegame file LOC.DAT

little_endian

while {![end]} {

    section "Location" {

        set location_type [uint8 "location type"]

        switch $location_type {
            0 { set type_str "Unused" }
            1 { set type_str "Alien Ship" }
            2 { set type_str "X-COM Ship" }
            3 { set type_str "X-COM Base" }
            4 { set type_str "Alien Base" }
            5 { set type_str "Crash Site" }
            6 { set type_str "Landed UFO" }
            7 { set type_str "Waypoint" }
            8 { set type_str "TerrorSite" }
            default { set type_str "Unknown" }
        }

        sectionvalue $type_str

        int8 "reference"
        uint16 "lon"
        int16 "lat"

        bytes 14 "data (tbd)"
    }
}
