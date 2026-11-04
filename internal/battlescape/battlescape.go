// Package battlescape
//
// These are savegame files used by the Battlescape (tactical) portion of the
// game. They are only present in savegames that were saved during a mission,
// in addition to the Geoscape files (see package geoscape).
//
// Supported files: BGLOB.DAT, GEODATA.DAT, MAP.DAT, MISDATA.DAT, OBPOS.DAT,
// ROUTES.DAT, SEEMAP.DAT, SMOKBIT.DAT, SMOKREF.DAT, SOURCEMP.DAT, TERMP.DAT,
// UNITPOS.DAT, UNITREF.DAT and WGLOB.DAT.
//
// The size of MAP.DAT, ROUTES.DAT, SEEMAP.DAT, SOURCEMP.DAT and TERMP.DAT
// depends on the size of the map (see WGLOB.DAT). The game does not truncate
// these files, so they may contain stale data from a larger map at the end.
// All bytes are kept, so that every file can be saved byte-identical to the
// original.
//
// The temporary files in MISSDAT are not supported:
//
//   - GEOSCAPE.EXE writes GEODATA.DAT, MISSION.DAT (format of MISDATA.DAT),
//     UNIPOS.DAT, UNIREF.DAT and OBPOSREF.DAT before a mission. The last three
//     use the record formats of UNITPOS.DAT, UNITREF.DAT and OBPOS.DAT, but
//     only contain the used records. TACTICAL.EXE loads them into its
//     UNITPOS/UNITREF/OBPOS tables.
//   - TACTICAL.EXE writes UNITPOS.DAT, UNIREF2.DAT, MISSION2.DAT, OBPOS2.DAT
//     (full tables), CODE.DAT and OPTION.DAT at the end of a mission.
//   - CODE.DAT is always written as uint16 1 and is not read by any program.
//   - OPTION.DAT holds the Battlescape options (layout of WGLOB.DAT bytes
//     32-43) and is read by TACTICAL.EXE at startup.
//
// References (linear addresses): GEOSCAPE.EXE writes the pre-mission files at
// 0x34f50, TACTICAL.EXE reads them at 0x31c7c and writes the post-mission
// files at 0x2b4c0 (CODE.DAT at 0x2b75f).
//
// https://www.ufopaedia.org/index.php/Saved_Game_Files
package battlescape

import (
	"encoding/binary"
	"fmt"

	"github.com/go-restruct/restruct"
)

// unpackRecords decodes buf into consecutive records of the given size.
// Trailing bytes that do not fill a complete record are an error.
func unpackRecords[T any](buf []byte, order binary.ByteOrder, size int) ([]T, error) {
	if len(buf)%size != 0 {
		return nil, fmt.Errorf("data size %d is not a multiple of record size %d", len(buf), size)
	}
	records := make([]T, len(buf)/size)
	for i := range records {
		if err := restruct.Unpack(buf[i*size:(i+1)*size], order, &records[i]); err != nil {
			return nil, fmt.Errorf("could not unpack record %d: %w", i, err)
		}
	}
	return records, nil
}

// packRecords encodes records of the given size into buf.
func packRecords[T any](buf []byte, order binary.ByteOrder, size int, records []T) error {
	for i := range records {
		data, err := restruct.Pack(order, &records[i])
		if err != nil {
			return fmt.Errorf("could not pack record %d: %w", i, err)
		}
		copy(buf[i*size:(i+1)*size], data)
	}
	return nil
}
