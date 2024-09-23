# X-COM Enemy Unknown
# savegame file CRAFT.DAT

little_endian

while {![end]} {

    section "Craft" {

        set craft_type [uint8 "craft type"]

        switch $craft_type {
            0 { set craft_str "Skyranger" }
            1 { set craft_str "Lightning" }
            2 { set craft_str "Avenger" }
            3 { set craft_str "Interceptor" }
            4 { set craft_str "Firestorm" }
            5 { set craft_str "Small Scout" }
            6 { set craft_str "Medium Scout" }
            7 { set craft_str "Large Scout" }
            8 { set craft_str "Harvester" }
            9 { set craft_str "Abductor" }
            10 { set craft_str "Terror Ship" }
            11 { set craft_str "Battleship" }
            12 { set craft_str "Supply Ship" }
            255 { set craft_str "not used" }
            default { set craft_str "Unknown" }
        }

        sectionvalue $craft_str

        bytes 103 "data (tbd)"
    }
}
