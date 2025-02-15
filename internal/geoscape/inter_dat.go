package geoscape

import (
	"encoding/binary"

	"github.com/go-restruct/restruct"
)

const (
	maxInterceptions       = 4
	interceptionByteLength = 142
)

// InterFile details the specifics of any interceptions
// currently occurring on the Geoscape.
type InterFile struct {
	Interceptions []InterceptionData
}

// Pack implements the restruct.Packer interface.
func (if_ InterFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < maxCrafts; i++ {
		data, err := restruct.Pack(order, &if_.Interceptions[i])
		if err != nil {
			return nil, err
		}
		offset := i * craftByteLength
		for j := 0; j < len(data); j++ {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

// Unpack implements the restruct.Unpacker interface.
func (s *InterFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	s.Interceptions = make([]InterceptionData, maxInterceptions)
	for i := 0; i < maxInterceptions; i++ {
		offset := i * interceptionByteLength
		data := buf[offset : offset+interceptionByteLength]
		if err := restruct.Unpack(data, order, &s.Interceptions[i]); err != nil {
			return nil, err
		}
	}
	return buf[s.SizeOf():], nil
}

// SizeOf implemtents the restruct.Sizer interface.
func (if_ InterFile) SizeOf() int {
	return maxInterceptions * interceptionByteLength
}

type InterceptionData struct {

	// Geoscape Window/Icon related bytes

	// 0x00 ( word): active/inactive window:
	// 0= active window (aircraft exists & is in intercept range)
	// 1= no active window (aircraft does not exist or does exist but is not in intercept range now)
	ActiveWindow bool `struct:"int8,invertedbool"`

	// 0x02 ( word): aircraft type for icon's picture
	// 0x04 ( word): byte A (craft number) from LOC.DAT for icon of minimized window
	// 0x06 ( word): craft LOC.DAT entry number.
	// 0x08 ( word): CRAFT.DAT entry number.
	// 0x0A ( word): UFO LOC.DAT entry number.

	// Active Window bytes

	// 0x0C ( word): UFO size (size is [5 - entry from UFO data]: vs= 0; s= 1; m= 2; l= 3; vl= 4)
	// 0x0e ( word): window minimized (=1) or active (=0)
	// 0x12 ( word): y position of window (accessed as a byte where I found it but it is most likely a word).
	// 0x14 ( word): x position of window (x position is accessed as a word)
	// 0x16 ( word): status message (starting at offset 0x397 in ENGLISH.DAT)
	// 0x18 ( word): count down before the status message is cleared
	// 0x1A ( word): attack mode [0 = Standoff; 1 = Cautious; 2 = Standard; 3 = Aggressive]
	// 0x1C ( word): current distance to target (km x8) - can exceed standoff range (if byte00 = 1)
	// 0x1E ( word): requested distance to target (km x8) eg 560 = 70km (stand off)
	// 0x20 ( word): change in position for projectiles due to Xcraft stance and/or UFO escape speed delta (difference between speed of target ufo and xcom craft), "approach speed "
	// 0x22 ( word): current air speed
	// 0x24-0x41 (word): 15 entries for the positions of launched projectiles from the left weapon
	// 0x42-0x5F (word): 15 entries for the positions of launched projectiles from the right weapon
	// 0x60 ( word): Left weapon type
	// 0x62 ( word): Right weapon type
	// 0x64 ( word): Left weapon reload timer
	// 0x66 ( word): Right weapon reload timer
	// 0x68 ( word): explosion/hit radius value (and timer to delay the interception window closing until this is zero).
	// 0x6A ( word): timer to prevent a craft from firing or being fired upon. Set when the UFO is downed/destroyed.
	// 0x6C ( word): UFO return fire intense level
	// 0x6E ( word): timer to prevent a craft from firing or being fired upon. Set when an Xcraft is destroyed.
	// 0x70 ( word): UFO side view preview opened
	// 0x72 ( word): XCOM craft type
	// 0x74 ( word): UFO craft type
	// 0x76 (dword): pointer to CRAFT.DAT offset of the attacking craft
	// 0x7A (dword): pointer to attacking craft stats offset
	// 0x7E (dword): pointer to attacking craft left weapon stats offset
	// 0x82 (dword): pointer to attacking craft right weapon stats offset
	// 0x86 (dword): pointer to CRAFT.DAT offset of the attacked UFO
	// 0x8A (dword): pointer to UFO stats offset
}

// SizeOf implemtents the restruct.Sizer interface.
func (id InterceptionData) SizeOf() int {
	return 142
}
