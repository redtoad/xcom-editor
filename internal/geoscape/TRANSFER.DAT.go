package geoscape

import (
	"encoding/binary"

	"github.com/go-restruct/restruct"
)

const transferByteLength = 8
const maxTransfers = 100

// TRANSFER_DAT contains all information about items in transit. Each record is
// 8 bytes long and is fixed at 100 entries thus a fixed size of 800.
type TRANSFER_DAT struct {
	Transfers []Transfer
}

func (s TRANSFER_DAT) SizeOf() int {
	return transferByteLength * maxTransfers
}

func (t TRANSFER_DAT) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < maxTransfers; i++ {
		data, err := restruct.Pack(order, &t.Transfers[i])
		if err != nil {
			return nil, err
		}
		offset := i * transferByteLength
		for j := 0; j < len(data); j++ {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

func (t *TRANSFER_DAT) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	t.Transfers = make([]Transfer, maxTransfers)
	for i := 0; i < maxTransfers; i++ {
		offset := i * transferByteLength
		data := buf[offset : offset+transferByteLength]
		if err := restruct.Unpack(data, order, &t.Transfers[i]); err != nil {
			return nil, err
		}
	}
	return buf[t.SizeOf():], nil
}

type Transfer struct {

	// Base the item is coming from (as indexed in LOC.DAT). 255 if the item is purchased and thus no base of origin
	Origin int `struct:"int8"`

	// Base the item is going to (again from LOC.DAT). 255 should not be used here
	Destination int `struct:"int8"`

	// Hours left in transit. NOTE: Setting this to 0 will make the game think it has been completed already
	HoursLeft int `struct:"int8"`

	// Offset 3 (1 Byte) - Item Type. This also affects what can be used in the next offset. Possible/observed values:
	Type int `stuct:"int8"`

	// Offset 4-5 (2 Bytes) - Reference number. The meaning of this value depends on the above Item Type value.
	ReferenceNumber int `stuct:"int16"`

	// Offset 6 (1 Byte) - Quantity. Also the entry is ignored if this value is 0, thus there can be invalid data in the other entries but they will always have this byte set to 0.
	Quantity int `stuct:"int8"`
}

// SizeOf implements restruct.Sizer
func (t Transfer) SizeOf() int {
	return transferByteLength
}
