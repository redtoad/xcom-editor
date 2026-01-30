package geoscape

import (
	"encoding/binary"

	"github.com/go-restruct/restruct"
)

const transferByteLength = 8
const maxTransfers = 100

// TransferFile contains all information about items in transit. Each record is
// 8 bytes long and is fixed at 100 entries thus a fixed size of 800.
type TransferFile struct {
	Transfers []TransferData
}

func (tf TransferFile) SizeOf() int {
	return transferByteLength * maxTransfers
}

func (tf TransferFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < maxTransfers; i++ {
		data, err := restruct.Pack(order, &tf.Transfers[i])
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

func (tf *TransferFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	tf.Transfers = make([]TransferData, maxTransfers)
	for i := 0; i < maxTransfers; i++ {
		offset := i * transferByteLength
		data := buf[offset : offset+transferByteLength]
		if err := restruct.Unpack(data, order, &tf.Transfers[i]); err != nil {
			return nil, err
		}
	}
	return buf[tf.SizeOf():], nil
}

type TransferData struct {

    // BaseData the item is coming from (as indexed in LOC.DAT). 255 if the item is purchased and thus no base of origin
	Origin int `struct:"int8"`

    // Destination the item is going to (again from LOC.DAT). 255 should not be used here
	Destination int `struct:"int8"`

	// HoursLeft in transit. NOTE: Setting this to 0 will make the game think it has been completed already
	HoursLeft int `struct:"int8"`

	// Offset 3 (1 Byte) - Item Type. This also affects what can be used in the next offset. Possible/observed values:
	Type int `struct:"int8"`

	// Offset 4-5 (2 Bytes) - Reference number. The meaning of this value depends on the above Item Type value.
	ReferenceNumber int `struct:"int16"`

	// Offset 6 (1 Byte) - Quantity. Also the entry is ignored if this value is 0, thus there can be invalid data in the other entries but they will always have this byte set to 0.
	Quantity int `struct:"int8"`
}

// SizeOf implements restruct.Sizer
func (td TransferData) SizeOf() int {
	return transferByteLength
}
