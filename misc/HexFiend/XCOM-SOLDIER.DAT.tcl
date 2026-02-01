# X-COM Enemy Unknown
# savegame file SOLDIER.DAT

little_endian

while {![end]} {

    section "Soldier" {

        uint16 "rank"
        uint16 "base"
        uint16 "craft"
        uint16 "craft before"
        uint16 "missions"
        uint16 "kills"
        uint16 "recovery days"
        uint16 "soldier value"
        set name [ascii 25 "name"]
        sectionvalue $name

        bytes 27 "data (tbd)"
    }
}
